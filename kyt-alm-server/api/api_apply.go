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
	"strings"

	d "github.com/ci4rail/kyt/kyt-alm-server/internal/deployment"
	m "github.com/ci4rail/kyt/kyt-alm-server/internal/deployment/manifest"
	token "github.com/ci4rail/kyt/kyt-server-common/token"
)

// ApplyPut applies a customer manifest to the tenants devices
func ApplyPut(w http.ResponseWriter, r *http.Request) {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenRead := authHeaderParts[1]

	hasScope := token.CheckScope("alm/Apply.write", tokenRead)

	tenants := token.GetTenants(tokenRead)

	if !hasScope {
		responseJSON("Error: insufficient scope.", w, http.StatusForbidden)
		return
	}

	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("Error: cannot read payload")
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}

	errStr := ""
	for _, tenant := range tenants {
		// Unmarshalling for every tenant currently is a workaround for a strange bug
		// that has not been identified yet. The bug modifies `manifest` somewhere within
		// function `CreateOrUpdateFromCustomerDeployment` even `manifest` gets passed by value.
		manifest, err := unmarshalCustomerManifest(jsonData)
		if err != nil {
			err = fmt.Errorf("Error: cannot read customer manifest")
			responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
			return
		}
		fmt.Printf("Writing deployment for tenant: %s\n", tenant)
		err = d.CreateOrUpdateFromCustomerDeployment(tenant, manifest)
		if err != nil {
			errStr += fmt.Sprintf("%s: %s\n", tenant, err)
			continue
		}
	}
	if len(errStr) > 0 {
		responseJSON(fmt.Sprintf("error: failed applying deployment for:\n%s", errStr), w, http.StatusInternalServerError)
		return
	}
	responseJSON("OK", w, http.StatusOK)
}

func unmarshalCustomerManifest(jsonData []byte) (m.CustomerManifest, error) {
	manifest := m.CustomerManifest{}
	err := json.Unmarshal(jsonData, &manifest)
	if err != nil {
		return manifest, err
	}
	return manifest, nil
}
