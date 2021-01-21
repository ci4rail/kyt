package main

import (
	"fmt"
	"log"

	"github.com/ci4rail/kyt-cli/kyt-api-server/cmd"
	"github.com/ci4rail/kyt-cli/kyt-api-server/controllerif"
	"github.com/ci4rail/kyt-cli/kyt-api-server/iothubservice"
)

func main() {
	log.Printf("Test Server with dummy data started")

	// switch controller to stub functions
	iothubservice.ControllerNewIOTHubServiceClient = newIOTHubStubClient
	cmd.Execute()

}

type stubbedIOTHubServiceClient struct {
	controllerif.IOTHubServices
}

func newIOTHubStubClient(connectionString string) (controllerif.IOTHubServices, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("Missing IoTHub connection string")
	}

	c := &stubbedIOTHubServiceClient{}
	return c, nil
}

// ListDeviceIDs returns a list with the device IDs of all devices of that IoT Hub
func (c *stubbedIOTHubServiceClient) ListDeviceIDs() (*[]string, error) {
	const numDevs = 1000
	var devList [numDevs]string

	for i := 0; i < numDevs; i++ {
		devList[i] = fmt.Sprintf("Device %5d", i)
	}
	slice := devList[:]
	return &slice, nil
}
