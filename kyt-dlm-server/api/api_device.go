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

// DevicesDidGet gets a specific device for given tenants
func DevicesDidGet(w http.ResponseWriter, r *http.Request) {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenRead := authHeaderParts[1]

	hasScope := token.CheckScope("dlm/DevicesDidGet.read", tokenRead)

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
	deviceIDFilter := params["did"]

	var deviceID string = ""
	deviceFound := false
	// check if the device is owed by any of the given tenants
	for _, t := range tenants {
		d, err := client.ListDeviceByID(t, deviceIDFilter)
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
	connection, err := client.GetConnectionState(deviceIDFilter)
	if err != nil {
		responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
		return
	}
	versions, err := client.GetVersions(deviceIDFilter)
	if err != nil {
		fmt.Printf("Info: device didn't repoart a version yet: %s\n", deviceIDFilter)
	}
	var firmwareVersion = ""
	f, ok := versions["firmwareVersion"]
	if ok {
		firmwareVersion = f
	}

	responseJSON(&Device{
		ID:              deviceID,
		Network:         connection,
		FirmwareVersion: firmwareVersion,
	}, w, http.StatusOK)
}

// DevicesGet list devices for given tenants
func DevicesGet(w http.ResponseWriter, r *http.Request) {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenRead := authHeaderParts[1]

	hasScope := token.CheckScope("dlm/DevicesGet.read", tokenRead)

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
		d, err := client.ListDeviceIDs(t)
		if err != nil {
			responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
			return
		}
		deviceIDs = append(deviceIDs, *d...)
	}
	var deviceList []Device
	for _, deviceID := range deviceIDs {
		connection, err := client.GetConnectionState(deviceID)
		if err != nil {
			responseJSON(fmt.Sprintf("error: %s", err), w, http.StatusInternalServerError)
			return
		}
		versions, err := client.GetVersions(deviceID)
		if err != nil {
			fmt.Printf("Info: device didn't repoart a version yet: %s\n", deviceID)
		}
		var firmwareVersion = ""
		f, ok := versions["firmwareVersion"]
		if ok {
			firmwareVersion = f
		}
		deviceList = append(deviceList, Device{
			ID:              deviceID,
			Network:         connection,
			FirmwareVersion: firmwareVersion,
		})
	}
	responseJSON(deviceList, w, http.StatusOK)
}
