package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/deepnative/engine/internal/engine"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	engine *engine.Engine
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(eng *engine.Engine) *HealthHandler {
	return &HealthHandler{engine: eng}
}

// Healthz is a liveness probe.
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz is a readiness probe that checks all provider health.
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	providerNames := h.engine.Registry().List()
	unhealthy := []string{}

	for _, name := range providerNames {
		p, err := h.engine.Registry().Get(name)
		if err != nil {
			unhealthy = append(unhealthy, name)
			continue
		}
		if err := p.Healthy(ctx); err != nil {
			unhealthy = append(unhealthy, name)
		}
	}

	if len(unhealthy) > 0 {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status":    "degraded",
			"unhealthy": unhealthy,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
