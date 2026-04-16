package provider

import (
	"context"
	"time"
)

// IncidentListOptions controls filtering for incident listing.
type IncidentListOptions struct {
	Status   string // "triggered", "acknowledged", "resolved"
	Urgency  string // "high", "low"
	Since    time.Time
	PageSize int
}

// Incident represents an SRE incident.
type Incident struct {
	ID          string
	Title       string
	Status      string
	Urgency     string
	Service     string
	AssignedTo  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ResolvedAt  time.Time
	Description string
}

// OnCallEntry represents who is currently on call.
type OnCallEntry struct {
	UserName   string
	UserEmail  string
	ScheduleID string
	StartAt    time.Time
	EndAt      time.Time
	Level      int
}

// SREProvider extends Provider with incident management operations.
type SREProvider interface {
	Provider
	ListIncidents(ctx context.Context, opts IncidentListOptions) ([]Incident, error)
	GetIncident(ctx context.Context, id string) (*Incident, error)
	AcknowledgeIncident(ctx context.Context, id string) error
	ResolveIncident(ctx context.Context, id string) error
	ListOnCall(ctx context.Context, scheduleID string) ([]OnCallEntry, error)
}
