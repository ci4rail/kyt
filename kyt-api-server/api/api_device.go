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

	iothubservice "github.com/ci4rail/kyt/kyt-api-server/internal/iothubservice"
	"github.com/gin-gonic/gin"
	"github.com/golangci/golangci-lint/pkg/sliceutil"
)

// DevicesDidGet -
func DevicesDidGet(c *gin.Context) {
	iotHubConnectionString, err := iothubservice.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothubservice.ControllerNewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	deviceIdFilter := c.Param("did")
	deviceID, err := client.ListDeviceById(deviceIdFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	connection, err := client.GetConnectionState(deviceIdFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	versions, err := client.GetVersions(deviceIdFilter)
	if err != nil {
		fmt.Printf("Info: device didn't repoart a version yet: %s\n", deviceIdFilter)
	}
	var firmwareVersion = ""
	f, ok := versions["firmwareVersion"]
	if ok {
		firmwareVersion = f
	}

	c.JSON(http.StatusOK, &Device{
		Id:              *deviceID,
		Network:         connection,
		FirmwareVersion: firmwareVersion,
	})
}

// DevicesGet - List devices for a tenant
func DevicesGet(c *gin.Context) {
	tokenValid, claims, err := validateToken(c.Request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	if !sliceutil.Contains(claims, "DevicesGet.read") {
		err = fmt.Errorf("Error: not allowed.")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	// If token is not valid it means that either it has expired or the signature cannot be validated.
	// In both cases return `Unauthorized`.
	if !tokenValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	iotHubConnectionString, err := iothubservice.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothubservice.ControllerNewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deviceIDs, err := client.ListDeviceIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var deviceList []Device
	for _, deviceID := range *deviceIDs {
		connection, err := client.GetConnectionState(deviceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
			Id:              deviceID,
			Network:         connection,
			FirmwareVersion: firmwareVersion,
		})
	}

	c.JSON(http.StatusOK, deviceList)
}
