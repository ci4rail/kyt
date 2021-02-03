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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:     "devices",
	Aliases: []string{"device", "dev"},
	Short:   "Display all kyt-devices",
	Long: `Display all kyt-devices

Prints a table of the most important information of all kyt-devices.
`,
	Run: getDevices,
}

func getDevicesAll() []openapi.Device {
	apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
	devices, resp, err := apiClient.DeviceApi.DevicesGet(ctx).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		RefreshToken()

		apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
		devices, _, err = apiClient.DeviceApi.DevicesGet(ctx).Execute()
		if err.Error() != "" {
			er(fmt.Sprintf("Error calling RefreshApi.RefreshToken: %v\n", err))
		}
	} else if err.Error() != "" {
		er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
	}
	return devices
}

func getDevicesById(deviceId string) (openapi.Device, error) {
	apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
	device, resp, err := apiClient.DeviceApi.DevicesDidGet(ctx, deviceId).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		RefreshToken()

		apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
		device, _, err = apiClient.DeviceApi.DevicesDidGet(ctx, deviceId).Execute()
		if resp.StatusCode == 404 {
			return openapi.Device{}, fmt.Errorf("No device found with deviceID: %s", deviceId)
		}
		if err.Error() != "" {
			fmt.Printf("Error calling RefreshApi.RefreshToken: %v\n", err)
		}
	} else if resp.StatusCode == 404 {
		fmt.Println("No device found with deviceID: ", deviceId)
	} else if err.Error() != "" {
		er(fmt.Sprintf("Error calling DeviceApi.DevicesGet: %v\n", err))
	}
	return device, nil
}

func getDevices(cmd *cobra.Command, args []string) {
	if !viper.IsSet("token") {
		fmt.Println("No access token set. Please run `login` command.")
		os.Exit(1)
	}
	var devices []openapi.Device
	if len(args) > 0 {
		for _, arg := range args {
			dev, err := getDevicesById(arg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			devices = append(devices, dev)
		}
	} else {
		devices = getDevicesAll()
	}

	fmt.Println("DEVICE ID")
	for _, dev := range devices {
		fmt.Println(dev.GetId())
	}
}

func init() {
	getCmd.AddCommand(devicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
