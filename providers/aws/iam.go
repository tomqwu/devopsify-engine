package aws

import (
	"context"

	"github.com/deepnative/engine/pkg/resource"
)

// listIAMUsers returns IAM users as universal resources.
// This is a scaffold - actual AWS SDK calls will replace the placeholder.
func listIAMUsers(_ context.Context, _ *AuthConfig, region string) ([]resource.Resource, error) {
	_ = region
	return []resource.Resource{}, nil
}
