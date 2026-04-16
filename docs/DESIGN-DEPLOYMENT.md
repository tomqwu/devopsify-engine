# Deep Native Engine - Deployment & Package Design

## 11. Deployment Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Kubernetes / Docker Host                     │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              dne-server container                         │  │
│  │              (distroless/static:nonroot)                  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────┐                  │  │
│  │  │         dne-server binary           │                  │  │
│  │  │         (single static binary)      │                  │  │
│  │  └──────────────┬──────────────────────┘                  │  │
│  │                 │                                         │  │
│  │           ┌─────┴─────┐                                   │  │
│  │           │  :8080    │                                   │  │
│  │           └─────┬─────┘                                   │  │
│  │                 │                                         │  │
│  │    ┌────────────┼────────────┐                            │  │
│  │    │            │            │                            │  │
│  │  /healthz    /readyz    /api/v1/*                         │  │
│  │  (liveness)  (readiness) (business)                       │  │
│  │                                                           │  │
│  │  Config: /etc/dne/config.yaml (mount)                     │  │
│  │  Secrets: Environment variables                           │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                 │
└─────────────────────┬───────────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────────────────────┐
        │             │                             │
        ▼             ▼                             ▼
  ┌───────────┐ ┌───────────┐               ┌───────────┐
  │Cloud APIs │ │  CI/CD    │               │ Incident  │
  │           │ │  Systems  │               │ Platforms │
  │ AWS API   │ │           │               │           │
  │ Azure API │ │ ArgoCD    │               │ PagerDuty │
  │ GCP API   │ │ Flux/K8s  │               │ OpsGenie  │
  │           │ │ GitHub    │               │           │
  │           │ │ GitLab    │               │           │
  │           │ │ Jenkins   │               │           │
  └───────────┘ └───────────┘               └───────────┘

Docker Build (multi-stage):
  Stage 1: golang:1.24-alpine  ──▶  make build
  Stage 2: distroless/static   ──▶  COPY dne-server
```

## 12. Package Dependency Graph

```
cmd/dne ──────────────┐
cmd/dne-server ───────┤
                      │
                      ▼
              ┌── internal/api ──────────────┐
              │       │                      │
              │       ▼                      ▼
              │  internal/engine      internal/insights
              │       │                  │       │
              │       ├──────────┐       │       │
              │       ▼          ▼       │       │
              │  internal/    internal/   │       │
              │  eventbus     config      │       │
              │       │          │        │       │
              │       ▼          │        │       │
              │  pkg/event       │        │       │
              │       │          │        │       │
              └───────┼──────────┼────────┼───────┘
                      │          │        │
                      ▼          ▼        ▼
                   pkg/provider ◀── pkg/resource
                      ▲
                      │
              providers/* (adapters)
                      │
                      └──▶ internal/config (ResolveEnvValue only)

DEPENDENCY RULES:
  pkg/        ──▶ Zero imports from internal/ (true ports)
  providers/* ──▶ Only imports pkg/ and internal/config
  internal/   ──▶ Only imports pkg/
  cmd/        ──▶ Imports everything (composition root)
```

## 13. API Reference

```
METHOD  PATH                        HANDLER                 DESCRIPTION
──────  ────                        ───────                 ───────────
GET     /healthz                    HealthHandler.Healthz   Liveness probe
GET     /readyz                     HealthHandler.Readyz    Readiness (checks providers)
GET     /api/v1/providers           CloudHandler.List       List all providers
GET     /api/v1/resources           CloudHandler.Resources  List resources (fan-out)
GET     /api/v1/costs               CloudHandler.Costs      Get cost data
GET     /api/v1/pipelines           PipelineHandler.List    List pipelines
POST    /api/v1/pipelines/{id}/sync PipelineHandler.Sync    Trigger pipeline sync
GET     /api/v1/insights            InsightsHandler.Get     Full analysis
GET     /api/v1/insights/cost       InsightsHandler.Cost    Cost-only insights

Query Parameters:
  /resources  ?type=ec2&region=us-east-1
  /costs      ?granularity=daily&group_by=service
  /pipelines  ?namespace=production
  /pipelines/{id}/sync  ?provider=argocd-main
```

## 14. Key Design Decisions

```
DECISION                           RATIONALE
────────                           ─────────
Go 1.24, single dependency         Minimal footprint, fast builds,
(gopkg.in/yaml.v3)                 no dependency management overhead

Hexagonal architecture             Provider implementations are
                                   fully decoupled from core

Synchronous event bus              Simple, predictable behavior.
                                   Async can be added later.

Universal Resource model           All clouds normalize to one
                                   struct for unified querying

Provider factory map               No reflection, explicit
                                   registration, easy to extend

Config env resolution              Secrets never in YAML files,
(${VAR} + _env suffix)             two-layer resolution for flexibility

Compile-time interface checks      Catch missing methods at build
(var _ Interface = (*T)(nil))      time, not runtime

Fan-out with silent skip           Resilience - one provider failure
on errors                          doesn't break the whole response

Distroless container               Minimal attack surface,
                                   no shell, no package manager
```
