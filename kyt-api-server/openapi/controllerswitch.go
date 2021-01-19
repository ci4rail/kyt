package openapi

import (
	"github.com/ci4rail/kyt-cli/kyt-api-server/controller"
	"github.com/ci4rail/kyt-cli/kyt-api-server/controllerif"
)

// ControllerNewIOTHubServiceClient points to the actual controller's NewIOTHubServiceClient function.
// Can be re-assigned to a stub for testing
var ControllerNewIOTHubServiceClient controllerif.NewIOTHubServiceClient = controller.NewIOTHubServiceClient