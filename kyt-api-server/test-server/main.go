package main

import (
	"log"

	"github.com/ci4rail/kyt-cli/kyt-api-server/controllerif"
	sw "github.com/ci4rail/kyt-cli/kyt-api-server/go"
)

func main() {
	log.Printf("Test Server started")

	sw.ControllerNewIOTHubServiceClient = NewIOTHubServiceClient

	router := sw.NewRouter()

	log.Fatal(router.Run(":9091"))
}

// IOTHubServiceClient is an Azure IoT Hub service client.
type IOTHubServiceClient struct {
	controllerif.IOTHubServices
}

// NewIOTHubServiceClient creates a new IOTHubServiceClient based on the connection string
// connection string can be determined with "az iot hub connection-string show"
func NewIOTHubServiceClient(connectionString string) (controllerif.IOTHubServices, error) {
	c := &IOTHubServiceClient{}
	return c, nil
}

// ListDeviceIDs returns a list with the device IDs of all devices of that IoT Hub
func (c *IOTHubServiceClient) ListDeviceIDs() (*[]string, error) {
	devList := []string{
		"1234",
		"5678",
	}
	return &devList, nil
}
