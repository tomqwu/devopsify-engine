package eventbus

import (
	"sync"
	"testing"

	"github.com/deepnative/engine/pkg/event"
)

func TestSubscribeAndPublish(t *testing.T) {
	bus := New()
	var received event.Event

	bus.Subscribe("test.event", func(e event.Event) {
		received = e
	})

	evt := event.New("test.event", "test", "payload")
	bus.Publish(evt)

	if received == nil {
		t.Fatal("expected event to be received")
	}
	if received.EventType() != "test.event" {
		t.Errorf("expected type test.event, got %s", received.EventType())
	}
	if received.EventSource() != "test" {
		t.Errorf("expected source test, got %s", received.EventSource())
	}
}

func TestMultipleSubscribers(t *testing.T) {
	bus := New()
	count := 0
	mu := sync.Mutex{}

	for range 3 {
		bus.Subscribe("multi", func(_ event.Event) {
			mu.Lock()
			count++
			mu.Unlock()
		})
	}

	bus.Publish(event.New("multi", "test", nil))

	mu.Lock()
	defer mu.Unlock()
	if count != 3 {
		t.Errorf("expected 3 handlers called, got %d", count)
	}
}

func TestNoSubscribers(t *testing.T) {
	bus := New()
	// Should not panic
	bus.Publish(event.New("nobody.listening", "test", nil))
}

func TestSubscriberCount(t *testing.T) {
	bus := New()

	if bus.SubscriberCount("test") != 0 {
		t.Error("expected 0 subscribers")
	}

	bus.Subscribe("test", func(_ event.Event) {})
	bus.Subscribe("test", func(_ event.Event) {})

	if bus.SubscriberCount("test") != 2 {
		t.Errorf("expected 2 subscribers, got %d", bus.SubscriberCount("test"))
	}
}

func TestDifferentEventTypes(t *testing.T) {
	bus := New()
	var typeA, typeB bool

	bus.Subscribe("type.a", func(_ event.Event) { typeA = true })
	bus.Subscribe("type.b", func(_ event.Event) { typeB = true })

	bus.Publish(event.New("type.a", "test", nil))

	if !typeA {
		t.Error("expected type.a handler to be called")
	}
	if typeB {
		t.Error("expected type.b handler to NOT be called")
	}
}
