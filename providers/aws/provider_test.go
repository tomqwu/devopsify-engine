package aws

import (
	"context"
	"testing"

	"github.com/deepnative/engine/pkg/provider"
)

func TestMetadata(t *testing.T) {
	p := &Provider{}
	meta := p.Metadata()

	if meta.Name != "aws" {
		t.Errorf("expected name aws, got %s", meta.Name)
	}
	if meta.Kind != "cloud" {
		t.Errorf("expected kind cloud, got %s", meta.Kind)
	}
	if meta.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", meta.Version)
	}
}

func TestInit(t *testing.T) {
	p := &Provider{}
	ctx := context.Background()

	err := p.Init(ctx, map[string]any{
		"region": "us-west-2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.region != "us-west-2" {
		t.Errorf("expected region us-west-2, got %s", p.region)
	}
}

func TestHealthy(t *testing.T) {
	p := &Provider{}

	// Unhealthy before init
	if err := p.Healthy(context.Background()); err == nil {
		t.Error("expected unhealthy before init")
	}

	// Healthy after init
	_ = p.Init(context.Background(), map[string]any{"region": "us-east-1"})
	if err := p.Healthy(context.Background()); err != nil {
		t.Errorf("expected healthy after init: %v", err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ provider.CloudProvider = (*Provider)(nil)
}

func TestMapEC2State(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"running", "running"},
		{"stopped", "stopped"},
		{"terminated", "terminated"},
		{"pending", "pending"},
		{"unknown-state", "unknown"},
	}

	for _, tt := range tests {
		got := mapEC2State(tt.input)
		if string(got) != tt.want {
			t.Errorf("mapEC2State(%s) = %s, want %s", tt.input, got, tt.want)
		}
	}
}
