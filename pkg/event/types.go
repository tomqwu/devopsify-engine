package event

import (
	"github.com/deepnative/engine/pkg/resource"
)

// Event type constants.
const (
	TypeProviderRegistered = "provider.registered"
	TypeProviderHealthy    = "provider.healthy"
	TypeProviderUnhealthy  = "provider.unhealthy"
	TypeResourceDiscovered = "resource.discovered"
	TypeDriftDetected      = "drift.detected"
	TypeCostReportReady    = "cost.report.ready"
	TypePipelineSynced     = "pipeline.synced"
	TypePipelineFailed     = "pipeline.failed"
	TypeIncidentTriggered  = "incident.triggered"
	TypeIncidentResolved   = "incident.resolved"
	TypeInsightGenerated   = "insight.generated"
	TypeEngineStarted      = "engine.started"
	TypeEngineStopped      = "engine.stopped"
)

// ProviderPayload carries provider lifecycle event data.
type ProviderPayload struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

// ResourcePayload carries resource discovery event data.
type ResourcePayload struct {
	Provider  string              `json:"provider"`
	Count     int                 `json:"count"`
	Resources []resource.Resource `json:"resources"`
}

// DriftPayload carries drift detection event data.
type DriftPayload struct {
	Provider string                 `json:"provider"`
	Count    int                    `json:"count"`
	Results  []resource.DriftResult `json:"results"`
}

// CostPayload carries cost report event data.
type CostPayload struct {
	Provider  string  `json:"provider"`
	TotalCost float64 `json:"total_cost"`
	Currency  string  `json:"currency"`
}

// PipelinePayload carries pipeline event data.
type PipelinePayload struct {
	Provider   string `json:"provider"`
	PipelineID string `json:"pipeline_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

// IncidentPayload carries incident event data.
type IncidentPayload struct {
	Provider   string `json:"provider"`
	IncidentID string `json:"incident_id"`
	Title      string `json:"title"`
	Status     string `json:"status"`
}

// InsightPayload carries insight generation event data.
type InsightPayload struct {
	Type    string `json:"type"`
	Count   int    `json:"count"`
	Summary string `json:"summary"`
}
