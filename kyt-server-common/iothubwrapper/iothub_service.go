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

package iothubwrapper

import (
	"context"
	"fmt"

	"github.com/amenzhinsky/iothub/iotservice"
)

// IOTHubServiceClient is an Azure IoT Hub service client.
type IOTHubServiceClient struct {
	iotClient   *iotservice.Client
	deviceIDArr []string // filled by callback of ListRuntimeIDs
	tenantID    string
}

// NewIOTHubServiceClient creates a new IOTHubServiceClient based on the connection string
// connection string can be determined with "az iot hub connection-string show"
func NewIOTHubServiceClient(connectionString string) (*IOTHubServiceClient, error) {
	c := &IOTHubServiceClient{}

	iotClient, err := iotservice.NewFromConnectionString(connectionString)

	if err != nil {
		return nil, fmt.Errorf("Can't create IoT Hub Client %s", err)
	}

	c.iotClient = iotClient
	return c, nil
}

// GetConnectionState gets the connection state from the Device Twin on IoT Hub
// returns bool: 0 -> disconnected, 1 -> connected
func (c *IOTHubServiceClient) GetConnectionState(deviceID string) (string, error) {
	ctx := context.Background()
	twin, err := c.iotClient.GetDeviceTwin(ctx, deviceID)
	if err != nil {
		return "", fmt.Errorf("Error reading device twin %s", err)
	}
	return string(twin.ConnectionState), nil
}
