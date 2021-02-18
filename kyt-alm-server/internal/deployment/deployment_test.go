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

	"github.com/ci4rail/kyt/kyt-alm-server/internal/deployment/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) ListDeployments(connectionString string) ([]string, error) {
	return []string{
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication-1613639624",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613639181",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613640525",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication-1613639181",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613638006",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613640495",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613636260",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613636235",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613636170",
		"c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613636087",
		"myTenant1_application1_1600000000",
		"myTenant1_application1_1600000001",
		"myTenant1_application2_1600000002",
		"myTenant2_application1_1600000003",
		"myTenant2_application1_1600000004",
		"myTenant3_application1_1600000005",
		"myTenant3_application2_1600000006",
		"myTenant3_application2_1600000007",
		"myTenant3_application2_1600000008",
		"myTenant4_application1_1600000009",
		"myTenant4_application1_1600000010",
		"myTenant1_application1_1600000011",
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
	d, err := NewDeployment(string(j), "mydeployment", "tag='myTag'", time.Now().Unix())
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
	assert.Equal(d.name, "mydeployment", "they should be equal")
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
	_, err = NewDeployment(string(j), "mydeployment", "tag='myTag'", time.Now().Unix())
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
	_, err := NewDeployment(deploymentStr, "test_deployment", "deviceId='myDeviceId'", time.Now().Unix())
	assert.NotNil(err)
}

func TestGetLatestDeploymentFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(MyMockedObject)
	d, err := testObj.ListDeployments("")
	assert.Nil(err)
	fmt.Println(d)
	latest, err := getLatestDeployment(d, "c5399437-e3d8-4f26-a011-e2e447815d9c", "myapplication")
	assert.Nil(err)
	assert.Equal(latest, "c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613640525")
}

func TestGetLatestDeploymentNotFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(MyMockedObject)
	d, err := testObj.ListDeployments("")
	assert.Nil(err)
	fmt.Println(d)
	latest, err := getLatestDeployment(d, "myTenant1", "application9")
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
	latest, err := getLatestDeployment(d, "myTenant9", "application1")
	assert.NotNil(err)
	assert.Equal(latest, "")
}

func TestGetTimestampFromDeploymentValid1(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "myTenant1_application1_1613595000"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.Nil(err)
	assert.Equal(timestamp, 1613595000)
}

func TestGetTimestampFromDeploymentValid2(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "application1_1324"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.Nil(err)
	assert.Equal(timestamp, 1324)
}

func TestGetTimestampFromDeploymentValid3(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "_1324"
	timestamp, err := getTimestampFromDeployment(deploymentName)
	assert.Nil(err)
	assert.Equal(timestamp, 1324)
}

func TestGetTimestampFromDeploymentInvalid1(t *testing.T) {
	assert := assert.New(t)
	deploymentName := "application1_adsf"
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
	deploymentName := "application1_"
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

func TestValidateDeploymentNameValid1(t *testing.T) {
	assert := assert.New(t)
	ok := deploymentNameValid("c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613639624")
	assert.True(ok)
}

func TestValidateDeploymentNameValid2(t *testing.T) {
	assert := assert.New(t)
	ok := deploymentNameValid("c5399437e2e447815d9c_myapp-lication_1613639624")
	assert.True(ok)
}

func TestValidateDeploymentNameInvalid1(t *testing.T) {
	assert := assert.New(t)
	ok := deploymentNameValid("c5399437e2e_447815d9c_myapp-lication_1613639624")
	assert.False(ok)
}

func TestValidateDeploymentNameInvalid2(t *testing.T) {
	assert := assert.New(t)
	ok := deploymentNameValid("c5399437e2e_447815d9c_myapp-lication")
	assert.False(ok)
}

func TestValidateDeploymentNameInvalid3(t *testing.T) {
	assert := assert.New(t)
	ok := deploymentNameValid("c5399437e2e_447815d9c_myapplication-1613639624")
	assert.False(ok)
}

func TestCreateFromCustomerDeployment(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=kyt-dev-iot-hub.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=7m+i8rSSQCyIJGjdBVcFjmjdBOxVPBcbb34iFrxeEcA=")
	assert.Nil(err)
	c := &manifest.CustomerManifest{
		Application: "myApplication",
		Modules: []manifest.ModuleType{
			{
				Name:            "module1",
				Image:           "image1",
				CreateOptions:   "{}",
				ImagePullPolicy: "image_pull_1",
				RestartPolicy:   "policy1",
				Status:          "status1",
				StartupOrder:    1,
			},
			{
				Name:            "module2",
				Image:           "image2",
				CreateOptions:   "{}",
				ImagePullPolicy: "image_pull_2",
				RestartPolicy:   "policy2",
				Status:          "status2",
				StartupOrder:    2,
			},
		},
	}

	ok, err := createFromCustomerDeployment("c5399437-e3d8-4f26-a011-e2e447815d9c", c)
	assert.Nil(err)
	assert.True(ok)
}
