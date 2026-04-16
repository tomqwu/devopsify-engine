package eventbus

import (
	"sync"

	"github.com/deepnative/engine/pkg/event"
)

// Handler is a function that processes events.
type Handler func(event.Event)

// Bus is a synchronous, in-process event bus.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

// New creates a new event bus.
func New() *Bus {
	return &Bus{
		subscribers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for a specific event type.
func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

// Publish sends an event to all registered handlers for that event type.
func (b *Bus) Publish(evt event.Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.subscribers[evt.EventType()]))
	copy(handlers, b.subscribers[evt.EventType()])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(evt)
	}
}

// SubscriberCount returns the number of subscribers for a given event type.
func (b *Bus) SubscriberCount(eventType string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.subscribers[eventType])
}
