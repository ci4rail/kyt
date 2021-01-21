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

	"context"
	"encoding/json"
	"fmt"
	"log"

	api "github.com/ci4rail/kyt-cli/kyt-cli/internal/api"
	"github.com/ci4rail/kyt-cli/kyt-cli/openapi"
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

	password, err := prompt.Run()
	if err != nil {
		log.Panicln(err)
	}
	inlineObject := &openapi.InlineObject{
		Username: &username,
		Password: &password,
	}
	apiClient, _ := api.NewAPI(serverURL)
	resp, openapierr := apiClient.DefaultApi.LoginPost(context.Background()).InlineObject(*inlineObject).Execute()
	if openapierr.Error() != "" {
		log.Fatalf("Error when calling `DefaultApi.LoginPost``: %v\n", openapierr)

	}
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatalf("Error: %e", err)
	}

	token := data["token"]
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
