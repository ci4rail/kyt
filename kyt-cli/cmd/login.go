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

	"github.com/ci4rail/kyt/kyt-cli/internal/configuration"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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
	req, err := createAccessTokenRequest(configuration.TokenEndpoint, configuration.ClientId, username, password)
	if err != nil {
		er(err)
	}
	resp, err := sendAccessTokenRequest(req)
	if err != nil {
		er(err)
	}
	token, refreshToken, err := extractAccessToken(resp)
	if err != nil {
		er(err)
	}
	claims, err := getTokenClaims(token)
	if err != nil {
		er(err)
	}
	writeTokensToConfig(token, refreshToken)

	fmt.Println("Login Succeeded")
	givenName := ""
	if givenNameClaims, ok := claims["given_name"]; ok {
		if str, ok := givenNameClaims.(string); ok {
			givenName = str
		}

	}
	familyName := ""
	if familyNameClaims, ok := claims["family_name"]; ok {
		if str, ok := familyNameClaims.(string); ok {
			familyName = str
		}
	}
	name := fmt.Sprintf("%s %s", givenName, familyName)
	if len(name) > 0 {
		fmt.Printf("Logged in as: %s\n", name)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
