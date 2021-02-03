/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package devicestate

import (
	"sync"
	"time"

	"github.com/warthog618/gpiod"
)

const (
	// baseTimeMilliSec is the minimum time the led stays in one state
	baseTimeMilliSec = 500
	// numberSteps number of blink pattern steps defined
	numberSteps = 2
)

// LedService implements an led service indicating the device connection state by
// applying a blink pattern on the device led selected.
//
// LED on: Device is connected.
// LED blinking: Device tries to connect.
// LED off: Service terminated.
type LedService struct {
	closed                 chan struct{}
	connectionStateChannel chan bool
	blinkPattern           blinkPattern
	chip                   *gpiod.Chip
	line                   *gpiod.Line
	wg                     sync.WaitGroup
}

type blinkPatterns int

const (
	off blinkPatterns = iota
	blink
	on
	exit
)

type blinkPattern struct {
	pattern blinkPatterns
	m       sync.Mutex
}

// NewLedService intialize led service
// GPIO used can be configured
func NewLedService(connectionStateChannel chan bool, gpioChip string, lineNr int) (ledService LedService, err error) {

	chip, err := gpiod.NewChip(gpioChip)
	if err != nil {
		return
	}

	line, err := chip.RequestLine(lineNr, gpiod.AsOutput(0))
	if err != nil {
		return
	}

	ledService = LedService{
		closed:                 make(chan struct{}),
		chip:                   chip,
		line:                   line,
		connectionStateChannel: connectionStateChannel}

	return
}

// Close Cleaup function for LedService
func (ledService *LedService) Close() {

	// terminate all goroutines
	close(ledService.closed)
	// wait for goroutines to finish
	ledService.wg.Wait()

	// close gpio ressources
	ledService.chip.Close()
	ledService.line.Reconfigure(gpiod.AsInput)
	ledService.line.Close()
}

// Run runs the led servie
func (ledService *LedService) Run() {

	// add both this goroutine and the controlLed goroutine to waitgroups
	ledService.wg.Add(2)

	// start with blinkpattern
	go ledService.controlLed()

	defer ledService.wg.Done()

	// wait for new data from channel from device state goroutine
	for {
		select {
		case <-ledService.closed: // close function was called
			ledService.blinkPattern.changePattern(exit) // set blink pattern to terminate
			return
		case connectionState := <-ledService.connectionStateChannel:
			// depending on the connection state change the blink pattern
			if connectionState == true {
				ledService.blinkPattern.changePattern(on)

			} else if connectionState == false {
				ledService.blinkPattern.changePattern(blink)
			}
		}
	}
}

// controlLed goroutine which executes the currently selected blink pattern
func (ledService *LedService) controlLed() {
	defer ledService.wg.Done()

	steps := map[blinkPatterns][numberSteps]int{
		off:   {0, 0},
		blink: {0, 1},
		on:    {1, 1},
		exit:  {},
	}
	stepIdx := 0
	curPattern := off
	for {
		pattern := ledService.blinkPattern.getPattern()

		if pattern == exit {
			ledService.line.SetValue(0)
			break
		}
		if pattern != curPattern {
			stepIdx = 0
			curPattern = pattern
		}

		ledVal := steps[curPattern][stepIdx]
		ledService.line.SetValue(ledVal)

		time.Sleep(baseTimeMilliSec * time.Millisecond)
		stepIdx++
		if stepIdx == numberSteps {
			stepIdx = 0
		}
	}
}

func (b *blinkPattern) changePattern(newPattern blinkPatterns) {
	b.m.Lock()
	b.pattern = newPattern
	b.m.Unlock()
}

func (b *blinkPattern) getPattern() blinkPatterns {
	b.m.Lock()
	pattern := b.pattern
	b.m.Unlock()
	return pattern
}
