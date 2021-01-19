package openapi

import (
	"fmt"
	"os"
)

// MapTenantToIOTHubSAS returns the SAS token of the IOT Hub for the specified tenant
// TODO: Either take the SAS from a DB or get it via "az iot hub connection-string show"
// TODO: tenant is currently ignored
func MapTenantToIOTHubSAS(tenant string) (string, error) {
	envName := fmt.Sprintf("IOTHUB_SERVICE_CONNECTION_STRING")
	val, ok := os.LookupEnv(envName)

	if !ok {
		return "", fmt.Errorf("IOTHUB_SERVICE_CONNECTION_STRING not set")
	}
	return val, nil
}
