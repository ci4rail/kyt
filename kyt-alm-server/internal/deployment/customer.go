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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ci4rail/kyt/kyt-alm-server/internal/deployment/manifest"
	iothub "github.com/ci4rail/kyt/kyt-server-common/iothubwrapper"
)

// CreateOrUpdateFromCustomerDeployment creates a new deployment and deletes the old one if it was
// already present.
func CreateOrUpdateFromCustomerDeployment(tenantID string, c manifest.CustomerManifest) error {
	cs, err := iothub.ReadConnectionStringFromEnv()
	if err != nil {
		return err
	}
	// Get all deployments fron backend service
	deployments, err := ListDeployments(cs)
	if err != nil {
		return err
	}
	// Get latest deployment for specific tenantID
	latestToBeDeletedOnSuccess, err := getLatestCustomerDeployment(deployments, tenantID, c.Application)
	if err != nil {
		fmt.Println(err)
	}
	// The new deployment needs to be created first to start the update process
	_, err = createFromCustomerDeployment(tenantID, c)
	if err != nil {
		return err
	}
	// Delete old deployment with the same name to finish the update process
	if len(latestToBeDeletedOnSuccess) > 0 {
		// create new dummy deployment with specific name to be deleted
		deleteDeployment, err := NewDeployment("{}", latestToBeDeletedOnSuccess, "", true, "0", 0)
		if err != nil {
			return err
		}
		_, err = deleteDeployment.DeleteDeployment()
		if err != nil {
			return err
		}
	}
	return nil
}

// createFromCustomerDeployment creates and applies from a customer deployment
func createFromCustomerDeployment(tenantID string, c manifest.CustomerManifest) (bool, error) {
	now := fmt.Sprintf("%d", time.Now().Unix())
	layered, err := manifest.CreateLayeredManifest(c, tenantID)
	if err != nil {
		return false, err
	}
	deploymentName := fmt.Sprintf("%s_%s", tenantID, c.Application)
	// Currently the target condition is fixed to the tenant's ID
	targetCondition := fmt.Sprintf("tags.alm=true AND tags.tenantId='%s'", tenantID)
	d, err := NewDeployment(layered, deploymentName, targetCondition, true, now, 0)
	if err != nil {
		return false, err
	}
	ok, err := d.applyDeployment()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// getLatestCustomerDeployment gets the last deployment from tenandId with application
func getLatestCustomerDeployment(deployments []string, tenantID string, application string) (string, error) {
	appName := fmt.Sprintf("%s_%s", tenantID, application)
	tenantDeployments := []string{}
	for _, d := range deployments {
		if strings.Contains(d, appName) {
			if customerDeploymentNameValid(d) {
				tenantDeployments = append(tenantDeployments, d)
			}
		}
	}
	sorted := make([]string, len(tenantDeployments))
	copy(sorted, tenantDeployments)
	sort.Sort(ByTimestamp(sorted))

	if len(sorted) > 0 {
		return sorted[0], nil
	}
	return "", fmt.Errorf("Info: no latest deployment found")
}

// customerDeploymentNameValid checks if the deployment name has the right pattern
func customerDeploymentNameValid(name string) bool {
	re := regexp.MustCompile(`^[a-z0-9-]+_[a-z0-9-]+_[0-9]+$`)
	if ok := re.MatchString(name); ok {
		return true
	}
	return false
}

// getTimestampFromDeployment returns the timestamp from the deployment name
func getTimestampFromDeployment(deployment string) (int, error) {
	// This regex pattern finds any numbers with leading `_` in a string.
	// e.g. tenant.application.321553 results in -> `_321553`
	re := regexp.MustCompile(`([_][0-9]+)$`)
	if ok := re.MatchString(deployment); ok {
		matches := re.FindAllString(deployment, -1)
		trimmed := strings.Replace(matches[0], "_", "", -1)
		timestamp, err := strconv.Atoi(trimmed)
		if err != nil {
			return 0, err
		}
		return timestamp, nil
	}
	return 0, fmt.Errorf("Error: no timestamp found")
}
