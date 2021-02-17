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
	"testing"
)

func TestIotHubNameFromConnecetionStringValid(t *testing.T) {
	connectionString := "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	cs, err := IotHubNameFromConnecetionString(connectionString)
	if err != nil {
		t.Error(err)
	}
	if cs != "myHub.azure-devices.net" {
		t.Errorf("expected 'myHub.azure-devices.net', got '%s'", cs)
	}
}

func TestIotHubNameFromConnecetionStringInvalid(t *testing.T) {
	connectionString := "SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	cs, err := IotHubNameFromConnecetionString(connectionString)
	if err != nil {
		t.Error(err)
	}
	if cs != "myHub.azure-devices.net" {
		t.Errorf("expected 'myHub.azure-devices.net', got '%s'", cs)
	}
}
