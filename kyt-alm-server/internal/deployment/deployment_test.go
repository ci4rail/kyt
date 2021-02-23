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

package deployment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testManifest struct {
	Key string `json:"key"`
}

func newCustomerDeploymentForTest() *Deployment {
	if err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="); err != nil {
		fmt.Println(err)
		return nil
	}
	m := &testManifest{
		Key: "myValue",
	}
	j, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d, err := NewDeployment(string(j), "mydeployment", "tag='myTag'", true, fmt.Sprintf("%d", time.Now().Unix()), 123)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	os.Unsetenv("IOTHUB_SERVICE_CONNECTION_STRING")
	return d
}

func newBaseDeploymentForTest() *Deployment {
	if err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="); err != nil {
		fmt.Println(err)
		return nil
	}
	m := &testManifest{
		Key: "myValue",
	}
	j, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d, err := NewDeployment(string(j), "mydeployment", "tag='myTag'", false, fmt.Sprintf("%d", time.Now().Unix()), 123)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	os.Unsetenv("IOTHUB_SERVICE_CONNECTION_STRING")
	return d
}

func TestNewCustomerDeployment(t *testing.T) {
	assert := assert.New(t)
	d := newCustomerDeploymentForTest()
	assert.Equal(d.connectionString, "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=", "They should be equal")
	assert.Equal(d.name, "mydeployment")
	assert.Equal(d.hubName, "myHub")
	assert.Equal(d.priority, 223)
	assert.Equal(d.layered, true)
	assert.Equal(d.targetCondition, "tag='myTag'")
	var u testManifest
	err := json.Unmarshal([]byte(d.manifest), &u)
	assert.Nil(err, "should be nil")
	assert.Equal(u.Key, "myValue", "they should be equal")
}

func TestNewBaseDeployment(t *testing.T) {
	assert := assert.New(t)
	d := newBaseDeploymentForTest()
	assert.Equal(d.connectionString, "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=", "They should be equal")
	assert.Equal(d.name, "mydeployment")
	assert.Equal(d.hubName, "myHub")
	assert.Equal(d.priority, 123)
	assert.Equal(d.layered, false)
	assert.Equal(d.targetCondition, "tag='myTag'")
	var u testManifest
	err := json.Unmarshal([]byte(d.manifest), &u)
	assert.Nil(err, "should be nil")
	assert.Equal(u.Key, "myValue", "they should be equal")
}

func TestNewDeploymentNoEnv(t *testing.T) {
	assert := assert.New(t)

	m := &testManifest{
		Key: "myValue",
	}
	j, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	_, err = NewDeployment(string(j), "mydeployment", "tag='myTag'", true, fmt.Sprintf("%d", time.Now().Unix()), 1)
	assert.NotNil(err, "should not be nil")
}

func TestCreateManifestFile(t *testing.T) {
	assert := assert.New(t)

	d := newCustomerDeploymentForTest()
	tmpFile, err := d.createManifestFile()
	assert.Nil(err)
	readFile, err := ioutil.ReadFile(tmpFile.Name())
	assert.Nil(err)
	var u testManifest
	err = json.Unmarshal([]byte(readFile), &u)
	assert.Nil(err)
	assert.Equal(u.Key, "myValue")

	defer os.Remove(tmpFile.Name())
}

func TestApplyDeploymentNoConnectionString(t *testing.T) {
	assert := assert.New(t)
	deploymentStr := "{\"emptyjson\": true}"
	_, err := NewDeployment(deploymentStr, "test_deployment", "deviceId='myDeviceId'", true, fmt.Sprintf("%d", time.Now().Unix()), 1)
	assert.NotNil(err)
}
