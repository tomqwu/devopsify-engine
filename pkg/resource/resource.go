package resource

import "time"

// State represents the lifecycle state of a resource.
type State string

const (
	StateRunning      State = "running"
	StateStopped      State = "stopped"
	StateTerminated   State = "terminated"
	StatePending      State = "pending"
	StateUnknown      State = "unknown"
	StateDegraded     State = "degraded"
	StateHealthy      State = "healthy"
	StateUnhealthy    State = "unhealthy"
)

// DriftType classifies the kind of drift detected.
type DriftType string

const (
	DriftTypeModified DriftType = "modified"
	DriftTypeAdded    DriftType = "added"
	DriftTypeDeleted  DriftType = "deleted"
	DriftTypeMissing  DriftType = "missing"
)

// Resource is the universal resource model that all providers normalize to.
type Resource struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Provider     string            `json:"provider"`
	Region       string            `json:"region,omitempty"`
	State        State             `json:"state"`
	Tags         map[string]string `json:"tags,omitempty"`
	Properties   map[string]any    `json:"properties,omitempty"`
	CostPerMonth float64           `json:"cost_per_month,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// DriftResult represents a detected drift between desired and actual state.
type DriftResult struct {
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	DriftType    DriftType `json:"drift_type"`
	Field        string    `json:"field"`
	Expected     string    `json:"expected"`
	Actual       string    `json:"actual"`
	Message      string    `json:"message"`
}
