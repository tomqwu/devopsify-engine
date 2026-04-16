package engine

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deepnative/engine/internal/config"
	"github.com/deepnative/engine/internal/eventbus"
	"github.com/deepnative/engine/pkg/event"
	"github.com/deepnative/engine/pkg/provider"
)

// Engine is the core orchestrator of the Deep Native Engine.
type Engine struct {
	config   *config.Config
	registry *Registry
	eventBus *eventbus.Bus
	logger   *log.Logger
}

// New creates a new Engine with the given options.
func New(opts ...Option) *Engine {
	e := &Engine{
		registry: NewRegistry(),
		eventBus: eventbus.New(),
		logger:   log.New(os.Stdout, "[engine] ", log.LstdFlags),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Registry returns the provider registry.
func (e *Engine) Registry() *Registry {
	return e.registry
}

// EventBus returns the event bus.
func (e *Engine) EventBus() *eventbus.Bus {
	return e.eventBus
}

// Config returns the engine configuration.
func (e *Engine) Config() *config.Config {
	return e.config
}

// RegisterProvider initializes and registers a provider with the engine.
func (e *Engine) RegisterProvider(ctx context.Context, name string, p provider.Provider, cfg map[string]any) error {
	if err := p.Init(ctx, cfg); err != nil {
		return fmt.Errorf("%w: %s: %v", provider.ErrProviderInit, name, err)
	}

	if err := e.registry.Register(name, p); err != nil {
		return err
	}

	meta := p.Metadata()
	e.eventBus.Publish(event.New(event.TypeProviderRegistered, "engine", &event.ProviderPayload{
		Name:    name,
		Kind:    meta.Kind,
		Message: fmt.Sprintf("provider %s (%s) registered", name, meta.Kind),
	}))

	e.logger.Printf("registered provider: %s (kind=%s)", name, meta.Kind)
	return nil
}

// Start starts the engine and publishes an engine started event.
func (e *Engine) Start(ctx context.Context) error {
	e.logger.Println("starting engine")
	e.eventBus.Publish(event.New(event.TypeEngineStarted, "engine", nil))
	return nil
}

// Shutdown gracefully shuts down all registered providers.
func (e *Engine) Shutdown(ctx context.Context) error {
	e.logger.Println("shutting down engine")

	var errs []error
	for _, name := range e.registry.List() {
		p, err := e.registry.Get(name)
		if err != nil {
			continue
		}
		if err := p.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown %s: %w", name, err))
		}
	}

	e.eventBus.Publish(event.New(event.TypeEngineStopped, "engine", nil))

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}
