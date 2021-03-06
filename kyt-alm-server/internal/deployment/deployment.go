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
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ci4rail/kyt/kyt-alm-server/internal/common"
	iothub "github.com/ci4rail/kyt/kyt-server-common/iothubwrapper"
)

const (
	customerDeploymentPriorityOffset = 100
	baseDeploymentName               = "base_deployment"
)

// Interface used for unit testing
type Interface interface {
	applyDeployment() (bool, error)
	ListDeployments(string) ([]string, error)
	getLatestCustomerDeployment(string, string, string, string) (string, error)
	createManifestFile() (*os.File, error)
}

// Deployment is the struct that contains all relevant data for a deployment
type Deployment struct {
	name             string
	connectionString string
	manifest         string
	hubName          string
	priority         int
	targetCondition  string
	layered          bool
	version          string
}

// NewDeployment is used to define a new deployment
func NewDeployment(manifest string, name string, targetCondition string, layered bool, version string, priority int) (*Deployment, error) {
	c, err := iothub.MapTenantToIOTHubSAS("")
	if err != nil {
		return nil, err
	}
	h, err := iothub.IotHubNameFromConnecetionString(c)
	if err != nil {
		return nil, err
	}
	d := &Deployment{
		name:             strings.ToLower(name),
		connectionString: c,
		manifest:         manifest,
		hubName:          h,
		priority:         priority,
		targetCondition:  targetCondition,
		layered:          layered,
		version:          version,
	}
	// Layered deployments are customer deployments. Customer deployments have priority offset.
	if d.layered {
		d.priority += customerDeploymentPriorityOffset
	}
	return d, nil
}

// DeleteDeployment deletes a specified deployment to the backend service
func (d *Deployment) DeleteDeployment() (bool, error) {
	azExecutable, err := exec.LookPath("az")
	if err != nil {
		return false, err
	}
	cmdArgs := fmt.Sprintf("%s iot edge deployment delete --hub-name %s --deployment-id %s --login '%s'", azExecutable, d.hubName, d.name, d.connectionString)
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
	c, err := iothub.NewIOTHubServiceClient(connectionString)
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

// applyDeployment writes the deployment to the backend service
func (d *Deployment) applyDeployment() (bool, error) {
	if !common.DryRun {
		manifestFile, err := d.createManifestFile()
		if err != nil {
			return false, err
		}
		defer os.Remove(manifestFile.Name())

		azExecutable, err := exec.LookPath("az")
		if err != nil {
			return false, err
		}
		nameWithVersion := fmt.Sprintf("%s_%s", d.name, d.version)
		cmdArgs := fmt.Sprintf("%s iot edge deployment create --hub-name %s --content %s --priority %d --target-condition \"%s\" --deployment-id %s --login '%s'", azExecutable, d.hubName, manifestFile.Name(), d.priority, d.targetCondition, nameWithVersion, d.connectionString)
		if d.layered {
			cmdArgs += " --layered"
		}
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
	} else {
		log.Printf("Dry run: not applying deployment: %s", d.name)
	}
	return true, nil
}

//createManifestFile writes the manifest string to a temp file
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
