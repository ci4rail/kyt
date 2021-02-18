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

	"github.com/ci4rail/kyt/kyt-alm-server/internal/controller"
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
}

// NewDeployment is used to define a new deployment
func NewDeployment(manifest string, name string, targetCondition string) (*Deployment, error) {
	c, err := controller.MapTenantToIOTHubSAS("")
	if err != nil {
		return nil, err
	}
	h, err := controller.IotHubNameFromConnecetionString(c)
	if err != nil {
		return nil, err
	}
	return &Deployment{
		name:             name,
		connectionString: c,
		manifest:         manifest,
		hubName:          h,
		priority:         defaultPriority,
		targetCondition:  targetCondition,
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

	priority := strconv.Itoa(defaultPriority)
	cmdArgs := fmt.Sprintf("%s iot edge deployment create --hub-name %s --content %s --priority %s --layered --target-condition \"%s\" --deployment-id %s --login '%s'", azExecutable, d.hubName, manifestFile.Name(), priority, d.targetCondition, d.name, d.connectionString)
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

// GetLatestDeployment gets the last deployment from tenandId with application
func GetLatestDeployment(deployments []string, tenantId string, application string) (string, error) {
	appName := fmt.Sprintf("%s.%s", tenantId, application)
	tenantDeployments := []string{}
	for _, d := range deployments {
		if strings.Contains(d, appName) {
			tenantDeployments = append(tenantDeployments, d)
		}
	}
	sort.Sort(ByTimestamp(tenantDeployments))

	if len(tenantDeployments) > 0 {
		return tenantDeployments[0], nil
	}
	return "", fmt.Errorf("Error: no deployments found")
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
	// This regex pattern finds any numbers with leading `.` in a string.
	// e.g. tenant.application.321553 results in -> `.321553`
	re := regexp.MustCompile(`([.][0-9]+)$`)
	if ok := re.MatchString(deployment); ok {
		matches := re.FindAllString(deployment, -1)
		trimmed := strings.Replace(matches[0], ".", "", -1)
		timestamp, err := strconv.Atoi(trimmed)
		if err != nil {
			return 0, err
		}
		return timestamp, nil
	} else {
		return 0, fmt.Errorf("Error: no timestamp found")
	}
}
