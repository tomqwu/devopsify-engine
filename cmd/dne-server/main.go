package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deepnative/engine/internal/api"
	"github.com/deepnative/engine/internal/config"
	"github.com/deepnative/engine/internal/engine"
	"github.com/deepnative/engine/internal/insights"
	"github.com/deepnative/engine/providers/argocd"
	"github.com/deepnative/engine/providers/aws"
	"github.com/deepnative/engine/providers/azure"
	"github.com/deepnative/engine/providers/flux"
	"github.com/deepnative/engine/providers/gcp"
	"github.com/deepnative/engine/providers/githubactions"
	"github.com/deepnative/engine/providers/gitlabci"
	"github.com/deepnative/engine/providers/jenkins"
	"github.com/deepnative/engine/providers/opsgenie"
	"github.com/deepnative/engine/providers/pagerduty"
	"github.com/deepnative/engine/pkg/provider"
)

var providerFactory = map[string]func() provider.Provider{
	"aws":            func() provider.Provider { return &aws.Provider{} },
	"azure":          func() provider.Provider { return &azure.Provider{} },
	"gcp":            func() provider.Provider { return &gcp.Provider{} },
	"argocd":         func() provider.Provider { return &argocd.Provider{} },
	"flux":           func() provider.Provider { return &flux.Provider{} },
	"github-actions": func() provider.Provider { return &githubactions.Provider{} },
	"gitlab-ci":      func() provider.Provider { return &gitlabci.Provider{} },
	"jenkins":        func() provider.Provider { return &jenkins.Provider{} },
	"pagerduty":      func() provider.Provider { return &pagerduty.Provider{} },
	"opsgenie":       func() provider.Provider { return &opsgenie.Provider{} },
}

func main() {
	configPath := os.Getenv("DNE_CONFIG")
	if configPath == "" {
		configPath = "configs/default.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create engine
	eng := engine.New(engine.WithConfig(cfg))
	ctx := context.Background()

	// Register providers from config
	for _, pc := range cfg.Providers {
		factory, ok := providerFactory[pc.Provider]
		if !ok {
			log.Printf("unknown provider type: %s", pc.Provider)
			continue
		}
		p := factory()
		if err := eng.RegisterProvider(ctx, pc.Name, p, pc.Config); err != nil {
			log.Printf("failed to register provider %s: %v", pc.Name, err)
			continue
		}
	}

	// Start engine
	if err := eng.Start(ctx); err != nil {
		log.Fatalf("failed to start engine: %v", err)
	}

	// Create insights engine
	insightsEng := insights.NewEngine(cfg.Insights.CostThreshold, cfg.Insights.AnomalyZScore)

	// Create and start API server
	srv := api.NewServer(cfg.API.Address, eng, insightsEng, cfg.API.ReadTimeout, cfg.API.WriteTimeout)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.API.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	if err := eng.Shutdown(shutdownCtx); err != nil {
		log.Printf("engine shutdown error: %v", err)
	}

	log.Println("shutdown complete")
}
