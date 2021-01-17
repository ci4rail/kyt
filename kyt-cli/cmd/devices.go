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

	"github.com/spf13/cobra"
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

func getDevices(cmd *cobra.Command, args []string) {
	apiClient, ctx := NewAPI()
	resp, _, err := apiClient.DeviceApi.GetDevices(ctx).Execute()
	if err.Error() != "" {
		er(fmt.Sprintf("Error calling DeviceApi.GetDevices: %v\n", err))
	}
	// response from `GetDevices`: []Device
	fmt.Println("DEVICE ID")
	for _, dev := range resp {
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