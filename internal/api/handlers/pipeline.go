package handlers

import (
	"net/http"

	"github.com/deepnative/engine/internal/engine"
	"github.com/deepnative/engine/pkg/provider"
)

// PipelineHandler handles pipeline endpoints.
type PipelineHandler struct {
	engine *engine.Engine
}

// NewPipelineHandler creates a new pipeline handler.
func NewPipelineHandler(eng *engine.Engine) *PipelineHandler {
	return &PipelineHandler{engine: eng}
}

// ListPipelines returns pipelines from all pipeline providers.
func (h *PipelineHandler) ListPipelines(w http.ResponseWriter, r *http.Request) {
	pipelineProviders := h.engine.Registry().ListByKind("pipeline")
	namespace := r.URL.Query().Get("namespace")

	opts := provider.PipelineListOptions{
		Namespace: namespace,
	}

	var allPipelines []provider.Pipeline
	for _, p := range pipelineProviders {
		pp, ok := p.(provider.PipelineProvider)
		if !ok {
			continue
		}
		pipelines, err := pp.ListPipelines(r.Context(), opts)
		if err != nil {
			continue
		}
		allPipelines = append(allPipelines, pipelines...)
	}

	writeJSON(w, http.StatusOK, allPipelines)
}

// TriggerSync triggers a sync for a specific pipeline.
func (h *PipelineHandler) TriggerSync(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "pipeline id is required")
		return
	}

	providerName := r.URL.Query().Get("provider")
	if providerName == "" {
		writeError(w, http.StatusBadRequest, "provider query parameter is required")
		return
	}

	pp, err := h.engine.Registry().GetPipeline(providerName)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	result, err := pp.TriggerSync(r.Context(), id, provider.SyncOptions{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
