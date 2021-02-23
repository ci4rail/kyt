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
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ci4rail/kyt/kyt-alm-server/internal/version"
	iothub "github.com/ci4rail/kyt/kyt-server-common/iothub_wrapper"
)

// This is not a comment! It is a directive that includes the json file in a string.
//go:embed assets/base_manifest.json
var baseDeploymentManifest string

// CreateOrUpdateBaseDeployment creates a new base deployment and deletes the old one if it was
// already present.
func CreateOrUpdateBaseDeployment() (bool, error) {
	cs, err := iothub.ReadConnectionStringFromEnv()
	if err != nil {
		return false, err
	}
	// Get all deployments fron backend service
	deployments, err := ListDeployments(cs)
	if err != nil {
		return false, err
	}
	// Get old base deployments
	baseDeploymentsToBeDeleted := getBaseDeployments(deployments)
	// The new deployment needs to be created first to start the update process
	_, err = createBaseDeployment(baseDeploymentManifest)
	if err != nil {
		return false, err
	}
	// Delete old base deployments to finish the update process
	for _, deleteName := range baseDeploymentsToBeDeleted {
		// create new dummy deployment with specific name to be deleted
		deleteDeployment, err := NewDeployment("{}", deleteName, "", false, "0")
		if err != nil {
			return false, err
		}
		_, err = deleteDeployment.DeleteDeployment()
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// createFromCustomerDeployment creates and applies from a customer deployment
func createBaseDeployment(manifest string) (bool, error) {
	now := fmt.Sprintf("%d", time.Now().Unix())
	deploymentName := fmt.Sprintf("%s_%s", baseDeploymentName, version.Version)
	// Currently the target condition is fixed to the tenant's ID
	targetCondition := "tags.alm=true"
	d, err := NewDeployment(manifest, deploymentName, targetCondition, false, now)
	if err != nil {
		return false, err
	}
	ok, err := d.applyDeployment()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// getLatestBaseDeployment gets the last base deployment
func getBaseDeployments(deployments []string) []string {
	baseDeployments := []string{}
	for _, d := range deployments {
		if strings.Contains(d, baseDeploymentName+"_") {
			if baseDeploymentNameValid(d) {
				baseDeployments = append(baseDeployments, d)
			}
		}
	}
	return baseDeployments
}

// baseDeploymentNameValid checks if the base deployment name has the right pattern
func baseDeploymentNameValid(name string) bool {
	re := regexp.MustCompile(`^base_deployment_(.+)$`)
	if ok := re.MatchString(name); ok {
		return true
	}
	return false
}

// baseDeploymentNameValid checks if the deployment name has the right pattern
func getBaseDeploymentVersion(name string) (string, error) {
	re := regexp.MustCompile(`^base_deployment_(.+)$`)
	res := re.FindAllStringSubmatch(name, -1)
	if len(res) > 0 {
		fmt.Printf("version: %s\n", res[0][1])
		if len(res[0]) > 0 {
			return res[0][1], nil
		}
	}
	return "", fmt.Errorf("No version found in base deployment")
}
