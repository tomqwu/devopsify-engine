package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/deepnative/engine/internal/api/handlers"
	"github.com/deepnative/engine/internal/api/middleware"
	"github.com/deepnative/engine/internal/engine"
	"github.com/deepnative/engine/internal/insights"
)

// Server is the HTTP API server for the Deep Native Engine.
type Server struct {
	httpServer     *http.Server
	engine         *engine.Engine
	insightsEngine *insights.Engine
	logger         *log.Logger
}

// NewServer creates a new API server.
func NewServer(addr string, eng *engine.Engine, insightsEng *insights.Engine, readTimeout, writeTimeout int) *Server {
	s := &Server{
		engine:         eng,
		insightsEngine: insightsEng,
		logger:         log.New(os.Stdout, "[api] ", log.LstdFlags),
	}

	mux := http.NewServeMux()
	s.registerRoutes(mux)

	handler := middleware.Recovery(
		middleware.Logging(
			middleware.CORS(mux),
		),
	)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}

	return s
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Health checks
	health := handlers.NewHealthHandler(s.engine)
	mux.HandleFunc("GET /healthz", health.Healthz)
	mux.HandleFunc("GET /readyz", health.Readyz)

	// Cloud resources
	cloud := handlers.NewCloudHandler(s.engine)
	mux.HandleFunc("GET /api/v1/providers", cloud.ListProviders)
	mux.HandleFunc("GET /api/v1/resources", cloud.ListResources)
	mux.HandleFunc("GET /api/v1/costs", cloud.GetCosts)

	// Pipelines
	pipeline := handlers.NewPipelineHandler(s.engine)
	mux.HandleFunc("GET /api/v1/pipelines", pipeline.ListPipelines)
	mux.HandleFunc("POST /api/v1/pipelines/{id}/sync", pipeline.TriggerSync)

	// Insights
	insightsHandler := handlers.NewInsightsHandler(s.insightsEngine)
	mux.HandleFunc("GET /api/v1/insights", insightsHandler.GetInsights)
	mux.HandleFunc("GET /api/v1/insights/cost", insightsHandler.GetCostInsights)
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	s.logger.Printf("starting API server on %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("shutting down API server")
	return s.httpServer.Shutdown(ctx)
}
