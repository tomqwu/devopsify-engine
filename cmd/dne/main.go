package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

const version = "0.1.0"

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
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("dne version %s\n", version)

	case "providers":
		fmt.Println("Available providers:")
		for name := range providerFactory {
			fmt.Printf("  - %s\n", name)
		}

	case "resources":
		runWithEngine(func(eng *engine.Engine) {
			ctx := context.Background()
			for _, name := range eng.Registry().List() {
				cp, err := eng.Registry().GetCloud(name)
				if err != nil {
					continue
				}
				resources, err := cp.ListResources(ctx, "", provider.ListOptions{})
				if err != nil {
					log.Printf("error listing resources from %s: %v", name, err)
					continue
				}
				fmt.Printf("\n[%s] %d resources:\n", name, len(resources))
				for _, r := range resources {
					fmt.Printf("  %s (%s) - %s [%s]\n", r.Name, r.ID, r.Type, r.State)
				}
			}
		})

	case "pipelines":
		runWithEngine(func(eng *engine.Engine) {
			ctx := context.Background()
			for _, name := range eng.Registry().List() {
				pp, err := eng.Registry().GetPipeline(name)
				if err != nil {
					continue
				}
				pipelines, err := pp.ListPipelines(ctx, provider.PipelineListOptions{})
				if err != nil {
					log.Printf("error listing pipelines from %s: %v", name, err)
					continue
				}
				fmt.Printf("\n[%s] %d pipelines:\n", name, len(pipelines))
				for _, p := range pipelines {
					fmt.Printf("  %s (%s) - %s\n", p.Name, p.ID, p.Status)
				}
			}
		})

	case "insights":
		runWithEngine(func(eng *engine.Engine) {
			ctx := context.Background()
			cfg := eng.Config()
			insightsEng := insights.NewEngine(cfg.Insights.CostThreshold, cfg.Insights.AnomalyZScore)
			result, err := insightsEng.RunAnalysis(ctx, nil, nil, nil)
			if err != nil {
				log.Fatalf("error running analysis: %v", err)
			}
			fmt.Printf("Analysis complete: %d findings, %d recommendations\n", len(result.Findings), len(result.Recommendations))
			for _, f := range result.Findings {
				fmt.Printf("  [%s] %s: %s\n", f.Severity, f.Type, f.Message)
			}
		})

	case "serve":
		runWithEngine(func(eng *engine.Engine) {
			cfg := eng.Config()
			insightsEng := insights.NewEngine(cfg.Insights.CostThreshold, cfg.Insights.AnomalyZScore)
			srv := api.NewServer(cfg.API.Address, eng, insightsEng, cfg.API.ReadTimeout, cfg.API.WriteTimeout)
			log.Fatal(srv.Start())
		})

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Deep Native Engine CLI v%s

Usage:
  dne <command>

Commands:
  version     Print version information
  providers   List available providers
  resources   List resources from all cloud providers
  pipelines   List pipelines from all pipeline providers
  insights    Run AI insights analysis
  serve       Start the API server
`, version)
}

func runWithEngine(fn func(*engine.Engine)) {
	configPath := "configs/default.yaml"
	if len(os.Args) > 2 {
		configPath = os.Args[2]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

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
			log.Printf("error registering provider %s: %v", pc.Name, err)
			continue
		}
	}

	if err := eng.Start(ctx); err != nil {
		log.Fatalf("error starting engine: %v", err)
	}

	fn(eng)

	if err := eng.Shutdown(ctx); err != nil {
		log.Printf("error during shutdown: %v", err)
	}
}
