package flux

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()
	if meta.Name != "flux" {
		t.Errorf("expected name flux, got %s", meta.Name)
	}
	if meta.Kind != "pipeline" {
		t.Errorf("expected kind pipeline, got %s", meta.Kind)
	}
}

func TestInitDefaultNamespace(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.namespace != "flux-system" {
		t.Errorf("expected namespace flux-system, got %s", p.namespace)
	}
}

func TestInitCustomNamespace(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"namespace": "custom-ns",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.namespace != "custom-ns" {
		t.Errorf("expected namespace custom-ns, got %s", p.namespace)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.PipelineProvider = (*Provider)(nil)
}
