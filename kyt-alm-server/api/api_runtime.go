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
)

// // RuntimesRidGet -
// func RuntimesRidGet(c *gin.Context) {
// 	token, err := t.ReadToken(c.Request)
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
// 		return
// 	}
// 	tokenValid, claims, err := t.ValidateToken(token)
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
// 		return
// 	}
// 	if !sliceutil.Contains(claims, "RuntimesGet.read") {
// 		err = fmt.Errorf("Error: not allowed")
// 		c.JSON(http.StatusForbidden, gin.H{"error": err})
// 		return
// 	}
// 	var tenantID string
// 	if tenantID = t.TenantIDFromToken(token); tenantID == "" {
// 		err = fmt.Errorf("Error: reading user ID `oid` from token")
// 		c.JSON(http.StatusForbidden, gin.H{"error": err})
// 		return
// 	}
// 	// If token is not valid it means that either it has expired or the signature cannot be validated.
// 	// In both cases return `Unauthorized`.
// 	if !tokenValid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
// 		return
// 	}

// 	iotHubConnectionString, err := iothub.MapTenantToIOTHubSAS("")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	client, err := iothub.NewIOTHubServiceClient(iotHubConnectionString)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	runtimeIDFilter := c.Param("rid")
// 	deviceID, err := client.ListRuntimeByID(tenantID, runtimeIDFilter)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}
// 	connection, err := client.GetConnectionState(tenantID, runtimeIDFilter)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err})
// 		return
// 	}

// 	c.JSON(http.StatusOK, &Runtime{
// 		ID:      *deviceID,
// 		Network: connection,
// 	})
// }

// RuntimesGet - List devices for a tenant
func RuntimesGet(w http.ResponseWriter, r *http.Request) {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenRead := authHeaderParts[1]

	hasScope := token.CheckScope("alm/RuntimesGet.read", tokenRead)

	tenants := token.GetTenants(tokenRead)

	if !hasScope {
		message := "Insufficient scope."
		responseJSON(message, w, http.StatusForbidden)
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
