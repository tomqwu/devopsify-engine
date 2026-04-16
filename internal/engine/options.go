package engine

import (
	"log"

	"github.com/deepnative/engine/internal/config"
	"github.com/deepnative/engine/internal/eventbus"
)

// Option configures the Engine.
type Option func(*Engine)

// WithConfig sets the engine configuration.
func WithConfig(cfg *config.Config) Option {
	return func(e *Engine) {
		e.config = cfg
	}
}

// WithEventBus sets the event bus for the engine.
func WithEventBus(bus *eventbus.Bus) Option {
	return func(e *Engine) {
		e.eventBus = bus
	}
}

// WithLogger sets the logger for the engine.
func WithLogger(logger *log.Logger) Option {
	return func(e *Engine) {
		e.logger = logger
	}
}
