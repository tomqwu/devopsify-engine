package aws

import (
	"context"
	"time"

	"github.com/deepnative/engine/pkg/resource"
)

// listIAMUsers returns IAM users as universal resources.
// This is a scaffold - actual AWS SDK calls will replace the placeholder.
func listIAMUsers(_ context.Context, _ *AuthConfig, region string) ([]resource.Resource, error) {
	_ = region
	return []resource.Resource{}, nil
}

// newIAMResource creates a universal resource from IAM user data.
func newIAMResource(id, name string, tags map[string]string) resource.Resource {
	return resource.Resource{
		ID:        id,
		Name:      name,
		Type:      "identity:user",
		Provider:  "aws",
		Region:    "global",
		State:     resource.StateHealthy,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
