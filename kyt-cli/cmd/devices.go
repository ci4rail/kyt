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
	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"
	openapi "github.com/ci4rail/kyt/kyt-cli/openapidlm"
	"github.com/spf13/viper"
)

// FetchDevices Returns information from passed devices. In case list is empty, infromation
// from all connected devices are returned.
func FetchDevices(list []string) []openapi.Device {
	var devices []openapi.Device
	if len(list) > 0 {
		for _, arg := range list {
			dev, err := fetchDevicesByID(arg)
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
	apiClient, ctx := api.NewDlmAPIWithToken(viper.GetString("dlmServerURL"), viper.GetString("token"))
	devices, resp, err := apiClient.DeviceApi.DevicesGet(ctx).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		err := RefreshToken()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiClient, ctx := api.NewDlmAPIWithToken(viper.GetString("dlmServerURL"), viper.GetString("token"))
		devices, _, err = apiClient.DeviceApi.DevicesGet(ctx).Execute()
		if resp.StatusCode == 403 {
			e.Er("Forbidden\n")
		} else if err.Error() != "" {
			e.Er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
		}
	} else if resp.StatusCode == 403 {
		e.Er("Forbidden\n")
	} else if err.Error() != "" {
		e.Er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
	}
	return devices
}

func fetchDevicesByID(deviceID string) (openapi.Device, error) {
	apiClient, ctx := api.NewDlmAPIWithToken(viper.GetString("dlmServerURL"), viper.GetString("token"))
	device, resp, err := apiClient.DeviceApi.DevicesDidGet(ctx, deviceID).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		err := RefreshToken()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiClient, ctx := api.NewDlmAPIWithToken(viper.GetString("dlmServerURL"), viper.GetString("token"))
		device, resp, err = apiClient.DeviceApi.DevicesDidGet(ctx, deviceID).Execute()
		if resp.StatusCode == 404 {
			return openapi.Device{}, fmt.Errorf("No device found with deviceID: %s", deviceID)
		}
		if err.Error() != "" {
			fmt.Printf("Unable to refresh access token. Please run `login` command again.")
		}
	} else if resp.StatusCode == 404 {
		return openapi.Device{}, fmt.Errorf("No device found with deviceID: %s", deviceID)
	} else if resp.StatusCode == 403 {
		e.Er("Forbidden\n")
	} else if resp.StatusCode == 401 {
		e.Er("Unable to refresh access token. Please run `login` command again.\n")
	} else if err.Error() != "" {
		e.Er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
	}
	return device, nil
}
