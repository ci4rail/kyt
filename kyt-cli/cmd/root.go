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
	"path/filepath"

	common "github.com/ci4rail/kyt/kyt-cli/internal/common"
	configuration "github.com/ci4rail/kyt/kyt-cli/internal/configuration"
	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	viperConfig map[string]interface{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyt",
	Short: "kyt cli",
	Long: `kyt cli controls kyt-services

Control the kyt-servies application lifecycle management (alm), device lifecycle management (dlm) and application data services (ads).`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e.Er(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kyt-cli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// override default server config with config file
	// priority 1: servers from config file
	// priority 2: default servers
	// If a config file is found, read it in.
	viper.SetDefault("dlm_server_url", configuration.DefaultDlmServer)
	viper.SetDefault("alm_server_url", configuration.DefaultAlmServer)

	viper.SetConfigType(common.KytCliConfigFileType)
	if cfgFile != "" {
		// Use config file from the flag.
		abs, err := filepath.Abs(cfgFile)
		if err != nil {
			e.Er(err)
		}
		common.KytConfigPath = abs
		viper.SetConfigFile(cfgFile)
		if err := checkAndCreateConfig(cfgFile); err != nil {
			e.Er(err)
		}
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			e.Er(err)
		}
		// Search config in home directory with name ".kyt-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(common.KytCliConfigFile)
		cfgFile := fmt.Sprintf("%s/%s.%s", home, common.KytCliConfigFile, common.KytCliConfigFileType)
		common.KytConfigPath = cfgFile
		if err := checkAndCreateConfig(cfgFile); err != nil {
			e.Er(err)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		e.Er(err)
	}
}

func checkAndCreateConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = viper.WriteConfigAs(path)
		if err != nil {
			e.Er(err)
		}
	}
	return nil
}
