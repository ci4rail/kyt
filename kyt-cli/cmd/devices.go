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

package cmd

import (
	"fmt"
	"os"

	api "github.com/ci4rail/kyt/kyt-cli/internal/api"
	"github.com/ci4rail/kyt/kyt-cli/openapi"
	"github.com/spf13/viper"
)

func fetchDevices(list []string) []openapi.Device {
	var devices []openapi.Device
	if len(list) > 0 {
		for _, arg := range list {
			dev, err := fetchDevicesById(arg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			devices = append(devices, dev)
		}
	} else {
		devices = fetchDevicesAll()
	}
	return devices
}

func fetchDevicesAll() []openapi.Device {
	apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
	devices, resp, err := apiClient.DeviceApi.DevicesGet(ctx).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		err := RefreshToken()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
		devices, _, err = apiClient.DeviceApi.DevicesGet(ctx).Execute()
		if err.Error() != "" {
			er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
		}
	} else if err.Error() != "" {
		er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
	}
	return devices
}

func fetchDevicesById(deviceId string) (openapi.Device, error) {
	apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
	device, resp, err := apiClient.DeviceApi.DevicesDidGet(ctx, deviceId).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		err := RefreshToken()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
		device, resp, err = apiClient.DeviceApi.DevicesDidGet(ctx, deviceId).Execute()
		if resp.StatusCode == 404 {
			return openapi.Device{}, fmt.Errorf("No device found with deviceID: %s", deviceId)
		}
		if err.Error() != "" {
			fmt.Printf("Unable to refresh access token. Please run `login` command again.")
		}
	} else if resp.StatusCode == 404 {
		return openapi.Device{}, fmt.Errorf("No device found with deviceID: %s", deviceId)
	} else if resp.StatusCode == 401 {
		fmt.Printf("Unable to refresh access token. Please run `login` command again.")
	} else if err.Error() != "" {
		er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
	}
	return device, nil
}
