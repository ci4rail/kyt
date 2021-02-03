package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	api "github.com/ci4rail/kyt/kyt-cli/internal/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func RefreshToken() {
	apiClient, ctx := api.NewAPIWithToken(serverURL, viper.GetString("token"))
	fmt.Println("Token expired. Refreshing...")
	resp, openapierr := apiClient.AuthApi.AuthRefreshTokenGet(ctx).Execute()
	if openapierr.Error() != "" {
		er(fmt.Sprintf("Error calling RefreshApi.RefreshToken: %v\n", openapierr))
	}

	var data map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&data)
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
}
