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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type BaseMockObject struct {
	mock.Mock
}

func (m *BaseMockObject) ListDeployments(connectionString string) ([]string, error) {
	return []string{
		"base_deployment_version1",
		"base_deployment_version2",
		"base_deployment_version3",
		"base_deployment_version4",
		"base_deployment",
		"my_fancy_base_deployment",
	}, nil
}

func TestGetBaseDeploymentVersionValid1(t *testing.T) {
	assert := assert.New(t)
	version, err := getBaseDeploymentVersion("base_deployment_0.1.0-205.Branch.feature-refactoring-servers-common.Sha.a0b0b5506c255377f19424127cefdf343591b3bb")
	assert.Nil(err)
	assert.Equal(version, "0.1.0-205.Branch.feature-refactoring-servers-common.Sha.a0b0b5506c255377f19424127cefdf343591b3bb")
}

func TestGetBaseDeploymentVersionValid2(t *testing.T) {
	assert := assert.New(t)
	version, err := getBaseDeploymentVersion("base_deployment_dev")
	assert.Nil(err)
	assert.Equal(version, "dev")
}

func TestGetBaseDeploymentVersionInvalid1(t *testing.T) {
	assert := assert.New(t)
	version, err := getBaseDeploymentVersion("base_deployment")
	assert.NotNil(err)
	assert.Equal(version, "")
}

func TestGetBaseDeploymentVersionInvalid2(t *testing.T) {
	assert := assert.New(t)
	version, err := getBaseDeploymentVersion("test_deployment_version")
	assert.NotNil(err)
	assert.Equal(version, "")
}

func TestGetBaseDeploymentVersionInvalid3(t *testing.T) {
	assert := assert.New(t)
	version, err := getBaseDeploymentVersion("base_deployment_")
	assert.NotNil(err)
	assert.Equal(version, "")
}

func TestGetBaseDeployments(t *testing.T) {
	assert := assert.New(t)
	testObj := new(BaseMockObject)
	d, err := testObj.ListDeployments("")
	assert.Nil(err)
	baseDeployments := getBaseDeployments(d)
	assert.Contains(baseDeployments, "base_deployment_version1")
	assert.Contains(baseDeployments, "base_deployment_version2")
	assert.Contains(baseDeployments, "base_deployment_version3")
	assert.Contains(baseDeployments, "base_deployment_version4")
	assert.NotContains(baseDeployments, "base_deployment")
	assert.NotContains(baseDeployments, "my_fancy_base_deployment")
}
