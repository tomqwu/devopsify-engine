package githubactions

import (
	"context"
	"fmt"
	"time"

	"github.com/deepnative/engine/internal/config"
	"github.com/deepnative/engine/pkg/provider"
)

var _ provider.PipelineProvider = (*Provider)(nil)

// Provider implements the PipelineProvider interface for GitHub Actions.
type Provider struct {
	owner string
	repo  string
	token string
}

func (p *Provider) Metadata() provider.Metadata {
	return provider.Metadata{
		Name:    "github-actions",
		Kind:    "pipeline",
		Version: "1.0.0",
	}
}

func (p *Provider) Init(_ context.Context, cfg map[string]any) error {
	p.owner = config.ResolveEnvValue(cfg, "owner")
	if p.owner == "" {
		return fmt.Errorf("%w: owner is required", provider.ErrProviderInit)
	}
	p.repo = config.ResolveEnvValue(cfg, "repo")
	if p.repo == "" {
		return fmt.Errorf("%w: repo is required", provider.ErrProviderInit)
	}
	p.token = config.ResolveEnvValue(cfg, "token")
	if p.token == "" {
		return fmt.Errorf("%w: token is required", provider.ErrProviderInit)
	}
	return nil
}

func (p *Provider) Healthy(_ context.Context) error {
	if p.token == "" {
		return provider.ErrProviderUnhealthy
	}
	return nil
}

func (p *Provider) Shutdown(_ context.Context) error {
	return nil
}

func (p *Provider) ListPipelines(_ context.Context, _ provider.PipelineListOptions) ([]provider.Pipeline, error) {
	return []provider.Pipeline{}, nil
}

func (p *Provider) GetPipelineStatus(_ context.Context, id string) (*provider.PipelineStatus, error) {
	return nil, fmt.Errorf("%w: GetPipelineStatus for %s", provider.ErrUnsupportedOperation, id)
}

func (p *Provider) TriggerSync(_ context.Context, id string, _ provider.SyncOptions) (*provider.SyncResult, error) {
	return &provider.SyncResult{
		ID:        id,
		Status:    "queued",
		Message:   "workflow dispatch triggered",
		StartedAt: time.Now(),
	}, nil
}

func (p *Provider) GetHistory(_ context.Context, _ string, _ int) ([]provider.PipelineRun, error) {
	return []provider.PipelineRun{}, nil
}
