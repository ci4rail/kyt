package moduleclient

import (
	"context"

	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
)

// DeviceInfo is the type for passing the device info tuples
type DeviceInfo map[string]interface{}

// NewModule creates a new Iotedge Module client from the environment
func NewModule() (*iotdevice.ModuleClient, error) {
	c, err := iotdevice.NewModuleFromEnvironment(
		// <transport>, <use iotedge gateway for connection>,
		iotmqtt.NewModuleTransport(), true)
	return c, err
}

// SetStaticDeviceInfo writes device info to Iothub device twin
func SetStaticDeviceInfo(c *iotdevice.ModuleClient, d DeviceInfo) error {
	// connect to the iothub
	if err := c.Connect(context.Background()); err != nil {
		return err
	}

	s := iotdevice.TwinState{
		"version": "abc",
	}
	if _, err := c.UpdateTwinState(context.Background(), s); err != nil {
		return err
	}
	return nil
}
