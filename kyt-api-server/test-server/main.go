package main

import (
	"fmt"
	"log"

	"github.com/ci4rail/kyt-cli/kyt-api-server/controllerif"
	sw "github.com/ci4rail/kyt-cli/kyt-api-server/go"
)

func main() {
	log.Printf("Test Server with dummy data started")

	// switch controller to stub functions
	sw.ControllerNewIOTHubServiceClient = newIOTHubStubClient

	router := sw.NewRouter()

	log.Fatal(router.Run(":9091"))
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
	devList := []string{
		"1234",
		"5678",
	}
	return &devList, nil
}
