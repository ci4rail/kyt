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

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ci4rail/kyt/kyt-dlm-devinfo-static/fwinfo"
	"github.com/ci4rail/kyt/kyt-dlm-devinfo-static/moduleclient"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyt-dlm-devinfo-static",
	Short: "Server to feed IoTHub Device twin with static device information",
	Long: `kyt-dlm-devinfo-static provides static device information to IoTHub Device Twin, such as
	- Firmware version
	- Inventory data
`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Println("kyt-dlm-devinfo-static RUN")
		c, err := moduleclient.NewModule()
		if err != nil {
			log.Println("failed to create module client", err)
			return
		}
		log.Println("client created")

		fwinfo, err := fwinfo.Read()
		if err != nil {
			log.Println("failed to read fwinfo", err)
			return
		}

		d := moduleclient.DeviceInfo{
			"firmwareversion": fwinfo,
		}
		err = moduleclient.SetStaticDeviceInfo(c, d)
		if err != nil {
			log.Println("failed to set device info", err)
			return
		}
		log.Println("set device info ok", err)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

}
