package controllerif

// NewIOTHubServiceClient points to the actual creator function
type NewIOTHubServiceClient func(connectionString string) (IOTHubServices, error)

// IOTHubServices has all the services the IOTHub Controller offers
type IOTHubServices interface {
	ListDeviceIDs() (*[]string, error)
}
