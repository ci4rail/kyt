/*
Copyright © 2021 Ci4Rail GmbH

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

package main

import (
	"fmt"
	"log"

	"github.com/ci4rail/kyt/kyt-api-server/cmd"
	"github.com/ci4rail/kyt/kyt-api-server/internal/controllerif"
	"github.com/ci4rail/kyt/kyt-api-server/internal/iothubservice"
)

func main() {
	log.Printf("Test Server with dummy data started")

	// switch controller to stub functions
	iothubservice.ControllerNewIOTHubServiceClient = newIOTHubStubClient
	cmd.Execute()

}

type stubbedIOTHubServiceClient struct {
	controllerif.IOTHubServices
}

func newIOTHubStubClient(connectionString string) (controllerif.IOTHubServices, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("Missing IoTHub connection string")
	}

	c := &stubbedIOTHubServiceClient{}
	return c, nil
}

// ListDeviceIDs returns a list with the device IDs of all devices of that IoT Hub
func (c *stubbedIOTHubServiceClient) ListDeviceIDs() (*[]string, error) {
	const numDevs = 1000
	var devList [numDevs]string

	for i := 0; i < numDevs; i++ {
		devList[i] = fmt.Sprintf("Device %5d", i)
	}
	slice := devList[:]
	return &slice, nil
}
