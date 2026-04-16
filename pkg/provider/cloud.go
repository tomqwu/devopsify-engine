package provider

import (
	"context"
	"time"

	"github.com/deepnative/engine/pkg/resource"
)

// ListOptions controls filtering and pagination for resource listing.
type ListOptions struct {
	Region    string
	Tags      map[string]string
	PageSize  int
	PageToken string
}

// CostQueryOptions controls the time range and granularity for cost queries.
type CostQueryOptions struct {
	StartDate   time.Time
	EndDate     time.Time
	Granularity string // "daily", "monthly"
	GroupBy     string // "service", "resource", "tag"
}

// CostReport contains aggregated cost data from a cloud provider.
type CostReport struct {
	Provider    string
	StartDate   time.Time
	EndDate     time.Time
	TotalCost   float64
	Currency    string
	CostsByItem map[string]float64
}

// CloudProvider extends Provider with cloud infrastructure operations.
type CloudProvider interface {
	Provider
	ListResources(ctx context.Context, resourceType string, opts ListOptions) ([]resource.Resource, error)
	GetResource(ctx context.Context, resourceType, id string) (*resource.Resource, error)
	GetCostData(ctx context.Context, opts CostQueryOptions) (*CostReport, error)
	DetectDrift(ctx context.Context, desired []resource.Resource) ([]resource.DriftResult, error)
}
