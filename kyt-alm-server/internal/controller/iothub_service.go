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

// ListRuntimeIDs returns a list with the device IDs of all devices of that IoT Hub
func (c *IOTHubServiceClient) ListRuntimeIDs(tenantID string) (*[]string, error) {
	ctx := context.Background()

	c.tenantID = tenantID
	c.deviceIDArr = nil

	// this query selects all devices and returns only the deviceId
	// Unfortunately, QueryDevices does not support paging
	query := "SELECT deviceId FROM DEVICES"
	err := c.iotClient.QueryDevices(ctx, query, c.listRuntimeIDs)

	if err != nil {
		return nil, fmt.Errorf("Error IoT Hub QueryDevices %s", err)
	}
	return &c.deviceIDArr, nil
}

// this gets called from QueryDevices once for each record (device)
func (c *IOTHubServiceClient) listRuntimeIDs(v map[string]interface{}) error {
	// This is the place where things read from IoT Hub get entered into &Device{}
	deviceID := fmt.Sprintf("%v", v["deviceId"])
	belongs, err := c.deviceBelongsToTenantAndAlm(deviceID, c.tenantID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if belongs {
		c.deviceIDArr = append(c.deviceIDArr, deviceID)
	}
	return nil
}

// ListRuntimeByID returns a list with the device IDs of all devices of that IoT Hub
func (c *IOTHubServiceClient) ListRuntimeByID(tenantID string, deviceID string) (*string, error) {
	ctx := context.Background()

	c.tenantID = tenantID
	c.deviceIDArr = nil

	// this query selects all devices and returns only the deviceId
	// Unfortunately, QueryDevices does not support paging
	query := fmt.Sprintf("SELECT * FROM DEVICES WHERE deviceId = '%s'", deviceID)
	err := c.iotClient.QueryDevices(ctx, query, c.listRuntimeIDs)

	if err != nil {
		return nil, fmt.Errorf("Error IoT Hub QueryDevices %s", err)
	}
	if len(c.deviceIDArr) > 0 {
		return &c.deviceIDArr[0], nil
	}
	return nil, fmt.Errorf("No device found with id: %s", deviceID)
}

func (c *IOTHubServiceClient) deviceBelongsToTenantAndAlm(deviceID, tenantID string) (bool, error) {
	ctx := context.Background()
	twin, err := c.iotClient.GetDeviceTwin(ctx, deviceID)
	if err != nil {
		return false, fmt.Errorf("Error reading device twin %s", err)
	}
	if twin.Tags["tenantId"] == tenantID && twin.Tags["alm"] == true {
		return true, nil
	}
	return false, nil
}

// GetConnectionState gets the connection state from the Device Twin on IoT Hub
// returns bool: 0 -> disconnected, 1 -> connected
func (c *IOTHubServiceClient) GetConnectionState(tenantID string, deviceID string) (string, error) {
	ctx := context.Background()
	twin, err := c.iotClient.GetDeviceTwin(ctx, deviceID)
	if err != nil {
		return "", fmt.Errorf("Error reading device twin %s", err)
	}
	return string(twin.ConnectionState), nil
}
