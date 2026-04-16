package drift

import (
	"fmt"

	"github.com/deepnative/engine/pkg/resource"
)

// Finding represents a drift-related insight.
type Finding struct {
	Severity   string
	ResourceID string
	Message    string
}

// Detector analyzes drift results for patterns and severity.
type Detector struct{}

// NewDetector creates a new drift detector.
func NewDetector() *Detector {
	return &Detector{}
}

// Analyze examines drift results and produces findings with severity levels.
func (d *Detector) Analyze(results []resource.DriftResult) []Finding {
	var findings []Finding

	for _, r := range results {
		severity := "info"

		switch r.DriftType {
		case resource.DriftTypeMissing:
			severity = "critical"
		case resource.DriftTypeDeleted:
			severity = "warning"
		case resource.DriftTypeModified:
			if r.Field == "state" {
				severity = "warning"
			}
		case resource.DriftTypeAdded:
			severity = "info"
		}

		// Missing tags escalate severity
		if r.Field != "" && len(r.Field) > 5 && r.Field[:5] == "tags." {
			if severity == "info" {
				severity = "warning"
			}
		}

		findings = append(findings, Finding{
			Severity:   severity,
			ResourceID: r.ResourceID,
			Message:    fmt.Sprintf("[%s] %s", r.DriftType, r.Message),
		})
	}

	return findings
}
