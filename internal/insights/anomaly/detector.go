package anomaly

import (
	"fmt"
	"math"

	"github.com/deepnative/engine/pkg/provider"
)

// Finding represents an anomaly-related insight.
type Finding struct {
	Severity string
	Provider string
	Message  string
}

// Detector uses z-score based analysis to detect cost anomalies.
type Detector struct {
	zScoreThreshold float64
}

// NewDetector creates an anomaly detector with the given z-score threshold.
func NewDetector(zScoreThreshold float64) *Detector {
	return &Detector{zScoreThreshold: zScoreThreshold}
}

// Analyze examines cost reports for anomalous spending patterns.
func (d *Detector) Analyze(reports []*provider.CostReport) []Finding {
	var findings []Finding

	if len(reports) < 2 {
		return findings
	}

	costs := make([]float64, len(reports))
	for i, r := range reports {
		costs[i] = r.TotalCost
	}

	mean, stddev := meanAndStddev(costs)
	if stddev == 0 {
		return findings
	}

	for _, report := range reports {
		zScore := (report.TotalCost - mean) / stddev
		if math.Abs(zScore) > d.zScoreThreshold {
			severity := "warning"
			if math.Abs(zScore) > d.zScoreThreshold*1.5 {
				severity = "critical"
			}
			findings = append(findings, Finding{
				Severity: severity,
				Provider: report.Provider,
				Message:  fmt.Sprintf("cost anomaly detected: %.2f (z-score: %.2f, mean: %.2f, stddev: %.2f)", report.TotalCost, zScore, mean, stddev),
			})
		}
	}

	return findings
}

func meanAndStddev(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	var varianceSum float64
	for _, v := range values {
		diff := v - mean
		varianceSum += diff * diff
	}
	stddev := math.Sqrt(varianceSum / float64(len(values)))

	return mean, stddev
}
