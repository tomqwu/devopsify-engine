package githubactions

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()
	if meta.Name != "github-actions" {
		t.Errorf("expected name github-actions, got %s", meta.Name)
	}
	if meta.Kind != "pipeline" {
		t.Errorf("expected kind pipeline, got %s", meta.Kind)
	}
}

func TestInitRequiresOwner(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"repo":  "my-repo",
		"token": "ghp_test",
	})
	if err == nil {
		t.Fatal("expected error for missing owner")
	}
}

func TestInitRequiresRepo(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"owner": "my-org",
		"token": "ghp_test",
	})
	if err == nil {
		t.Fatal("expected error for missing repo")
	}
}

func TestInitSuccess(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"owner": "my-org",
		"repo":  "my-repo",
		"token": "ghp_test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.PipelineProvider = (*Provider)(nil)
}
