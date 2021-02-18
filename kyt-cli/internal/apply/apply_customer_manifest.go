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
	openapi "github.com/ci4rail/kyt/kyt-cli/openapialm"
)

// CustomerManifest -
func CustomerManifest(*openapi.CustomerManifest) error {

	return nil
}

// func apply_file() {
// 	apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("almServerURL"), viper.GetString("alm_token"))
// 	runtimes, resp, err := apiClient.AlmApi.
// 		RuntimesGet(ctx).Execute()
// 	// 401 mean 'Unauthorized'. Let's try to refresh the token once.
// 	if resp.StatusCode == 401 {
// 		err := token.RefreshToken("alm")
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		apiClient, ctx := api.NewAlmAPIWithToken(viper.GetString("almServerURL"), viper.GetString("alm_token"))
// 		runtimes, _, err = apiClient.AlmApi.RuntimesGet(ctx).Execute()
// 		if err.Error() != "" {
// 			e.Er(fmt.Sprintf("Error calling AlmApi.RuntimesGet: %v\n", err))
// 		}
// 	} else if err.Error() != "" {
// 		e.Er(fmt.Sprintf("Error calling AlmApi.RuntimesGet: %v\n", err))
// 	}
// 	return runtimes, nil
// }
