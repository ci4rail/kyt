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

var (
	output string
)

// devicesCmd represents the devices command
var describeDevicesCmd = &cobra.Command{
	Use:     "devices",
	Aliases: []string{"device", "dev"},
	Short:   "Display detailed information about kyt-devices",
	Long: `Display detailed information about kyt-devices

Prints detailed information about kyt-devices.
`,
	Run: describeDevices,
}

func describeDevices(cmd *cobra.Command, args []string) {
	tokenConfigCheck()

	devices := fetchDevices(args)

	if len(devices) > 0 {
		y, err := convertToYaml(&devices)
		if err != nil {
			er(err)
		}
		fmt.Println(y)
	}
}

func init() {
	describeCmd.AddCommand(describeDevicesCmd)
}
