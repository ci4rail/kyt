package controller

import (
	"context"
	"fmt"

	"github.com/amenzhinsky/iothub/iotservice"
)

// IOTHubServiceClient is an Azure IoT Hub service client.
type IOTHubServiceClient struct {
	iotClient *iotservice.Client
}

// NewIOTHubServiceClient creates a new IOTHubServiceClient based on the connection string
// connection string can be determined with "az iot hub connection-string show"
func NewIOTHubServiceClient(connectionString string) (*IOTHubServiceClient, error) {
	c := &IOTHubServiceClient{}

	iotClient, err := iotservice.NewFromConnectionString(connectionString)

	if err != nil {
		return nil, fmt.Errorf("Can't create IoT Hub Client %s", err)
	}

	c.iotClient = iotClient
	return c, nil
}

// ListDeviceIDs returns a list with the device IDs of all devices of that IoT Hub
func (c *IOTHubServiceClient) ListDeviceIDs() (*[]string, error) {
	var deviceIDArr []string

	ctx := context.Background()

	devices, err := c.iotClient.ListDevices(ctx)

	if err != nil {
		return nil, fmt.Errorf("Error IoT Hub ListDevices %s", err)
	}

	for _, device := range devices {
		deviceIDArr = append(deviceIDArr, device.DeviceID)
	}
	return &deviceIDArr, nil
}
