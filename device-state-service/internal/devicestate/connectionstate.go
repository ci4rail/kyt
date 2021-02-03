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
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	// cycleTime_msec The connection status is checked every cycleTimeSec
	cycleTimeMilliSec = 500
)

// ConnectionState implements a service which continously checks the connection state
// and provides the information via channel connectionStateChannel
type ConnectionState struct {
	closed                 chan struct{}
	connectionStateChannel chan bool
	wg                     sync.WaitGroup
}

// NewConnectionState intialize connection state
func NewConnectionState(connectionStateChannel chan bool) (connectionState ConnectionState) {
	connectionState = ConnectionState{
		connectionStateChannel: connectionStateChannel,
		closed:                 make(chan struct{})}
	return
}

// Run runs the connection state service
func (connectionState *ConnectionState) Run() {

	// add this goroutine to waitgroups
	connectionState.wg.Add(1)
	defer connectionState.wg.Done()

	// pre-initialize connection state
	connState := false

	for {

		select {
		case <-connectionState.closed: // close function was called
			// terminate goroutine
			close(connectionState.connectionStateChannel)
			return
		default:
			newConnState := getNetworkStatus()
			if newConnState != connState {
				connState = newConnState
				// Log connection state change
				if connState == true {
					fmt.Println("Device switched to connected state.")
				} else {
					fmt.Println("Device switched to disconnected state.")
				}
				// Provide connection change via channel
				connectionState.connectionStateChannel <- connState
			}
		}
		time.Sleep(cycleTimeMilliSec * time.Millisecond)
	}
}

// Close Cleaup function for LedService
func (connectionState *ConnectionState) Close() {

	// terminate goroutine
	close(connectionState.closed)
	// wait for goroutine to finish
	connectionState.wg.Wait()
}

// getNetworkStatus checks if a valid internet connection is established.
// This is checked by executing a dns lookup.
func getNetworkStatus() (status bool) {
	_, err := net.LookupIP("portal.azure.com")
	if err != nil {
		status = false
	} else {
		status = true
	}
	return
}
