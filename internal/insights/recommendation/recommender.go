package recommendation

// Recommendation is an actionable suggestion generated from findings.
type Recommendation struct {
	Type        string `json:"type"`
	Priority    string `json:"priority"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

// Finding is the subset of insight data needed for recommendations.
type Finding struct {
	Type     string
	Severity string
	Message  string
}

// Recommender generates actionable recommendations from findings.
type Recommender struct{}

// NewRecommender creates a new recommender.
func NewRecommender() *Recommender {
	return &Recommender{}
}

// Generate produces recommendations from analysis findings.
func (r *Recommender) Generate(findings []Finding) []Recommendation {
	var recs []Recommendation

	costCount := 0
	driftCount := 0
	anomalyCount := 0
	criticalCount := 0

	for _, f := range findings {
		switch f.Type {
		case "cost":
			costCount++
		case "drift":
			driftCount++
		case "anomaly":
			anomalyCount++
		}
		if f.Severity == "critical" {
			criticalCount++
		}
	}

	if costCount > 0 {
		recs = append(recs, Recommendation{
			Type:        "cost",
			Priority:    priorityFromCount(costCount),
			Title:       "Review Cloud Spending",
			Description: "Cost analysis detected spending patterns that may need attention.",
			Action:      "Review cost reports and consider rightsizing or terminating idle resources.",
		})
	}

	if driftCount > 0 {
		recs = append(recs, Recommendation{
			Type:        "drift",
			Priority:    priorityFromCount(driftCount),
			Title:       "Resolve Configuration Drift",
			Description: "Drift detection found discrepancies between desired and actual state.",
			Action:      "Run reconciliation to align actual state with desired configuration.",
		})
	}

	if anomalyCount > 0 {
		recs = append(recs, Recommendation{
			Type:        "anomaly",
			Priority:    "high",
			Title:       "Investigate Cost Anomalies",
			Description: "Unusual spending patterns detected that deviate from normal behavior.",
			Action:      "Review recent changes and investigate root cause of cost spikes.",
		})
	}

	if criticalCount > 0 {
		recs = append(recs, Recommendation{
			Type:        "general",
			Priority:    "critical",
			Title:       "Address Critical Findings",
			Description: "Critical issues detected that require immediate attention.",
			Action:      "Review and resolve all critical findings before they impact operations.",
		})
	}

	return recs
}

func priorityFromCount(count int) string {
	switch {
	case count >= 10:
		return "critical"
	case count >= 5:
		return "high"
	case count >= 2:
		return "medium"
	default:
		return "low"
	}
}
