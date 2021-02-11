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

package alm

import (
	"encoding/json"
	"fmt"

	openapi "github.com/ci4rail/kyt/kyt-cli/openapialm"
	"github.com/ghodss/yaml"
)

// ConvertToJSON Convert runtimes infromation to JSON
func ConvertToJSON(runtimes *[]openapi.Runtime) (string, error) {
	c, err := json.MarshalIndent(runtimes, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Cannot output as format json")
	}
	return string(c), nil
}

// ConvertToYaml Convert runtimes infromation to yaml
func ConvertToYaml(runtimes *[]openapi.Runtime) (string, error) {
	j, err := ConvertToJSON(runtimes)
	if err != nil {
		return "", err
	}
	c, err := yaml.JSONToYAML([]byte(j))
	if err != nil {
		return "", fmt.Errorf("Cannot output as format yaml")
	}
	return string(c), nil
}
