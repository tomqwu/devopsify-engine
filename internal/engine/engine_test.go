package engine

import (
	"context"
	"testing"

	"github.com/deepnative/engine/internal/config"
	"github.com/deepnative/engine/internal/eventbus"
	"github.com/deepnative/engine/pkg/provider"
)

type mockProvider struct {
	meta        provider.Metadata
	initErr     error
	healthyErr  error
	shutdownErr error
}

func (m *mockProvider) Metadata() provider.Metadata             { return m.meta }
func (m *mockProvider) Init(_ context.Context, _ map[string]any) error { return m.initErr }
func (m *mockProvider) Healthy(_ context.Context) error         { return m.healthyErr }
func (m *mockProvider) Shutdown(_ context.Context) error        { return m.shutdownErr }

func TestNewEngine(t *testing.T) {
	cfg := config.Default()
	bus := eventbus.New()

	e := New(WithConfig(cfg), WithEventBus(bus))

	if e.Config() != cfg {
		t.Error("expected config to be set")
	}
	if e.EventBus() != bus {
		t.Error("expected event bus to be set")
	}
	if e.Registry() == nil {
		t.Error("expected registry to be non-nil")
	}
}

func TestRegisterProvider(t *testing.T) {
	e := New()
	ctx := context.Background()

	mock := &mockProvider{
		meta: provider.Metadata{Name: "test-aws", Kind: "cloud", Version: "1.0.0"},
	}

	err := e.RegisterProvider(ctx, "test-aws", mock, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, err := e.Registry().Get("test-aws")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Metadata().Name != "test-aws" {
		t.Errorf("expected name test-aws, got %s", p.Metadata().Name)
	}
}

func TestRegisterProviderInitError(t *testing.T) {
	e := New()
	ctx := context.Background()

	mock := &mockProvider{
		meta:    provider.Metadata{Name: "bad", Kind: "cloud"},
		initErr: provider.ErrProviderUnhealthy,
	}

	err := e.RegisterProvider(ctx, "bad", mock, nil)
	if err == nil {
		t.Fatal("expected error from Init")
	}
}

func TestRegisterDuplicateProvider(t *testing.T) {
	e := New()
	ctx := context.Background()

	mock := &mockProvider{
		meta: provider.Metadata{Name: "dup", Kind: "cloud"},
	}

	if err := e.RegisterProvider(ctx, "dup", mock, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := e.RegisterProvider(ctx, "dup", mock, nil)
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestStartAndShutdown(t *testing.T) {
	e := New()
	ctx := context.Background()

	mock := &mockProvider{
		meta: provider.Metadata{Name: "svc", Kind: "pipeline"},
	}
	if err := e.RegisterProvider(ctx, "svc", mock, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := e.Start(ctx); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}
}
