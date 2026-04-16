package handlers

import (
	"net/http"

	"github.com/deepnative/engine/internal/insights"
)

// InsightsHandler handles insights endpoints.
type InsightsHandler struct {
	engine *insights.Engine
}

// NewInsightsHandler creates a new insights handler.
func NewInsightsHandler(eng *insights.Engine) *InsightsHandler {
	return &InsightsHandler{engine: eng}
}

// GetInsights returns a full analysis result.
func (h *InsightsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	result, err := h.engine.RunAnalysis(r.Context(), nil, nil, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetCostInsights returns cost-specific insights.
func (h *InsightsHandler) GetCostInsights(w http.ResponseWriter, r *http.Request) {
	result, err := h.engine.RunAnalysis(r.Context(), nil, nil, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var costFindings []insights.Finding
	for _, f := range result.Findings {
		if f.Type == "cost" {
			costFindings = append(costFindings, f)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"findings": costFindings,
		"run_at":   result.RunAt,
	})
}
