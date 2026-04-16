package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deepnative/engine/internal/engine"
	"github.com/deepnative/engine/pkg/provider"
	"github.com/deepnative/engine/pkg/resource"
)

// CloudHandler handles cloud resource endpoints.
type CloudHandler struct {
	engine *engine.Engine
}

// NewCloudHandler creates a new cloud handler.
func NewCloudHandler(eng *engine.Engine) *CloudHandler {
	return &CloudHandler{engine: eng}
}

// ListProviders returns all registered providers.
func (h *CloudHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	names := h.engine.Registry().List()

	type providerInfo struct {
		Name    string `json:"name"`
		Kind    string `json:"kind"`
		Version string `json:"version"`
	}

	result := make([]providerInfo, 0, len(names))
	for _, name := range names {
		p, err := h.engine.Registry().Get(name)
		if err != nil {
			continue
		}
		meta := p.Metadata()
		result = append(result, providerInfo{
			Name:    meta.Name,
			Kind:    meta.Kind,
			Version: meta.Version,
		})
	}

	writeJSON(w, http.StatusOK, result)
}

// ListResources returns resources from all cloud providers.
func (h *CloudHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	cloudProviders := h.engine.Registry().ListByKind("cloud")
	resourceType := r.URL.Query().Get("type")
	region := r.URL.Query().Get("region")

	opts := provider.ListOptions{
		Region: region,
	}

	var allResources []resource.Resource
	for _, p := range cloudProviders {
		cp, ok := p.(provider.CloudProvider)
		if !ok {
			continue
		}
		resources, err := cp.ListResources(r.Context(), resourceType, opts)
		if err != nil {
			continue
		}
		allResources = append(allResources, resources...)
	}

	writeJSON(w, http.StatusOK, allResources)
}

// GetCosts returns cost data from all cloud providers.
func (h *CloudHandler) GetCosts(w http.ResponseWriter, r *http.Request) {
	cloudProviders := h.engine.Registry().ListByKind("cloud")

	opts := provider.CostQueryOptions{
		StartDate:   time.Now().AddDate(0, -1, 0),
		EndDate:     time.Now(),
		Granularity: r.URL.Query().Get("granularity"),
		GroupBy:     r.URL.Query().Get("group_by"),
	}
	if opts.Granularity == "" {
		opts.Granularity = "daily"
	}

	var reports []*provider.CostReport
	for _, p := range cloudProviders {
		cp, ok := p.(provider.CloudProvider)
		if !ok {
			continue
		}
		report, err := cp.GetCostData(r.Context(), opts)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}

	writeJSON(w, http.StatusOK, reports)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data) //nolint:errcheck
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
