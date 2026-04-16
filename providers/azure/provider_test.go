package azure

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()

	if meta.Name != "azure" {
		t.Errorf("expected name azure, got %s", meta.Name)
	}
	if meta.Kind != "cloud" {
		t.Errorf("expected kind cloud, got %s", meta.Kind)
	}
}

func TestInitRequiresSubscription(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing subscription_id")
	}
}

func TestInitSuccess(t *testing.T) {
	p := &Provider{}
	err := p.Init(context.Background(), map[string]any{
		"subscription_id": "sub-123",
		"tenant_id":       "tenant-456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.CloudProvider = (*Provider)(nil)
}
