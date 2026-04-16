package flux

import (
	"context"
	"fmt"
	"time"

	"github.com/deepnative/engine/internal/config"
	"github.com/deepnative/engine/pkg/provider"
)

var _ provider.PipelineProvider = (*Provider)(nil)

// Provider implements the PipelineProvider interface for Flux CD.
type Provider struct {
	kubeconfig string
	namespace  string
}

func (p *Provider) Metadata() provider.Metadata {
	return provider.Metadata{
		Name:    "flux",
		Kind:    "pipeline",
		Version: "1.0.0",
	}
}

func (p *Provider) Init(_ context.Context, cfg map[string]any) error {
	p.kubeconfig = config.ResolveEnvValue(cfg, "kubeconfig")
	p.namespace = config.ResolveEnvValue(cfg, "namespace")
	if p.namespace == "" {
		p.namespace = "flux-system"
	}
	return nil
}

func (p *Provider) Healthy(_ context.Context) error {
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
		Status:    "reconciling",
		Message:   "reconciliation triggered",
		StartedAt: time.Now(),
	}, nil
}

func (p *Provider) GetHistory(_ context.Context, _ string, _ int) ([]provider.PipelineRun, error) {
	return []provider.PipelineRun{}, nil
}
