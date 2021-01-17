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

	sw "github.com/ci4rail/kyt-cli/kyt-api-server/go"
	"github.com/spf13/cobra"
)

var serverAddr string = ":8080"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyt-api-server",
	Short: "REST API Server to control KYT",
	Long: `kyt api server is the central service to control Ci4Rails KYT.

	KYT consists of application lifecycle management (alm), device lifecycle management (dlm) and application data services (ads).
`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Printf("Run args %s\n", args)
		router := sw.NewRouter()

		log.Fatal(router.Run(serverAddr))

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

	rootCmd.PersistentFlags().StringVar(&serverAddr, "addr", "", "address the server shall listen to (default is \":8080\")")

}
