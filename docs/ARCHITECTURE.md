# Deep Native Engine - Architecture Design

## 1. System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        ENTRY POINTS                                 │
│                                                                     │
│    ┌──────────────┐                    ┌──────────────────┐         │
│    │   dne CLI    │                    │   dne-server     │         │
│    │  (cmd/dne)   │                    │ (cmd/dne-server) │         │
│    └──────┬───────┘                    └────────┬─────────┘         │
│           │                                     │                   │
│           └──────────────┬──────────────────────┘                   │
│                          ▼                                          │
├─────────────────────────────────────────────────────────────────────┤
│                       REST API LAYER                                │
│                                                                     │
│   Recovery ─▶ Logging ─▶ CORS ─▶ ServeMux                          │
│                                                                     │
│   /healthz  /readyz  /api/v1/providers  /api/v1/resources           │
│   /api/v1/costs  /api/v1/pipelines  /api/v1/insights               │
├─────────────────────────────────────────────────────────────────────┤
│                       CORE ENGINE                                   │
│                                                                     │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐         │
│   │   Registry   │  │  Event Bus   │  │ Config Manager   │         │
│   │              │  │              │  │                  │         │
│   │ Register()   │  │ Subscribe()  │  │ Load()           │         │
│   │ Get()        │  │ Publish()    │  │ Validate()       │         │
│   │ GetCloud()   │  │              │  │ ResolveEnvValue()│         │
│   │ GetPipeline()│  │              │  │                  │         │
│   │ GetSRE()     │  │              │  │                  │         │
│   │ ListByKind() │  │              │  │                  │         │
│   └──────────────┘  └──────────────┘  └──────────────────┘         │
├─────────────────────────────────────────────────────────────────────┤
│                     INSIGHTS ENGINE                                 │
│                                                                     │
│   ┌────────────┐ ┌────────────┐ ┌──────────┐ ┌────────────────┐   │
│   │    Cost    │ │   Drift    │ │ Anomaly  │ │ Recommendation │   │
│   │  Analyzer  │ │ Detector   │ │ Detector │ │   Recommender  │   │
│   │            │ │            │ │          │ │                │   │
│   │ threshold  │ │ severity   │ │ z-score  │ │ prioritize     │   │
│   │ idle check │ │ classify   │ │ analysis │ │ actionable     │   │
│   └────────────┘ └────────────┘ └──────────┘ └────────────────┘   │
├─────────────────────────────────────────────────────────────────────┤
│                        PROVIDERS                                    │
│                                                                     │
│   CLOUD              PIPELINE               SRE                    │
│   ┌─────┐ ┌───────┐ ┌────────┐ ┌──────┐   ┌───────────┐          │
│   │ AWS │ │ Azure │ │ ArgoCD │ │ Flux │   │ PagerDuty │          │
│   └─────┘ └───────┘ └────────┘ └──────┘   └───────────┘          │
│   ┌─────┐            ┌────────┐ ┌────────┐ ┌───────────┐          │
│   │ GCP │            │ GH Act │ │GitLab │ │ OpsGenie  │          │
│   └─────┘            └────────┘ └────────┘ └───────────┘          │
│                      ┌─────────┐                                   │
│                      │ Jenkins │                                   │
│                      └─────────┘                                   │
└─────────────────────────────────────────────────────────────────────┘
```

## 2. Hexagonal Architecture (Ports & Adapters)

```
                    DRIVING ADAPTERS                        DRIVEN ADAPTERS
                    (inbound)                               (outbound)

                                  ┌─────────────────┐
                ┌────────────────▶│                 │──────────────┐
  ┌──────────┐  │                 │    PORTS        │              │  ┌─────────────┐
  │ REST API │──┘    ┌───────────▶│  (interfaces)   │─────────┐   ├─▶│ AWS SDK     │
  │ Handlers │       │            │                 │         │   │  │ Azure SDK   │
  └──────────┘       │            │ CloudProvider   │         │   │  │ GCP SDK     │
                     │            │ PipelineProvider│         │   │  └─────────────┘
  ┌──────────┐       │            │ SREProvider     │         │   │
  │   CLI    │───────┘            │                 │         │   │  ┌─────────────┐
  │ Commands │                    │    DOMAIN       │         ├───├─▶│ ArgoCD API  │
  └──────────┘                    │    CORE         │         │   │  │ Flux K8s    │
                                  │                 │         │   │  │ GH Actions  │
                                  │ Engine          │         │   │  │ GitLab API  │
                                  │ Registry        │         │   │  │ Jenkins API │
                                  │ EventBus        │         │   │  └─────────────┘
                                  │ Resource Model  │         │   │
                                  │ Diff Algorithm  │         │   │  ┌─────────────┐
                                  │                 │         └───┴─▶│ PagerDuty   │
                                  └─────────────────┘                │ OpsGenie    │
                                                                     └─────────────┘
    pkg/ = PORTS (zero internal/ imports)
    internal/ = DOMAIN CORE
    providers/ = DRIVEN ADAPTERS
    internal/api/ = DRIVING ADAPTER
    cmd/ = COMPOSITION ROOT
```

## 3. Provider Interface Hierarchy

```
                    ┌──────────────────────────┐
                    │     Provider (base)       │
                    │──────────────────────────│
                    │ + Metadata() Metadata     │
                    │ + Init(ctx, config) error │
                    │ + Healthy(ctx) error      │
                    │ + Shutdown(ctx) error     │
                    └─────────┬────────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
┌──────────────────┐ ┌────────────────┐ ┌────────────────┐
│  CloudProvider   │ │PipelineProvider│ │  SREProvider   │
│──────────────────│ │────────────────│ │────────────────│
│ ListResources()  │ │ ListPipelines()│ │ ListIncidents()│
│ GetResource()    │ │ GetStatus()    │ │ GetIncident()  │
│ GetCostData()    │ │ TriggerSync()  │ │ Acknowledge()  │
│ DetectDrift()    │ │ GetHistory()   │ │ Resolve()      │
└───────┬──────────┘ └───────┬────────┘ │ ListOnCall()   │
        │                    │          └───────┬────────┘
        │                    │                  │
   ┌────┼────┐        ┌──┬──┼──┬──┬──┐    ┌────┼────┐
   ▼    ▼    ▼        ▼  ▼  ▼  ▼  ▼  ▼    ▼         ▼
  AWS Azure GCP    Argo Flux GH GL Jen   PagerDuty OpsGenie

All providers use compile-time check:
  var _ provider.CloudProvider = (*Provider)(nil)
```

## 4. Universal Resource Model

```
┌──────────────────────────────────────┐
│          resource.Resource           │
│──────────────────────────────────────│
│ ID           string                  │
│ Name         string                  │
│ Type         string                  │  Types:
│ Provider     string                  │    compute:instance
│ Region       string                  │    storage:bucket
│ State        State ──────────────┐   │    identity:user
│ Tags         map[string]string   │   │    compute:vm
│ Properties   map[string]any      │   │    storage:account
│ CostPerMonth float64             │   │    iam:serviceaccount
│ CreatedAt    time.Time           │   │
│ UpdatedAt    time.Time           │   │
└──────────────────────────────────┘   │
                                       ▼
                              ┌─────────────────┐
                              │     State        │
                              │─────────────────│
                              │ running          │
                              │ stopped          │
                              │ terminated       │
                              │ pending          │
                              │ unknown          │
                              │ degraded         │
                              │ healthy          │
                              │ unhealthy        │
                              └─────────────────┘

┌──────────────────────────────────────┐
│        resource.DriftResult          │
│──────────────────────────────────────│
│ ResourceID   string                  │
│ ResourceType string                  │  DriftType:
│ DriftType    DriftType ──────────┐   │    modified
│ Field        string              │   │    added
│ Expected     string              │   │    deleted
│ Actual       string              │   │    missing
│ Message      string              │   │
└──────────────────────────────────┘   │
                                       ▼
         resource.Diff(desired, actual []Resource) ──▶ []DriftResult
```
