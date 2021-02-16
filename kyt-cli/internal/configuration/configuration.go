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

package configuration

import "fmt"

const (
	// DefaultDlmServer Production URL to kyt-dlm-server
	DefaultDlmServer = "https://dlm.ci4rail.com/v1"
	// DefaultAlmServer Production URL to kyt-alm-server
	DefaultAlmServer = "https://alm.ci4rail.com/v1"
)

const (
	// TokenEndpoint Token request URL
	TokenEndpoint = "https://ci4railtesting.b2clogin.com/ci4railtesting.onmicrosoft.com/B2C_1_signin_native/oauth2/v2.0/token"
	// ClientID kyt-cli client id
	ClientID = "2c9a4ac6-c0ad-4bd4-bc3d-544ff94a2471"
)

// GetConstScopes returns the scopes for the requested ressource that are configured for the application. At least one scope is needed for a successfull login.
// If no scopes are defined, there will be no token assigned and returns with error code 400.
func GetConstScopes(ressource string) ([]string, error) {
	if ressource == "dlm" {
		return []string{"https://ci4railtesting.onmicrosoft.com/794d32c1-8515-4daf-be13-4c914593bbfc/DevicesGet.read",
			"https://ci4railtesting.onmicrosoft.com/794d32c1-8515-4daf-be13-4c914593bbfc/DevicesDidGet.read"}, nil

	} else if ressource == "alm" {
		return []string{"https://ci4railtesting.onmicrosoft.com/0c2e91da-821c-4d38-aa8c-02ee539cdd3e/RuntimesGet.read",
			"https://ci4railtesting.onmicrosoft.com/0c2e91da-821c-4d38-aa8c-02ee539cdd3e/RuntimesRidGet.read"}, nil
	} else {
		return nil, fmt.Errorf("scopes for invalid ressource requested")
	}
}
