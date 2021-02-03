/*
Copyright © 2021 Ci4Rail GmbH

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

package controllerif

// NewIOTHubServiceClient points to the actual creator function
type NewIOTHubServiceClient func(connectionString string) (IOTHubServices, error)

// IOTHubServices has all the services the IOTHub Controller offers
type IOTHubServices interface {
	ListDeviceIDs() (*[]string, error)
	ListDeviceById(string) (*string, error)
}
