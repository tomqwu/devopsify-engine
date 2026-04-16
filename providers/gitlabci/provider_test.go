package gitlabci

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()
	if meta.Name != "gitlab-ci" {
		t.Errorf("expected name gitlab-ci, got %s", meta.Name)
	}
	if meta.Kind != "pipeline" {
		t.Errorf("expected kind pipeline, got %s", meta.Kind)
	}
}

func TestInitRequiresProjectID(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"token": "glpat-test",
	})
	if err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestInitRequiresToken(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"project_id": "12345",
	})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestInitSuccess(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"project_id": "12345",
		"token":      "glpat-test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.server != "https://gitlab.com" {
		t.Errorf("expected default server, got %s", p.server)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.PipelineProvider = (*Provider)(nil)
}
