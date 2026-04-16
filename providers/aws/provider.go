package aws

import (
	"context"
	"fmt"

	"github.com/deepnative/engine/pkg/provider"
	"github.com/deepnative/engine/pkg/resource"
)

// Compile-time interface compliance check.
var _ provider.CloudProvider = (*Provider)(nil)

// Provider implements the CloudProvider interface for AWS.
type Provider struct {
	auth   *AuthConfig
	region string
}

func (p *Provider) Metadata() provider.Metadata {
	return provider.Metadata{
		Name:    "aws",
		Kind:    "cloud",
		Version: "1.0.0",
	}
}

func (p *Provider) Init(_ context.Context, cfg map[string]any) error {
	auth, err := resolveAuth(cfg)
	if err != nil {
		return fmt.Errorf("resolving auth: %w", err)
	}
	if err := validateAuth(auth); err != nil {
		return fmt.Errorf("validating auth: %w", err)
	}
	p.auth = auth
	p.region = auth.Region
	return nil
}

func (p *Provider) Healthy(_ context.Context) error {
	if p.auth == nil {
		return provider.ErrProviderUnhealthy
	}
	return nil
}

func (p *Provider) Shutdown(_ context.Context) error {
	return nil
}

func (p *Provider) ListResources(ctx context.Context, resourceType string, opts provider.ListOptions) ([]resource.Resource, error) {
	region := opts.Region
	if region == "" {
		region = p.region
	}

	switch resourceType {
	case "compute:instance", "ec2":
		return listEC2Instances(ctx, p.auth, region)
	case "storage:bucket", "s3":
		return listS3Buckets(ctx, p.auth, region)
	case "identity:user", "iam":
		return listIAMUsers(ctx, p.auth, region)
	case "":
		var all []resource.Resource
		ec2, _ := listEC2Instances(ctx, p.auth, region)
		all = append(all, ec2...)
		s3, _ := listS3Buckets(ctx, p.auth, region)
		all = append(all, s3...)
		iam, _ := listIAMUsers(ctx, p.auth, region)
		all = append(all, iam...)
		return all, nil
	default:
		return nil, fmt.Errorf("%w: resource type %s", provider.ErrUnsupportedOperation, resourceType)
	}
}

func (p *Provider) GetResource(_ context.Context, resourceType, id string) (*resource.Resource, error) {
	return nil, fmt.Errorf("%w: GetResource for %s/%s", provider.ErrUnsupportedOperation, resourceType, id)
}

func (p *Provider) GetCostData(ctx context.Context, opts provider.CostQueryOptions) (*provider.CostReport, error) {
	return getCostExplorerData(ctx, p.auth, opts)
}

func (p *Provider) DetectDrift(ctx context.Context, desired []resource.Resource) ([]resource.DriftResult, error) {
	actual, err := p.ListResources(ctx, "", provider.ListOptions{})
	if err != nil {
		return nil, err
	}
	return resource.Diff(desired, actual), nil
}
