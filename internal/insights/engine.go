package insights

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/deepnative/engine/internal/insights/anomaly"
	"github.com/deepnative/engine/internal/insights/cost"
	"github.com/deepnative/engine/internal/insights/drift"
	"github.com/deepnative/engine/internal/insights/recommendation"
	"github.com/deepnative/engine/pkg/provider"
	"github.com/deepnative/engine/pkg/resource"
)

// Finding represents a single insight finding.
type Finding struct {
	Type       string    `json:"type"`
	Severity   string    `json:"severity"`
	Provider   string    `json:"provider"`
	ResourceID string    `json:"resource_id,omitempty"`
	Message    string    `json:"message"`
	Details    any       `json:"details,omitempty"`
	DetectedAt time.Time `json:"detected_at"`
}

// AnalysisResult contains all findings from a complete analysis run.
type AnalysisResult struct {
	Findings        []Finding                       `json:"findings"`
	Recommendations []recommendation.Recommendation `json:"recommendations"`
	RunAt           time.Time                       `json:"run_at"`
	Duration        time.Duration                   `json:"duration"`
}

// Engine orchestrates all insight analyzers.
type Engine struct {
	costAnalyzer  *cost.Analyzer
	driftDetector *drift.Detector
	anomalyDetect *anomaly.Detector
	recommender   *recommendation.Recommender
	logger        *log.Logger
}

// NewEngine creates a new insights engine.
func NewEngine(costThreshold, anomalyZScore float64) *Engine {
	return &Engine{
		costAnalyzer:  cost.NewAnalyzer(costThreshold),
		driftDetector: drift.NewDetector(),
		anomalyDetect: anomaly.NewDetector(anomalyZScore),
		recommender:   recommendation.NewRecommender(),
		logger:        log.New(os.Stdout, "[insights] ", log.LstdFlags),
	}
}

// RunAnalysis executes all analyzers against the provided data.
func (e *Engine) RunAnalysis(ctx context.Context, resources []resource.Resource, costReports []*provider.CostReport, driftResults []resource.DriftResult) (*AnalysisResult, error) {
	start := time.Now()
	var findings []Finding

	// Cost analysis
	costFindings := e.costAnalyzer.Analyze(resources, costReports)
	for _, f := range costFindings {
		findings = append(findings, Finding{
			Type:       "cost",
			Severity:   f.Severity,
			Provider:   f.Provider,
			ResourceID: f.ResourceID,
			Message:    f.Message,
			DetectedAt: time.Now(),
		})
	}

	// Drift analysis
	driftFindings := e.driftDetector.Analyze(driftResults)
	for _, f := range driftFindings {
		findings = append(findings, Finding{
			Type:       "drift",
			Severity:   f.Severity,
			ResourceID: f.ResourceID,
			Message:    f.Message,
			DetectedAt: time.Now(),
		})
	}

	// Anomaly detection
	anomalyFindings := e.anomalyDetect.Analyze(costReports)
	for _, f := range anomalyFindings {
		findings = append(findings, Finding{
			Type:       "anomaly",
			Severity:   f.Severity,
			Provider:   f.Provider,
			Message:    f.Message,
			DetectedAt: time.Now(),
		})
	}

	// Generate recommendations
	recFindings := make([]recommendation.Finding, len(findings))
	for i, f := range findings {
		recFindings[i] = recommendation.Finding{
			Type:     f.Type,
			Severity: f.Severity,
			Message:  f.Message,
		}
	}
	recs := e.recommender.Generate(recFindings)

	e.logger.Printf("analysis complete: %d findings, %d recommendations in %v", len(findings), len(recs), time.Since(start))

	return &AnalysisResult{
		Findings:        findings,
		Recommendations: recs,
		RunAt:           start,
		Duration:        time.Since(start),
	}, nil
}
