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

package apply

import (
	"fmt"
	"os"

	"github.com/ci4rail/kyt/kyt-cli/internal/api"
	"github.com/ci4rail/kyt/kyt-cli/internal/configuration"
	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"
	"github.com/ci4rail/kyt/kyt-cli/internal/token"
	openapi "github.com/ci4rail/kyt/kyt-cli/openapialm"
	"github.com/spf13/viper"
)

// CustomerManifest -
func CustomerManifest(c openapi.CustomerManifest) {
	apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("alm_server_url"), viper.GetString("alm_token"))
	resp, err := apiClient.DeploymentApi.ApplyPut(ctx).CustomerManifest(c).Execute()
	if resp == nil && len(err.Error()) > 0 {
		e.Er("DLM Server unreachable\n")
	}
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		err := token.RefreshToken(configuration.AlmScope)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("alm_server_url"), viper.GetString("alm_token"))
		resp, err = apiClient.DeploymentApi.ApplyPut(ctx).CustomerManifest(c).Execute()
		if resp.StatusCode == 401 {
			e.Er("Unable to refresh access token. Please run `login` command again.\n")
		} else if resp.StatusCode == 403 {
			e.Er("Forbidden\n")
		} else if resp.StatusCode == 500 {
			e.Er("Internal server error\n")
		} else if err.Error() != "" {
			e.Er(fmt.Sprintf("Error calling DeploymentApi.ApplyPut: %v\n", err))
		}
	} else if resp.StatusCode == 403 {
		e.Er("Forbidden\n")
	} else if err.Error() != "" {
		e.Er(fmt.Sprintf("Error calling DeploymentApi.ApplyPut: %v\n", err))
	}
}
