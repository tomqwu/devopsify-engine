package gcp

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()

	if meta.Name != "gcp" {
		t.Errorf("expected name gcp, got %s", meta.Name)
	}
	if meta.Kind != "cloud" {
		t.Errorf("expected kind cloud, got %s", meta.Kind)
	}
}

func TestInitRequiresProject(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestInitSuccess(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"project": "my-project",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.project != "my-project" {
		t.Errorf("expected project my-project, got %s", p.project)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.CloudProvider = (*Provider)(nil)
}
