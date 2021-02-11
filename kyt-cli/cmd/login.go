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
	// "fmt"

	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Ci4Rail services",
	Long: `Login to Ci4Rail services

Log in with user name and password.

Log in interactively.

Not implemented yet.`,
	Run: login,
}

var (
	username string
	password string
)

func login(cmd *cobra.Command, args []string) {
	prompt := promptui.Prompt{
		Label:    "Username",
		Validate: nil,
	}

	username, err := prompt.Run()
	if err != nil {
		log.Panicln(err)
	}
	prompt = promptui.Prompt{
		Label:    "Password",
		Validate: nil,
		Mask:     ' ',
	}

	password, err = prompt.Run()
	if err != nil {
		log.Panicln(err)
	}
	req, err := createTokenRequest(viper.GetString("token_endpoint"), viper.GetString("client_id"), username, password)
	if err != nil {
		er(err)
	}
	resp, err := sendTokenRequest(req)
	if err != nil {
		er(err)
	}
	token, err := extractAccessToken(resp)
	if err != nil {
		er(err)
	}
	claims, err := getTokenClaims(token)
	if err != nil {
		er(err)
	}
	fmt.Printf("Token: %s\n", token)
	RefreshToken()
	viper.Set("token", token)
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		er(err)
	}
	err = viper.WriteConfigAs(fmt.Sprintf("%s/%s.%s", home, kytCliConfigFile, kytCliConfigFileType))
	if err != nil {
		log.Println("Cannot save config file")
	}
	fmt.Println("Login Succeeded")
	fmt.Printf("Logged in as: %s %s\n", claims["given_name"], claims["family_name"])
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
