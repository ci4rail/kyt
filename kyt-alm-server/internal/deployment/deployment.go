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
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ci4rail/kyt/kyt-alm-server/internal/controller"
	"github.com/ci4rail/kyt/kyt-alm-server/internal/deployment/manifest"
)

const (
	defaultPriority = 1
)

type DeploymentInterface interface {
	ApplyDeployment() (bool, error)
	ListDeployments(string) ([]string, error)
	GetLatestDeployment(string, string, string, string) (string, error)
	createManifestFile() (*os.File, error)
}

type Deployment struct {
	name             string
	connectionString string
	manifest         string
	hubName          string
	priority         int
	targetCondition  string
	now              string
}

// NewDeployment is used to define a new deployment
func NewDeployment(manifest string, name string, targetCondition string, now int64) (*Deployment, error) {
	c, err := controller.MapTenantToIOTHubSAS("")
	if err != nil {
		return nil, err
	}
	h, err := controller.IotHubNameFromConnecetionString(c)
	if err != nil {
		return nil, err
	}
	return &Deployment{
		name:             strings.ToLower(name),
		connectionString: c,
		manifest:         manifest,
		hubName:          h,
		priority:         defaultPriority,
		targetCondition:  targetCondition,
		now:              fmt.Sprintf("%d", now),
	}, nil
}

// ApplyDeployment applies the deployment to the backend service
func (d *Deployment) ApplyDeployment() (bool, error) {
	manifestFile, err := d.createManifestFile()
	if err != nil {
		return false, err
	}
	defer os.Remove(manifestFile.Name())

	azExecutable, err := exec.LookPath("az")
	if err != nil {
		return false, err
	}
	nameWithTimestamp := fmt.Sprintf("%s_%s", d.name, d.now)
	priority := strconv.Itoa(defaultPriority)
	cmdArgs := fmt.Sprintf("%s iot edge deployment create --hub-name %s --content %s --priority %s --layered --target-condition \"%s\" --deployment-id %s --login '%s'", azExecutable, d.hubName, manifestFile.Name(), priority, d.targetCondition, nameWithTimestamp, d.connectionString)
	fmt.Println("sh", "-c", cmdArgs)
	cmd := exec.Command("sh", "-c", cmdArgs)

	err = cmd.Start()
	if err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				err = fmt.Errorf("Exit Status: %d", status.ExitStatus())
				return false, err
			}
		} else {
			return false, err
		}
	}
	return true, nil
}

// DeleteDeployment deletes a specified deployment to the backend service
func (d *Deployment) DeleteDeployment() (bool, error) {
	azExecutable, err := exec.LookPath("az")
	if err != nil {
		return false, err
	}
	cmdArgs := fmt.Sprintf("%s iot edge deployment delete --hub-name %s --deployment-id %s --login '%s'", azExecutable, d.hubName, d.name, d.connectionString)
	fmt.Println("sh", "-c", cmdArgs)
	cmd := exec.Command("sh", "-c", cmdArgs)

	err = cmd.Start()
	if err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				err = fmt.Errorf("Exit Status: %d", status.ExitStatus())
				return false, err
			}
		} else {
			return false, err
		}
	}
	return true, nil
}

// ListDeployments gets all deployments from the backend service
func ListDeployments(connectionString string) ([]string, error) {
	c, err := controller.NewIOTHubServiceClient(connectionString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	deployments, err := c.ListDeployments()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []string{}
	for _, d := range deployments {
		result = append(result, d.ID)
	}
	return result, nil
}

func deploymentNameValid(name string) bool {
	re := regexp.MustCompile(`^[a-z0-9-]+_[a-z0-9-]+_[0-9]+$`)
	if ok := re.MatchString(name); ok {
		return true
	}
	return false
}

// GetLatestDeployment gets the last deployment from tenandId with application
func GetLatestDeployment(deployments []string, tenantId string, application string) (string, error) {
	appName := fmt.Sprintf("%s_%s", tenantId, application)
	tenantDeployments := []string{}
	for _, d := range deployments {
		if strings.Contains(d, appName) {
			if deploymentNameValid(d) {
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

func (d *Deployment) createManifestFile() (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "manifest")
	if err != nil {
		return nil, err
	}
	if _, err := tmpfile.Write([]byte(d.manifest)); err != nil {
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}
	return tmpfile, nil
}

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
	} else {
		return 0, fmt.Errorf("Error: no timestamp found")
	}
}

// CreateOrUpdateFromCustomerDeployment creates a new deployment and deletes the old one if it was
// already present.
func CreateOrUpdateFromCustomerDeployment(tenantId string, c *manifest.CustomerManifest) (bool, error) {
	cs, err := controller.ReadConnectionStringFromEnv()
	if err != nil {
		return false, err
	}
	// Get all deployments fron backend service
	deployments, err := ListDeployments(cs)
	if err != nil {
		return false, err
	}
	// Get latest deployment for specific tenantId
	latestToBeDeletedOnSuccess, err := GetLatestDeployment(deployments, tenantId, c.Application)
	if err != nil {
		fmt.Println(err)
	}
	// The new deployment needs to be created first to start the update process
	_, err = CreateDeploymentFromCustomerDeployment(tenantId, c)
	if err != nil {
		return false, err
	}
	// Delete old deployment with the same name to finish the update process
	if len(latestToBeDeletedOnSuccess) > 0 {
		// create new dummy deployment with specific name to be deleted
		deleteDeployment, err := NewDeployment("{}", latestToBeDeletedOnSuccess, "", 0)
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

func CreateDeploymentFromCustomerDeployment(tenantId string, c *manifest.CustomerManifest) (bool, error) {
	now := time.Now().Unix()
	layered, err := manifest.CreateLayeredManifest(c, tenantId)
	if err != nil {
		return false, err
	}
	deploymentName := fmt.Sprintf("%s_%s", tenantId, c.Application)
	// Currently the target condition is fixed to the tenant's ID
	targetCondition := fmt.Sprintf("tags.alm=true AND tags.tenantId='%s'", tenantId)
	d, err := NewDeployment(layered, deploymentName, targetCondition, now)
	if err != nil {
		return false, err
	}
	ok, err := d.ApplyDeployment()
	if err != nil {
		return false, err
	}
	return ok, nil
}
