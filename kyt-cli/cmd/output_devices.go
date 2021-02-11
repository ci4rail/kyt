/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

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

package cmd

import (
	"encoding/json"
	"fmt"

	openapi "github.com/ci4rail/kyt/kyt-cli/openapidlm"
	"github.com/ghodss/yaml"
)

// ConvertToJSON Convert devices infromation to JSON
func ConvertToJSON(devices *[]openapi.Device) (string, error) {
	c, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Cannot output as format json")
	}
	return string(c), nil
}

// ConvertToYaml Convert devices infromation to yaml
func ConvertToYaml(devices *[]openapi.Device) (string, error) {
	j, err := ConvertToJSON(devices)
	if err != nil {
		return "", err
	}
	c, err := yaml.JSONToYAML([]byte(j))
	if err != nil {
		return "", fmt.Errorf("Cannot output as format yaml")
	}
	return string(c), nil
}
