package aws

import (
	"context"
	"time"

	"github.com/deepnative/engine/pkg/resource"
)

// listS3Buckets returns S3 buckets as universal resources.
// This is a scaffold - actual AWS SDK calls will replace the placeholder.
func listS3Buckets(_ context.Context, _ *AuthConfig, region string) ([]resource.Resource, error) {
	_ = region
	return []resource.Resource{}, nil
}

// newS3Resource creates a universal resource from S3 bucket data.
func newS3Resource(name, region string, tags map[string]string) resource.Resource {
	return resource.Resource{
		ID:        name,
		Name:      name,
		Type:      "storage:bucket",
		Provider:  "aws",
		Region:    region,
		State:     resource.StateHealthy,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
