package provider

import "context"

// Metadata describes a provider's identity and capabilities.
type Metadata struct {
	Name    string
	Kind    string // "cloud", "pipeline", "sre"
	Version string
}

// Provider is the base lifecycle interface all providers must implement.
type Provider interface {
	Metadata() Metadata
	Init(ctx context.Context, config map[string]any) error
	Healthy(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
