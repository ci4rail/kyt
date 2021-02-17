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

package controller

import (
	"fmt"
	"log"
	"os"

	iothub "github.com/amenzhinsky/iothub/common"
)

// MapTenantToIOTHubSAS returns the SAS token of the IOT Hub for the specified tenant
// TODO: Either take the SAS from a DB or get it via "az iot hub connection-string show"
// TODO: tenant is currently ignored
func MapTenantToIOTHubSAS(tenant string) (string, error) {
	return ReadConnectionStringFromEnv()
}

func ReadConnectionStringFromEnv() (string, error) {
	envName := fmt.Sprintf("IOTHUB_SERVICE_CONNECTION_STRING")
	val, ok := os.LookupEnv(envName)

	if !ok {
		return "", fmt.Errorf("IOTHUB_SERVICE_CONNECTION_STRING not set")
	}
	return val, nil
}

func IotHubNameFromConnecetionString(connectionString string) (string, error) {
	csMap, err := iothub.ParseConnectionString(connectionString)
	if err != nil {
		log.Panicln(err)
	}
	if value, ok := csMap["HostName"]; ok {
		return value, nil
	}
	return "", fmt.Errorf("Error: 'HostName' not found in connection string")
}
