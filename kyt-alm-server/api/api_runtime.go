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

	iothubservice "github.com/ci4rail/kyt/kyt-alm-server/internal/iothubservice"
	"github.com/gin-gonic/gin"
)

// RuntimesRidGet - Get runtime by id rid
func RuntimesRidGet(c *gin.Context) {
	iotHubConnectionString, err := iothubservice.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothubservice.ControllerNewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	runtimeIDFilter := c.Param("rid")
	runtimeID, err := client.ListRuntimeByID(runtimeIDFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	connection, err := client.GetConnectionState(runtimeIDFilter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, &Runtime{
		ID:      *runtimeID,
		Network: connection,
	})
}

// RuntimesGet - List runtimes for a tenant
func RuntimesGet(c *gin.Context) {
	iotHubConnectionString, err := iothubservice.MapTenantToIOTHubSAS("")
	if err != nil {
		log.Fatal(err)
	}

	client, err := iothubservice.ControllerNewIOTHubServiceClient(iotHubConnectionString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	runtimeIDs, err := client.ListRuntimeIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var runtimeList []Runtime
	for _, runtimeID := range *runtimeIDs {
		connection, err := client.GetConnectionState(runtimeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		runtimeList = append(runtimeList, Runtime{
			ID:      runtimeID,
			Network: connection,
		})
	}

	c.JSON(http.StatusOK, runtimeList)
}
