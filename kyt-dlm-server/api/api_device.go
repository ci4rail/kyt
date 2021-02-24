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
	t "github.com/ci4rail/kyt/kyt-server-common/token"
	"github.com/gin-gonic/gin"
	"github.com/golangci/golangci-lint/pkg/sliceutil"
)

// DevicesDidGet -
func DevicesDidGet(c *gin.Context) {
	token, err := t.ReadToken(c.Request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	tokenValid, claims, err := t.ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	if !sliceutil.Contains(claims, "DevicesDidGet.read") {
		err = fmt.Errorf("Error: not allowed")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	var tenantID string
	if tenantID = t.TenantIDFromToken(token); tenantID == "" {
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
	deviceIDFilter := c.Param("did")
	deviceID, err := client.ListDeviceByID(tenantID, deviceIDFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	connection, err := client.GetConnectionState(tenantID, deviceIDFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	versions, err := client.GetVersions(tenantID, deviceIDFilter)
	if err != nil {
		fmt.Printf("Info: device didn't repoart a version yet: %s\n", deviceIDFilter)
	}
	var firmwareVersion = ""
	f, ok := versions["firmwareVersion"]
	if ok {
		firmwareVersion = f
	}

	c.JSON(http.StatusOK, &Device{
		ID:              *deviceID,
		Network:         connection,
		FirmwareVersion: firmwareVersion,
	})
}

// DevicesGet - List devices for a tenant
func DevicesGet(c *gin.Context) {
	token, err := t.ReadToken(c.Request)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "expired") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	tokenValid, claims, err := t.ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	if !sliceutil.Contains(claims, "DevicesGet.read") {
		err = fmt.Errorf("Error: not allowed")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	var tenantID string
	if tenantID = t.TenantIDFromToken(token); tenantID == "" {
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

	deviceIDs, err := client.ListDeviceIDs(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var deviceList []Device
	for _, deviceID := range *deviceIDs {
		connection, err := client.GetConnectionState(tenantID, deviceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		versions, err := client.GetVersions(tenantID, deviceID)
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

	c.JSON(http.StatusOK, deviceList)
}
