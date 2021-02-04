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
var getDevicesCmd = &cobra.Command{
	Use:     "devices",
	Aliases: []string{"device", "dev"},
	Short:   "Display kyt-devices",
	Long: `Display kyt-devices

Prints a table of the most important information of kyt-devices.
`,
	Run: getDevices,
}

func getDevices(cmd *cobra.Command, args []string) {
	tokenConfigCheck()

	devices := fetchDevices(args)

	if len(devices) > 0 {
		switch o := output; o {
		case "json":
			j, err := convertToJson(&devices)
			if err != nil {
				er(err)
			}
			fmt.Println(j)
		case "yaml":
			y, err := convertToYaml(&devices)
			if err != nil {
				er(err)
			}
			fmt.Println(y)
		case "wide":
			// wide: Add here some more information for the table
			fmt.Printf("%-40s\t%-16s  %s\n", "DEVICE ID", "CONNECTION STATE", "FIRMWARE VERSION")
			for _, dev := range devices {
				fmt.Printf("%-40s\t%-16s  %s\n", dev.GetId(), dev.GetNetwork(), dev.GetFirmwareVersion())
			}
		default:
			// short: only the most important information
			fmt.Printf("%-40s\t%-16s\n", "DEVICE ID", "CONNECTION STATE")
			for _, dev := range devices {
				fmt.Printf("%-40s\t%s\n", dev.GetId(), dev.GetNetwork())
			}
		}
	}
}

func init() {
	getCmd.AddCommand(getDevicesCmd)
	getDevicesCmd.Flags().StringVarP(&output, "output", "o", "short", "Output format. One of: short|json|yaml|wide")
}
