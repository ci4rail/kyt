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

package iothubwrapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIotHubNameFromConnecetionStringValid(t *testing.T) {
	assert := assert.New(t)
	connectionString := "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	cs, err := IotHubNameFromConnecetionString(connectionString)
	assert.Nil(err)
	assert.Equal(cs, "myHub")
}

func TestIotHubNameFromConnecetionStringInvalid(t *testing.T) {
	assert := assert.New(t)
	connectionString := "SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	cs, err := IotHubNameFromConnecetionString(connectionString)
	assert.NotNil(err)
	assert.Empty(cs)
}

func TestSplitSubdomainValid(t *testing.T) {
	assert := assert.New(t)
	sub := splitSubdomain("mysubdomain.example.com")
	assert.Equal(sub, "mysubdomain")
}

func TestSplitSubdomainInvalid(t *testing.T) {
	assert := assert.New(t)
	sub := splitSubdomain("mysubdomain")
	assert.Equal(sub, "")
}

func TestReadConnectionStringFromEnvValid(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=")
	assert.Nil(err)
	c, err := ReadConnectionStringFromEnv()
	assert.Nil(err)
	assert.Equal(c, "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=")
}

func TestReadConnectionStringFromEnvInvalid(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=")
	assert.Nil(err)
	err = os.Unsetenv("IOTHUB_SERVICE_CONNECTION_STRING")
	assert.Nil(err)
	c, err := ReadConnectionStringFromEnv()
	assert.NotNil(err)
	assert.Equal(c, "")
}

func TestMapTenantToIOTHubSASValid(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=")
	assert.Nil(err)
	s, err := MapTenantToIOTHubSAS("")
	assert.Nil(err)
	assert.Equal(s, "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=")
}

func TestMapTenantToIOTHubSASInvalid(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=")
	assert.Nil(err)
	err = os.Unsetenv("IOTHUB_SERVICE_CONNECTION_STRING")
	assert.Nil(err)
	s, err := MapTenantToIOTHubSAS("")
	assert.NotNil(err)
	assert.Equal(s, "")
}
