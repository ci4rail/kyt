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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CustomerMockObject struct {
	mock.Mock
}

func (m *CustomerMockObject) ListDeployments(connectionString string) ([]string, error) {
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

func TestGetLatestCustomerDeploymentFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(CustomerMockObject)
	d, err := testObj.ListDeployments("")
	assert.Nil(err)
	fmt.Println(d)
	latest, err := getLatestCustomerDeployment(d, "c5399437-e3d8-4f26-a011-e2e447815d9c", "myapplication")
	assert.Nil(err)
	assert.Equal(latest, "c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613640525")
}

func TestGetLatestCustomerDeploymentNotFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(CustomerMockObject)
	d, err := testObj.ListDeployments("")
	assert.Nil(err)
	fmt.Println(d)
	latest, err := getLatestCustomerDeployment(d, "myTenant1", "application9")
	assert.NotNil(err)
	assert.Equal(latest, "")
}

func TestListDeploymentsTenantNotFound(t *testing.T) {
	assert := assert.New(t)
	testObj := new(CustomerMockObject)
	connectionString := "HostName=HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="
	d, err := testObj.ListDeployments(connectionString)
	assert.Nil(err)
	fmt.Println(d)
	latest, err := getLatestCustomerDeployment(d, "myTenant9", "application1")
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

func TestValidatecustomerDeploymentNameValid1(t *testing.T) {
	assert := assert.New(t)
	ok := customerDeploymentNameValid("c5399437-e3d8-4f26-a011-e2e447815d9c_myapplication_1613639624")
	assert.True(ok)
}

func TestValidatecustomerDeploymentNameValid2(t *testing.T) {
	assert := assert.New(t)
	ok := customerDeploymentNameValid("c5399437e2e447815d9c_myapp-lication_1613639624")
	assert.True(ok)
}

func TestValidateDeploymentNameInvalid1(t *testing.T) {
	assert := assert.New(t)
	ok := customerDeploymentNameValid("c5399437e2e_447815d9c_myapp-lication_1613639624")
	assert.False(ok)
}

func TestValidateDeploymentNameInvalid2(t *testing.T) {
	assert := assert.New(t)
	ok := customerDeploymentNameValid("c5399437e2e_447815d9c_myapp-lication")
	assert.False(ok)
}

func TestValidateDeploymentNameInvalid3(t *testing.T) {
	assert := assert.New(t)
	ok := customerDeploymentNameValid("c5399437e2e_447815d9c_myapplication-1613639624")
	assert.False(ok)
}
