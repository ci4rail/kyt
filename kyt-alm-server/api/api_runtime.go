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

	"github.com/ci4rail/kyt/kyt-serveres-common/iothub"
	"github.com/gin-gonic/gin"
	"github.com/golangci/golangci-lint/pkg/sliceutil"
)

// RuntimesRidGet -
func RuntimesRidGet(c *gin.Context) {
	token, err := ReadToken(c.Request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	tokenValid, claims, err := ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	if !sliceutil.Contains(claims, "RuntimesGet.read") {
		err = fmt.Errorf("Error: not allowed")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	var tenantID string
	if tenantID = tenantIDFromToken(token); tenantID == "" {
		err = fmt.Errorf("Error: reading user ID `oid` from token")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	// If token is not valid it means that either it has expired or the signature cannot be validated.
	// In both cases return `Unauthorized`.
	if !tokenValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	iotHubConnectionString, err := iothub.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothub.NewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	runtimeIDFilter := c.Param("rid")
	deviceID, err := client.ListRuntimeByID(tenantID, runtimeIDFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	connection, err := client.GetConnectionState(tenantID, runtimeIDFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, &Runtime{
		ID:      *deviceID,
		Network: connection,
	})
}

// RuntimesGet - List devices for a tenant
func RuntimesGet(c *gin.Context) {
	token, err := ReadToken(c.Request)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "expired") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	tokenValid, claims, err := ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	if !sliceutil.Contains(claims, "RuntimesGet.read") {
		err = fmt.Errorf("Error: not allowed")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	var tenantID string
	if tenantID = tenantIDFromToken(token); tenantID == "" {
		err = fmt.Errorf("Error: reading user ID `oid` from token")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	// If token is not valid it means that either it has expired or the signature cannot be validated.
	// In both cases return `Unauthorized`.
	if !tokenValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	iotHubConnectionString, err := iothub.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothub.NewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deviceIDs, err := client.ListRuntimeIDs(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var runtimeList []Runtime
	for _, deviceID := range *deviceIDs {
		connection, err := client.GetConnectionState(tenantID, deviceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		runtimeList = append(runtimeList, Runtime{
			ID:      deviceID,
			Network: connection,
		})
	}

	c.JSON(http.StatusOK, runtimeList)
}
