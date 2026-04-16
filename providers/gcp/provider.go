package gcp

import (
	"context"
	"fmt"

	"github.com/deepnative/engine/pkg/provider"
	"github.com/deepnative/engine/pkg/resource"
)

var _ provider.CloudProvider = (*Provider)(nil)

// Provider implements the CloudProvider interface for GCP.
type Provider struct {
	auth    *AuthConfig
	project string
}

func (p *Provider) Metadata() provider.Metadata {
	return provider.Metadata{
		Name:    "gcp",
		Kind:    "cloud",
		Version: "1.0.0",
	}
}

func (p *Provider) Init(_ context.Context, cfg map[string]any) error {
	auth := resolveAuth(cfg)
	if auth.Project == "" {
		return fmt.Errorf("%w: project is required", provider.ErrProviderInit)
	}
	p.auth = auth
	p.project = auth.Project
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

func (p *Provider) ListResources(_ context.Context, resourceType string, _ provider.ListOptions) ([]resource.Resource, error) {
	switch resourceType {
	case "compute:instance", "":
		return []resource.Resource{}, nil
	case "storage:bucket":
		return []resource.Resource{}, nil
	case "iam:serviceaccount":
		return []resource.Resource{}, nil
	default:
		return nil, fmt.Errorf("%w: resource type %s", provider.ErrUnsupportedOperation, resourceType)
	}
}

func (p *Provider) GetResource(_ context.Context, resourceType, id string) (*resource.Resource, error) {
	return nil, fmt.Errorf("%w: GetResource for %s/%s", provider.ErrUnsupportedOperation, resourceType, id)
}

func (p *Provider) GetCostData(_ context.Context, opts provider.CostQueryOptions) (*provider.CostReport, error) {
	return &provider.CostReport{
		Provider:    "gcp",
		StartDate:   opts.StartDate,
		EndDate:     opts.EndDate,
		TotalCost:   0,
		Currency:    "USD",
		CostsByItem: make(map[string]float64),
	}, nil
}

func (p *Provider) DetectDrift(ctx context.Context, desired []resource.Resource) ([]resource.DriftResult, error) {
	actual, err := p.ListResources(ctx, "", provider.ListOptions{})
	if err != nil {
		return nil, err
	}
	return resource.Diff(desired, actual), nil
}
