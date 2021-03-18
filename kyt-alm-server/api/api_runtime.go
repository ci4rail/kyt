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
	"fmt"
	"log"
	"net/http"
	"strings"

	iothub "github.com/ci4rail/kyt/kyt-server-common/iothubwrapper"
	token "github.com/ci4rail/kyt/kyt-server-common/token"
	"github.com/gorilla/mux"
)

// RuntimesRidGet lists specific filtered devices for a list of tenants
func RuntimesRidGet(w http.ResponseWriter, r *http.Request) {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenRead := authHeaderParts[1]

	hasScope := token.CheckScope("alm/RuntimesGet.read", tokenRead)

	tenants := token.GetTenants(tokenRead)

	if !hasScope {
		responseJSON("Error: insufficient scope.", w, http.StatusForbidden)
		return
	}

	iotHubConnectionString, err := iothub.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothub.NewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	runtimeIDFilter := params["rid"]

	var deviceID string = ""
	deviceFound := false
	// check if the device is owed by any of the given tenants
	for _, t := range tenants {
		d, err := client.ListRuntimeByID(t, runtimeIDFilter)
		// move to next tenant in the list and try this tenant
		if err != nil {
			continue
		}
		// device found
		if d != nil {
			deviceID = *d
			deviceFound = true
			break
		}
	}
	if !deviceFound {
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusNotFound)
		return
	}

	connection, err := client.GetConnectionState(deviceID)
	if err != nil {
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}

	runtime := &Runtime{
		ID:      deviceID,
		Network: connection,
	}
	responseJSON(runtime, w, http.StatusOK)
}

// RuntimesGet lists all devices for a list of tenants
func RuntimesGet(w http.ResponseWriter, r *http.Request) {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenRead := authHeaderParts[1]

	hasScope := token.CheckScope("alm/RuntimesGet.read", tokenRead)

	tenants := token.GetTenants(tokenRead)

	if !hasScope {
		responseJSON("Error: insufficient scope.", w, http.StatusForbidden)
		return
	}

	iotHubConnectionString, err := iothub.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothub.NewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}

	var deviceIDs []string
	// get devices for owned by all given tenants
	for _, t := range tenants {
		d, err := client.ListRuntimeIDs(t)
		if err != nil {
			responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
			return
		}
		deviceIDs = append(deviceIDs, *d...)
	}
	var runtimeList []Runtime
	for _, deviceID := range deviceIDs {
		connection, err := client.GetConnectionState(deviceID)
		if err != nil {
			responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
			return
		}

		runtimeList = append(runtimeList, Runtime{
			ID:      deviceID,
			Network: connection,
		})
	}
	responseJSON(runtimeList, w, http.StatusOK)
}
