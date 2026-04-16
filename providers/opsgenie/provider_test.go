package opsgenie

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()
	if meta.Name != "opsgenie" {
		t.Errorf("expected name opsgenie, got %s", meta.Name)
	}
	if meta.Kind != "sre" {
		t.Errorf("expected kind sre, got %s", meta.Kind)
	}
}

func TestInitRequiresAPIKey(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestInitSuccess(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"api_key": "test-key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.SREProvider = (*Provider)(nil)
}
