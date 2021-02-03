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
	"log"
	"net/http"

	iothubservice "github.com/ci4rail/kyt/kyt-api-server/internal/iothubservice"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, &Device{Id: *deviceID})
}

// DevicesGet - List devices for a tenant
func DevicesGet(c *gin.Context) {
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
		deviceList = append(deviceList, Device{Id: deviceID})
	}

	c.JSON(http.StatusOK, deviceList)
}
