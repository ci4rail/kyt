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

const (
	TokenEndpoint = "https://ci4railtesting.b2clogin.com/ci4railtesting.onmicrosoft.com/B2C_1_signin_native/oauth2/v2.0/token"
	ClientId      = "2c9a4ac6-c0ad-4bd4-bc3d-544ff94a2471"
)

// GetConstScopes returns the scopes that are configured for the application. At least one scope is needed for a successfull login.
// If no scopes are defined, there will be no token assigned and returns with error code 400.
func GetConstScopes() []string {
	return []string{"https://ci4railtesting.onmicrosoft.com/794d32c1-8515-4daf-be13-4c914593bbfc/DevicesGet.read",
		"https://ci4railtesting.onmicrosoft.com/794d32c1-8515-4daf-be13-4c914593bbfc/DevicesDidGet.read"}
}
