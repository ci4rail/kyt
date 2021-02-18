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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) ListDeployments(connectionString string) ([]string, error) {
	return []string{
		"myTenant1.application1.1600000000",
		"myTenant1.application1.1600000001",
		"myTenant1.application2.1600000002",
		"myTenant2.application1.1600000003",
		"myTenant2.application1.1600000004",
		"myTenant3.application1.1600000005",
		"myTenant3.application2.1600000006",
		"myTenant3.application2.1600000007",
		"myTenant3.application2.1600000008",
		"myTenant4.application1.1600000009",
		"myTenant4.application1.1600000010",
		"myTenant1.application1.1600000011",
	}, nil
}

type testManifest struct {
	Key string `json:"key"`
}

func newDeploymentForTest() *Deployment {
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
	d, err := NewDeployment(string(j), "myDeployment", "tag='myTag'")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	os.Unsetenv("IOTHUB_SERVICE_CONNECTION_STRING")
	return d
}

func TestNewDeployment(t *testing.T) {
	assert := assert.New(t)
	d := newDeploymentForTest()
	assert.Equal(d.connectionString, "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=", "They should be equal")
	assert.Equal(d.name, "myDeployment", "they should be equal")
	assert.Equal(d.priority, defaultPriority)
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
	_, err = NewDeployment(string(j), "myDeployment", "tag='myTag'")
	assert.NotNil(err, "should not be nil")
}

func TestCreateManifestFile(t *testing.T) {
	assert := assert.New(t)

	d := newDeploymentForTest()
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
	_, err := NewDeployment(deploymentStr, "test_deployment", "deviceId='myDeviceId'")
	assert.NotNil(err)
}

func TestListDeploymentsFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(MyMockedObject)
	connectionString := "HostName=HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	d, err := testObj.ListDeployments(connectionString)
	assert.Nil(err)
	fmt.Println(d)
	latest, err := GetLatestDeployment(d, "myTenant1", "application1")
	assert.Nil(err)
	assert.Equal(latest, "myTenant1.application1.1600000011")
}

func TestListDeploymentsApplicationNotFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(MyMockedObject)
	connectionString := "HostName=HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	d, err := testObj.ListDeployments(connectionString)
	assert.Nil(err)
	fmt.Println(d)
	latest, err := GetLatestDeployment(d, "myTenant1", "application9")
	assert.NotNil(err)
	assert.Equal(latest, "")
}

func TestListDeploymentsTenantNotFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(MyMockedObject)
	connectionString := "HostName=HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	d, err := testObj.ListDeployments(connectionString)
	assert.Nil(err)
	fmt.Println(d)
	latest, err := GetLatestDeployment(d, "myTenant9", "application1")
	assert.NotNil(err)
	assert.Equal(latest, "")
}

func TestGetTimestampFromDeploymentValid1(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "myTenant1.application1.1613595000"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.Nil(err)
	assert.Equal(timestamp, 1613595000)
}

func TestGetTimestampFromDeploymentValid2(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "application1.1324"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.Nil(err)
	assert.Equal(timestamp, 1324)
}

func TestGetTimestampFromDeploymentValid3(t *testing.T) {
	assert := assert.New(t)
	deploymentName := ".1324"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.Nil(err)
	assert.Equal(timestamp, 1324)
}

func TestGetTimestampFromDeploymentInvalid1(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "application1.adsf"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.NotNil(err)
	assert.Equal(timestamp, 0)
}

func TestGetTimestampFromDeploymentInvalid2(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "application1"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.NotNil(err)
	assert.Equal(timestamp, 0)
}

func TestGetTimestampFromDeploymentInvalid3(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "application1."
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.NotNil(err)
	assert.Equal(timestamp, 0)
}

func TestGetTimestampFromDeploymentInvalid4(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "32156"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.NotNil(err)
	assert.Equal(timestamp, 0)
}
