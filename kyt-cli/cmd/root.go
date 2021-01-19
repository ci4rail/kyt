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

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	DEFAULT_KYT_SERVER = "https://kyt.ci4rail.com/v1"
)

var (
	cfgFile            string
	serverURL          string
	serverURLParameter string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyt",
	Short: "kyt cli",
	Long: `kyt cli controls kyt-services

Control the kyt-servies application lifecycle management (alm), device lifecycle management (dlm) and application data services (ads).`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

func er(msg interface{}) {
	fmt.Fprintf(os.Stderr, "Error: %v", msg)
	os.Exit(1)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		er(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kyt-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&serverURLParameter, "server", DEFAULT_KYT_SERVER, "use alternative server URL")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".kyt-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kyt-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())

		// override default server config with config file
		// priority 1: '--server' argument that differs from defailt
		// priority 2: 'server' from config file
		// priority 3: default server
		if serverURLParameter != DEFAULT_KYT_SERVER {
			serverURL = serverURLParameter
		} else {
			serverConfig := viper.GetString("server")
			if len(serverConfig) > 0 {
				serverURL = serverConfig
			} else {
				serverURL = DEFAULT_KYT_SERVER
			}
		}
	}
}
