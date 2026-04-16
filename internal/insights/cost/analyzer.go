package cost

import (
	"fmt"

	"github.com/deepnative/engine/pkg/provider"
	"github.com/deepnative/engine/pkg/resource"
)

// Finding represents a cost-related insight.
type Finding struct {
	Severity   string
	Provider   string
	ResourceID string
	Message    string
}

// Analyzer detects cost-related insights such as spend concentration and idle resources.
type Analyzer struct {
	threshold float64
}

// NewAnalyzer creates a cost analyzer with the given threshold.
func NewAnalyzer(threshold float64) *Analyzer {
	return &Analyzer{threshold: threshold}
}

// Analyze examines resources and cost reports for cost insights.
func (a *Analyzer) Analyze(resources []resource.Resource, reports []*provider.CostReport) []Finding {
	var findings []Finding

	// Check for spend concentration in cost reports
	for _, report := range reports {
		if report.TotalCost > a.threshold {
			findings = append(findings, Finding{
				Severity: "warning",
				Provider: report.Provider,
				Message:  fmt.Sprintf("total spend %.2f %s exceeds threshold %.2f", report.TotalCost, report.Currency, a.threshold),
			})
		}

		// Detect spend concentration (single item > 80% of total)
		for item, itemCost := range report.CostsByItem {
			if report.TotalCost > 0 && itemCost/report.TotalCost > 0.8 {
				findings = append(findings, Finding{
					Severity: "warning",
					Provider: report.Provider,
					Message:  fmt.Sprintf("spend concentration: %s accounts for %.1f%% of total spend", item, (itemCost/report.TotalCost)*100),
				})
			}
		}
	}

	// Check for potentially idle resources (stopped but still costing money)
	for _, r := range resources {
		if r.State == resource.StateStopped && r.CostPerMonth > 0 {
			findings = append(findings, Finding{
				Severity:   "info",
				Provider:   r.Provider,
				ResourceID: r.ID,
				Message:    fmt.Sprintf("resource %s is stopped but has cost %.2f/month", r.Name, r.CostPerMonth),
			})
		}
	}

	return findings
}
