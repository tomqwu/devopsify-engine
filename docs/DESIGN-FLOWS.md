# Deep Native Engine - Design Flows

## 5. Request Flow: GET /api/v1/resources

```
Client                API Server             Registry          Cloud Providers
  │                      │                      │              │
  │  GET /resources      │                      │              │
  │  ?type=ec2           │                      │              │
  │─────────────────────▶│                      │              │
  │                      │                      │              │
  │               ┌──────┴──────┐               │              │
  │               │  Middleware  │               │              │
  │               │  Recovery   │               │              │
  │               │  Logging    │               │              │
  │               │  CORS       │               │              │
  │               └──────┬──────┘               │              │
  │                      │                      │              │
  │                      │  ListByKind("cloud") │              │
  │                      │─────────────────────▶│              │
  │                      │                      │              │
  │                      │  []Provider (cloud)  │              │
  │                      │◀─────────────────────│              │
  │                      │                      │              │
  │                      │       Fan-out: ListResources(ctx, "ec2", opts)
  │                      │─────────────────────────────────────▶│ AWS
  │                      │─────────────────────────────────────▶│ Azure
  │                      │─────────────────────────────────────▶│ GCP
  │                      │                                      │
  │                      │       []resource.Resource (each)     │
  │                      │◀─────────────────────────────────────│
  │                      │                                      │
  │                      │  Aggregate all resources             │
  │  200 OK              │                                      │
  │  [{resources}]       │                                      │
  │◀─────────────────────│                                      │
```

## 6. Request Flow: POST /api/v1/pipelines/{id}/sync

```
Client               API Server            Registry         ArgoCD Provider
  │                      │                    │                  │
  │ POST /pipelines/     │                    │                  │
  │   app-1/sync         │                    │                  │
  │ ?provider=argocd     │                    │                  │
  │─────────────────────▶│                    │                  │
  │                      │                    │                  │
  │                      │ GetPipeline        │                  │
  │                      │ ("argocd")         │                  │
  │                      │───────────────────▶│                  │
  │                      │                    │                  │
  │                      │ PipelineProvider   │                  │
  │                      │◀───────────────────│                  │
  │                      │                    │                  │
  │                      │  TriggerSync(ctx, "app-1", opts)     │
  │                      │─────────────────────────────────────▶│
  │                      │                                      │
  │                      │  SyncResult{status: "syncing"}       │
  │                      │◀─────────────────────────────────────│
  │                      │                                      │
  │  200 OK              │                                      │
  │  {sync_result}       │                                      │
  │◀─────────────────────│                                      │
```

## 7. Event System Flow

```
PUBLISHERS                    EVENT BUS                    SUBSCRIBERS
                         (sync, in-process)

Engine ──────────────┐   ┌─────────────────┐
 engine.started      ├──▶│                 │──▶ Insights Engine
 engine.stopped      │   │   Subscribe()   │──▶ API WebSocket (future)
 provider.registered │   │   Publish()     │──▶ Audit Logger (future)
                     │   │                 │
Cloud Providers ─────┤   │  Event Types:   │
 resource.discovered │   │                 │
 drift.detected      ├──▶│  ┌───────────┐ │
 cost.report.ready   │   │  │ type      │ │
                     │   │  │ time      │ │
Pipeline Providers ──┤   │  │ source    │ │
 pipeline.synced     ├──▶│  │ payload   │ │
 pipeline.failed     │   │  └───────────┘ │
                     │   │                 │
SRE Providers ───────┤   │  Payloads:      │
 incident.triggered  ├──▶│  Provider       │
 incident.resolved   │   │  Resource       │
                     │   │  Drift          │
Insights Engine ─────┤   │  Cost           │
 insight.generated   ├──▶│  Pipeline       │
                     │   │  Incident       │
                     │   │  Insight        │
                     │   └─────────────────┘
```

## 8. Insights Engine Data Flow

```
  DATA SOURCES                  ANALYZERS                    OUTPUT
                                                        
  Cloud Providers               ┌──────────────┐         ┌──────────┐
       │                   ┌───▶│ Cost Analyzer│────────▶│          │
       │                   │    │              │         │          │
       ▼                   │    │ - threshold  │         │          │
  ┌──────────┐             │    │ - idle check │         │          │
  │Resources │─────────────┤    │ - concentrate│         │ Findings │
  │[]Resource│             │    └──────────────┘         │          │
  └──────────┘             │                             │ type     │
                           │    ┌──────────────┐         │ severity │
  ┌──────────┐             │    │Drift Detector│────────▶│ provider │
  │  Costs   │─────────┐  ├───▶│              │         │ message  │
  │[]CostRpt │         │  │    │ - missing=crt│         │          │
  └──────────┘         │  │    │ - deleted=wrn│         └─────┬────┘
                       │  │    │ - tag drift  │               │
  ┌──────────┐         │  │    └──────────────┘               │
  │  Drift   │─────────┤  │                                   ▼
  │[]DriftRs │         │  │    ┌──────────────┐         ┌──────────────┐
  └──────────┘         │  └───▶│Anomaly Detect│────────▶│Recommendation│
                       │       │              │         │  Recommender │
                       └──────▶│ - z-score    │         │              │
                               │ - mean/stdev │         │ - cost review│
                               │ - spike find │         │ - fix drift  │
                               └──────────────┘         │ - investigate│
                                                        │ - critical   │
                                                        └──────────────┘
```

## 9. Startup & Registration Sequence

```
main()
  │
  ├─ config.Load(path)
  │    ├─ os.ReadFile(path)
  │    ├─ Expand ${VAR} via regex
  │    ├─ yaml.Unmarshal onto Default()
  │    └─ Validate()
  │         ├─ Check name required
  │         ├─ Check provider required
  │         ├─ Check kind in [cloud, pipeline, sre]
  │         └─ Check unique names
  │
  ├─ engine.New(WithConfig(cfg))
  │    ├─ Creates Registry (map + RWMutex)
  │    ├─ Creates EventBus (map + RWMutex)
  │    └─ Creates Logger
  │
  ├─ For each ProviderConfig:
  │    ├─ providerFactory[type]() ──▶ new Provider
  │    └─ engine.RegisterProvider(ctx, name, provider, config)
  │         ├─ provider.Init(ctx, config)
  │         │    └─ resolveAuth(config)
  │         │         └─ ResolveEnvValue() for secrets
  │         ├─ registry.Register(name, provider)
  │         └─ eventBus.Publish("provider.registered")
  │
  ├─ engine.Start(ctx)
  │    └─ eventBus.Publish("engine.started")
  │
  ├─ [Server mode] api.NewServer() ──▶ srv.Start()
  │    └─ http.ListenAndServe(:8080)
  │
  ├─ [Wait for SIGINT/SIGTERM]
  │
  └─ Shutdown
       ├─ srv.Shutdown(ctx)
       └─ engine.Shutdown(ctx)
            └─ For each provider: provider.Shutdown(ctx)
```

## 10. Config Resolution

```
Two-layer environment variable resolution:

LAYER 1: YAML-level expansion (during Load)
──────────────────────────────────────────────
  config.yaml:
    server: "${ARGOCD_SERVER}"     ──▶  server: "https://argocd.internal"
                 │
                 └─ regex: \$\{([^}]+)\}
                    os.LookupEnv("ARGOCD_SERVER")

LAYER 2: Provider-level _env suffix (during Init)
──────────────────────────────────────────────────
  config.yaml:
    token_env: ARGOCD_AUTH_TOKEN

  Provider.Init() calls:
    config.ResolveEnvValue(cfg, "token")

    Decision tree:
    ┌──────────────────────────┐
    │ cfg["token"] exists?     │
    │ (direct value)           │
    ├──── YES ──▶ return value │
    ├──── NO                   │
    │                          │
    │ cfg["token_env"] exists? │
    │ (env var name)           │
    ├──── YES ──▶ os.Getenv()  │
    ├──── NO ──▶ return ""     │
    └──────────────────────────┘

  Secrets never stored in YAML ─ only env var names.
```
