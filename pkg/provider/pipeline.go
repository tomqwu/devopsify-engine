package provider

import (
	"context"
	"time"
)

// PipelineListOptions controls filtering for pipeline listing.
type PipelineListOptions struct {
	Namespace string
	Labels    map[string]string
}

// Pipeline represents a deployment pipeline.
type Pipeline struct {
	ID        string
	Name      string
	Namespace string
	Status    string
	Source    string
	Labels    map[string]string
}

// PipelineStatus represents the current state of a pipeline.
type PipelineStatus struct {
	ID          string
	Name        string
	Phase       string // "synced", "progressing", "degraded", "failed", "unknown"
	Health      string
	Message     string
	LastSyncAt  time.Time
	SyncedAt    time.Time
}

// SyncOptions controls how a pipeline sync is triggered.
type SyncOptions struct {
	Prune   bool
	DryRun  bool
	Force   bool
	Timeout time.Duration
}

// SyncResult contains the outcome of a pipeline sync operation.
type SyncResult struct {
	ID        string
	Status    string
	Message   string
	StartedAt time.Time
}

// PipelineRun represents a historical pipeline execution.
type PipelineRun struct {
	ID         string
	PipelineID string
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
	Trigger    string
	Message    string
}

// PipelineProvider extends Provider with CI/CD pipeline operations.
type PipelineProvider interface {
	Provider
	ListPipelines(ctx context.Context, opts PipelineListOptions) ([]Pipeline, error)
	GetPipelineStatus(ctx context.Context, id string) (*PipelineStatus, error)
	TriggerSync(ctx context.Context, id string, opts SyncOptions) (*SyncResult, error)
	GetHistory(ctx context.Context, id string, limit int) ([]PipelineRun, error)
}
