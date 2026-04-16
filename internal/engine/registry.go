package engine

import (
	"fmt"
	"sync"

	"github.com/deepnative/engine/pkg/provider"
)

// Registry manages provider registration and lookup.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]provider.Provider
}

// NewRegistry creates a new empty provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]provider.Provider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(name string, p provider.Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %q already registered", name)
	}
	r.providers[name] = p
	return nil
}

// Get returns a provider by name.
func (r *Registry) Get(name string) (provider.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", provider.ErrProviderNotFound, name)
	}
	return p, nil
}

// GetCloud returns a cloud provider by name.
func (r *Registry) GetCloud(name string) (provider.CloudProvider, error) {
	p, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	cp, ok := p.(provider.CloudProvider)
	if !ok {
		return nil, fmt.Errorf("%w: %s is not a cloud provider", provider.ErrUnsupportedOperation, name)
	}
	return cp, nil
}

// GetPipeline returns a pipeline provider by name.
func (r *Registry) GetPipeline(name string) (provider.PipelineProvider, error) {
	p, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	pp, ok := p.(provider.PipelineProvider)
	if !ok {
		return nil, fmt.Errorf("%w: %s is not a pipeline provider", provider.ErrUnsupportedOperation, name)
	}
	return pp, nil
}

// GetSRE returns an SRE provider by name.
func (r *Registry) GetSRE(name string) (provider.SREProvider, error) {
	p, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	sp, ok := p.(provider.SREProvider)
	if !ok {
		return nil, fmt.Errorf("%w: %s is not an SRE provider", provider.ErrUnsupportedOperation, name)
	}
	return sp, nil
}

// List returns the names of all registered providers.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// ListByKind returns providers that match a specific kind (cloud, pipeline, sre).
func (r *Registry) ListByKind(kind string) []provider.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []provider.Provider
	for _, p := range r.providers {
		if p.Metadata().Kind == kind {
			result = append(result, p)
		}
	}
	return result
}
