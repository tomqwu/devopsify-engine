package aws

import (
	"context"

	"github.com/deepnative/engine/pkg/provider"
)

// getCostExplorerData retrieves cost data from AWS Cost Explorer.
// This is a scaffold - actual AWS SDK calls will replace the placeholder.
func getCostExplorerData(_ context.Context, _ *AuthConfig, opts provider.CostQueryOptions) (*provider.CostReport, error) {
	return &provider.CostReport{
		Provider:    "aws",
		StartDate:   opts.StartDate,
		EndDate:     opts.EndDate,
		TotalCost:   0,
		Currency:    "USD",
		CostsByItem: make(map[string]float64),
	}, nil
}
