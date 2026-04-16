package event

import "time"

// Event is the interface for all events in the system.
type Event interface {
	EventType() string
	EventTime() time.Time
	EventSource() string
	EventPayload() any
}

// Base is a concrete implementation of the Event interface.
type Base struct {
	Type    string    `json:"type"`
	Time    time.Time `json:"time"`
	Source  string    `json:"source"`
	Payload any       `json:"payload"`
}

func (e *Base) EventType() string    { return e.Type }
func (e *Base) EventTime() time.Time { return e.Time }
func (e *Base) EventSource() string  { return e.Source }
func (e *Base) EventPayload() any    { return e.Payload }

// New creates a new event with the given type, source, and payload.
func New(eventType, source string, payload any) Event {
	return &Base{
		Type:    eventType,
		Time:    time.Now(),
		Source:  source,
		Payload: payload,
	}
}
