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

package apply

import (
	"fmt"
	"io/ioutil"

	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"
	openapi "github.com/ci4rail/kyt/kyt-cli/openapialm"
	"github.com/ghodss/yaml"
)

// ReadCustomerManifest -
func ReadCustomerManifest(filename string) *openapi.CustomerManifest {

	c := openapi.NewCustomerManifestWithDefaults()
	c = readYaml(filename)

	return c
}

func readYaml(filename string) *openapi.CustomerManifest {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		e.Er(fmt.Sprintf("read YAML file error: %v\n", err))
	}

	c := openapi.NewCustomerManifestWithDefaults()

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		e.Er(fmt.Sprintf("parsing YAML failed: %v\n", err))
	}

	return c
}
