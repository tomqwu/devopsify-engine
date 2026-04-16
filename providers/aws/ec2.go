package aws

import (
	"context"
	"time"

	"github.com/deepnative/engine/pkg/resource"
)

// listEC2Instances returns EC2 instances as universal resources.
// This is a scaffold - actual AWS SDK calls will replace the placeholder.
func listEC2Instances(_ context.Context, _ *AuthConfig, region string) ([]resource.Resource, error) {
	// Placeholder: In production, this calls AWS EC2 DescribeInstances
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

// newEC2Resource creates a universal resource from EC2 instance data.
func newEC2Resource(id, name, state, instanceType, region string, tags map[string]string) resource.Resource {
	return resource.Resource{
		ID:       id,
		Name:     name,
		Type:     "compute:instance",
		Provider: "aws",
		Region:   region,
		State:    mapEC2State(state),
		Tags:     tags,
		Properties: map[string]any{
			"instance_type": instanceType,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
