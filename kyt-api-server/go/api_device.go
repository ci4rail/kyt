/*
 * KYT API Server
 *
 * This is the KYT API Server
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// GetDevices - List devices for a tenant
func GetDevices(c *gin.Context) {
	client, err := ControllerNewIOTHubServiceClient(os.Getenv("IOTHUB_SERVICE_CONNECTION_STRING"))

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
