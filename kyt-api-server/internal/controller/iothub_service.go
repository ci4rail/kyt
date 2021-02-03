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

package controller

import (
	"context"
	"fmt"

	"github.com/amenzhinsky/iothub/iotservice"
	"github.com/ci4rail/kyt/kyt-api-server/internal/controllerif"
)

// IOTHubServiceClient is an Azure IoT Hub service client.
type IOTHubServiceClient struct {
	controllerif.IOTHubServices
	iotClient   *iotservice.Client
	deviceIDArr []string // filled by callback of ListDeviceIDs
}

// NewIOTHubServiceClient creates a new IOTHubServiceClient based on the connection string
// connection string can be determined with "az iot hub connection-string show"
func NewIOTHubServiceClient(connectionString string) (controllerif.IOTHubServices, error) {
	c := &IOTHubServiceClient{}

	iotClient, err := iotservice.NewFromConnectionString(connectionString)

	if err != nil {
		return nil, fmt.Errorf("Can't create IoT Hub Client %s", err)
	}

	c.iotClient = iotClient
	return c, nil
}

// ListDeviceIDs returns a list with the device IDs of all devices of that IoT Hub
func (c *IOTHubServiceClient) ListDeviceIDs() (*[]string, error) {
	ctx := context.Background()

	c.deviceIDArr = nil

	// this query selects all devices and returns only the deviceId
	// Unfortunately, QueryDevices does not support paging
	err := c.iotClient.QueryDevices(ctx, "SELECT deviceId FROM DEVICES", c.listDeviceIDsCB)

	if err != nil {
		return nil, fmt.Errorf("Error IoT Hub QueryDevices %s", err)
	}
	return &c.deviceIDArr, nil
}

// this gets called from QueryDevices once for each record (device)
func (c *IOTHubServiceClient) listDeviceIDsCB(v map[string]interface{}) error {
	deviceID := fmt.Sprintf("%v", v["deviceId"])
	c.deviceIDArr = append(c.deviceIDArr, deviceID)
	return nil
}
