package argocd

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()
	if meta.Name != "argocd" {
		t.Errorf("expected name argocd, got %s", meta.Name)
	}
	if meta.Kind != "pipeline" {
		t.Errorf("expected kind pipeline, got %s", meta.Kind)
	}
}

func TestInitRequiresServer(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"token": "test-token",
	})
	if err == nil {
		t.Fatal("expected error for missing server")
	}
}

func TestInitRequiresToken(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"server": "https://argocd.example.com",
	})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestInitSuccess(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"server": "https://argocd.example.com",
		"token":  "test-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.PipelineProvider = (*Provider)(nil)
}
