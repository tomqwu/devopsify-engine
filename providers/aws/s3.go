package aws

import (
	"context"

	"github.com/deepnative/engine/pkg/resource"
)

// listS3Buckets returns S3 buckets as universal resources.
// This is a scaffold - actual AWS SDK calls will replace the placeholder.
func listS3Buckets(_ context.Context, _ *AuthConfig, region string) ([]resource.Resource, error) {
	_ = region
	return []resource.Resource{}, nil
}
