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
	DefaultDlmServer = "https://testing.dlm.edgefarm.io/v1"
	// DefaultAlmServer Production URL to kyt-alm-server
	DefaultAlmServer = "https://testing.alm.edgefarm.io/v1"
)

const (
	// TokenEndpoint Token request URL
	TokenEndpoint = "https://edgefarm-dev.eu.auth0.com/oauth/token"
	// ClientID kyt-cli client id
	ClientID = "S0uKZJ00Mj4gRNdj1uA889lvc5Vnh0QI"
)

const (
	// DlmScope is the name of the scope for dlm apis
	DlmScope = "dlm"
	// AlmScope is the name of the scope for alm apis
	AlmScope = "alm"
)

// GetConstScopes returns the scopes for the requested ressource that are configured for the application. At least one scope is needed for a successfull login.
// If no scopes are defined, there will be no token assigned and returns with error code 400.
func GetConstScopes(ressource string) ([]string, error) {
	if ressource == DlmScope {
		return []string{DlmScope + "/DevicesGet.read",
			DlmScope + "/DevicesDidGet.read"}, nil

	} else if ressource == AlmScope {
		return []string{AlmScope + "/RuntimesGet.read",
			AlmScope + "/RuntimesRidGet.read",
			AlmScope + "/Apply.write"}, nil
	} else {
		return nil, fmt.Errorf("scopes for invalid ressource requested")
	}
}
