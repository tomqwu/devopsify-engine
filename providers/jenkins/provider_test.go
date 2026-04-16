package jenkins

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()
	if meta.Name != "jenkins" {
		t.Errorf("expected name jenkins, got %s", meta.Name)
	}
	if meta.Kind != "pipeline" {
		t.Errorf("expected kind pipeline, got %s", meta.Kind)
	}
}

func TestInitRequiresServer(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"username":  "admin",
		"api_token": "token123",
	})
	if err == nil {
		t.Fatal("expected error for missing server")
	}
}

func TestInitSuccess(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"server":    "https://jenkins.example.com",
		"username":  "admin",
		"api_token": "token123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.PipelineProvider = (*Provider)(nil)
}
