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

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
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

func login(cmd *cobra.Command, args []string) {
	app, err := public.New("config.ClientID")
	if err != nil {
		panic(err)
	}

	// var userAccount shared.Account
	accounts := app.Accounts()
	for _, account := range accounts {
		if account.PreferredUsername == "" {

		}
		
				// if account.PreferredUsername == config.Username {
		// 	userAccount = account
		// }
	}

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
