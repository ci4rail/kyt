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

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	d "github.com/ci4rail/kyt/kyt-alm-server/internal/deployment"
	m "github.com/ci4rail/kyt/kyt-alm-server/internal/deployment/manifest"
)

// ApplyPut -
func ApplyPut(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("Error: cannot read payload")
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}
	manifest := m.CustomerManifest{}
	err = json.Unmarshal(jsonData, &manifest)
	if err != nil {
		err = fmt.Errorf("Error: cannot read customer manifest")
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}
	_, err = d.CreateOrUpdateFromCustomerDeployment("TODO_TenantID", &manifest)
	if err != nil {
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}
	responseJSON("OK", w, http.StatusOK)
}
