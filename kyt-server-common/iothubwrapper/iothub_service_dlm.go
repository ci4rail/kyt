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
)

// ListDeviceIDs returns a list with the device IDs of all devices of that IoT Hub
func (c *IOTHubServiceClient) ListDeviceIDs(tenantID string) (*[]string, error) {
	ctx := context.Background()

	c.tenantID = tenantID
	c.deviceIDArr = nil

	// this query selects all devices and returns only the deviceId
	// Unfortunately, QueryDevices does not support paging
	query := "SELECT deviceId FROM DEVICES"
	err := c.iotClient.QueryDevices(ctx, query, c.listDeviceIDs)

	if err != nil {
		return nil, fmt.Errorf("Error IoT Hub QueryDevices %s", err)
	}
	return &c.deviceIDArr, nil
}

// this gets called from QueryDevices once for each record (device)
func (c *IOTHubServiceClient) listDeviceIDs(v map[string]interface{}) error {
	// This is the place where things read from IoT Hub get entered into &Device{}
	deviceID := fmt.Sprintf("%v", v["deviceId"])
	belongs, err := c.deviceBelongsToTenant(deviceID, c.tenantID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if belongs {
		c.deviceIDArr = append(c.deviceIDArr, deviceID)
	}
	return nil
}

// ListDeviceByID returns a list with the device IDs of all devices of that IoT Hub
func (c *IOTHubServiceClient) ListDeviceByID(tenantID string, deviceID string) (*string, error) {
	ctx := context.Background()

	c.tenantID = tenantID
	c.deviceIDArr = nil

	// this query selects all devices and returns only the deviceId
	// Unfortunately, QueryDevices does not support paging
	query := fmt.Sprintf("SELECT * FROM DEVICES WHERE deviceId = '%s'", deviceID)
	err := c.iotClient.QueryDevices(ctx, query, c.listDeviceIDs)

	if err != nil {
		return nil, fmt.Errorf("Error IoT Hub QueryDevices %s", err)
	}
	if len(c.deviceIDArr) > 0 {
		return &c.deviceIDArr[0], nil
	}
	return nil, fmt.Errorf("No device found with id: %s", deviceID)
}

func (c *IOTHubServiceClient) deviceBelongsToTenant(deviceID, tenantID string) (bool, error) {
	ctx := context.Background()
	twin, err := c.iotClient.GetDeviceTwin(ctx, deviceID)
	if err != nil {
		return false, fmt.Errorf("Error reading device twin %s", err)
	}
	if twin.Tags["tenantId"] == tenantID {
		return true, nil
	}
	return false, nil
}

// GetVersions gets the device versiones stored in IoT Hub device twin
func (c *IOTHubServiceClient) GetVersions(deviceID string) (map[string]string, error) {
	versionsMap := make(map[string]string)
	ctx := context.Background()
	twin, err := c.iotClient.GetDeviceTwin(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("Error reading device twin %s", err)
	}
	versionJSON, ok := twin.Properties.Reported["versions"]
	if ok {
		v, ok := versionJSON.(map[string]interface{})
		if ok {
			firmwareVersion, ok := v["firmwareVersion"].(string)
			if ok {
				versionsMap["firmwareVersion"] = firmwareVersion
			} else {
				return nil, fmt.Errorf("Error IoT Hub GetFirmwareVersion: no key 'version.firmwareVersion' found")
			}
		} else {
			return nil, fmt.Errorf("Error IoT Hub GetFirmwareVersion: no key 'version' found")
		}
	} else {
		return nil, fmt.Errorf("Error IoT Hub GetFirmwareVersion: no key 'version' found")
	}
	return versionsMap, nil
}
