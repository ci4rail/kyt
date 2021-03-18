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
	"net/http"
	"os"

	sw "github.com/ci4rail/kyt/kyt-dlm-server/api"
	common "github.com/ci4rail/kyt/kyt-dlm-server/internal/common"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

var serverAddr string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyt-dlm-server",
	Short: "REST API Server to control KYT dlm",
	Long: `kyt dlm api server is the central service to control Ci4Rails KYT device livecycle management.

	KYT consists of application lifecycle management (alm), device lifecycle management (dlm) and application data services (ads).
`,
	Run: func(cmd *cobra.Command, args []string) {

		c := cors.New(cors.Options{
			AllowCredentials: true,
			AllowedHeaders:   []string{"Authorization"},
		})

		router, err := sw.NewRouter()
		if err != nil {
			log.Fatal(err)
		}
		handler := c.Handler(router)
		http.Handle("/", router)
		err = http.ListenAndServe(serverAddr, handler)
		if err != nil {
			log.Fatal(err)
		}
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

	rootCmd.PersistentFlags().StringVar(&serverAddr, "addr", fmt.Sprintf(":%d", common.ServicePort), "address and port the server shall listen to")

}
