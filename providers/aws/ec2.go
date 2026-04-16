package aws

import (
	"context"

	"github.com/deepnative/engine/pkg/resource"
)

// listEC2Instances returns EC2 instances as universal resources.
// This is a scaffold - actual AWS SDK calls will replace the placeholder.
func listEC2Instances(_ context.Context, _ *AuthConfig, region string) ([]resource.Resource, error) {
	_ = region
	return []resource.Resource{}, nil
}

// mapEC2State maps AWS EC2 instance states to universal resource states.
func mapEC2State(state string) resource.State {
	switch state {
	case "running":
		return resource.StateRunning
	case "stopped":
		return resource.StateStopped
	case "terminated":
		return resource.StateTerminated
	case "pending":
		return resource.StatePending
	case "shutting-down":
		return resource.StateStopped
	default:
		return resource.StateUnknown
	}
}
