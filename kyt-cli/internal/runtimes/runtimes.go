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

package runtimes

import (
	"fmt"
	"os"

	api "github.com/ci4rail/kyt/kyt-cli/internal/api"
	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"
	"github.com/ci4rail/kyt/kyt-cli/internal/token"
	openapi "github.com/ci4rail/kyt/kyt-cli/openapialm"
	"github.com/spf13/viper"
)

// FetchRuntimes Returns information from passed runtimes. In case list is empty, infromation
// from all connected runtimes are returned.
func FetchRuntimes(list []string) []openapi.Runtime {
	var runtimes []openapi.Runtime
	if len(list) > 0 {
		for _, arg := range list {
			dev, err := fetchRuntimesByID(arg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			runtimes = append(runtimes, dev)
		}
	} else {
		runtimes = fetchRuntimesAll()
	}
	return runtimes
}

func fetchRuntimesAll() []openapi.Runtime {
	apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("almServerURL"), viper.GetString("alm_token"))
	runtimes, resp, err := apiClient.AlmApi.RuntimesGet(ctx).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		err := token.RefreshToken("alm")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("almServerURL"), viper.GetString("alm_token"))
		runtimes, _, err = apiClient.AlmApi.RuntimesGet(ctx).Execute()
		if err.Error() != "" {
			e.Er(fmt.Sprintf("Error calling AlmApi.RuntimesGet: %v\n", err))
		}
	} else if err.Error() != "" {
		e.Er(fmt.Sprintf("Error calling AlmApi.RuntimesGet: %v\n", err))
	}
	return runtimes
}

func fetchRuntimesByID(runtimeID string) (openapi.Runtime, error) {
	apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("almServerURL"), viper.GetString("alm_token"))
	runtime, resp, err := apiClient.AlmApi.RuntimesRidGet(ctx, runtimeID).Execute()
	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
	if resp.StatusCode == 401 {
		err := token.RefreshToken("alm")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("almServerURL"), viper.GetString("alm_token"))
		runtime, resp, err = apiClient.AlmApi.RuntimesRidGet(ctx, runtimeID).Execute()
		if resp.StatusCode == 404 {
			return openapi.Runtime{}, fmt.Errorf("No runtime found with runtimeID: %s", runtimeID)
		}
		if err.Error() != "" {
			fmt.Printf("Unable to refresh access token. Please run `login` command again.")
		}
	} else if resp.StatusCode == 404 {
		return openapi.Runtime{}, fmt.Errorf("No runtime found with runtimeID: %s", runtimeID)
	} else if resp.StatusCode == 401 {
		fmt.Printf("Unable to refresh access token. Please run `login` command again.")
	} else if err.Error() != "" {
		e.Er(fmt.Sprintf("Error calling AlmApi.RuntimesGet: %v\n", err))
	}
	return runtime, nil
}
