package opsgenie

import (
	"context"
	"fmt"

	"github.com/deepnative/engine/internal/config"
	"github.com/deepnative/engine/pkg/provider"
)

var _ provider.SREProvider = (*Provider)(nil)

// Provider implements the SREProvider interface for OpsGenie.
type Provider struct {
	apiKey string
}

func (p *Provider) Metadata() provider.Metadata {
	return provider.Metadata{
		Name:    "opsgenie",
		Kind:    "sre",
		Version: "1.0.0",
	}
}

func (p *Provider) Init(_ context.Context, cfg map[string]any) error {
	p.apiKey = config.ResolveEnvValue(cfg, "api_key")
	if p.apiKey == "" {
		return fmt.Errorf("%w: api_key is required", provider.ErrProviderInit)
	}
	return nil
}

func (p *Provider) Healthy(_ context.Context) error {
	if p.apiKey == "" {
		return provider.ErrProviderUnhealthy
	}
	return nil
}

func (p *Provider) Shutdown(_ context.Context) error {
	return nil
}

func (p *Provider) ListIncidents(_ context.Context, _ provider.IncidentListOptions) ([]provider.Incident, error) {
	return []provider.Incident{}, nil
}

func (p *Provider) GetIncident(_ context.Context, id string) (*provider.Incident, error) {
	return nil, fmt.Errorf("%w: GetIncident for %s", provider.ErrUnsupportedOperation, id)
}

func (p *Provider) AcknowledgeIncident(_ context.Context, id string) error {
	return fmt.Errorf("%w: AcknowledgeIncident for %s", provider.ErrUnsupportedOperation, id)
}

func (p *Provider) ResolveIncident(_ context.Context, id string) error {
	return fmt.Errorf("%w: ResolveIncident for %s", provider.ErrUnsupportedOperation, id)
}

func (p *Provider) ListOnCall(_ context.Context, _ string) ([]provider.OnCallEntry, error) {
	return []provider.OnCallEntry{}, nil
}
