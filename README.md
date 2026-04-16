# Deep Native Engine

A unified Cloud DevOps Control Plane built in Go. The missing integration layer between cloud infrastructure, GitOps pipelines, and SRE incident management.

```
┌─────────────────────────────────────────────┐
│              API Layer (REST)                │
├─────────────────────────────────────────────┤
│              Core Engine                     │
│  Registry │ Event Bus │ Config Management    │
├─────────────────────────────────────────────┤
│           Insights Engine                    │
│  Cost │ Drift │ Anomaly │ Recommendations    │
├─────────────────────────────────────────────┤
│              Providers                       │
│  AWS │ Azure │ GCP │ ArgoCD │ Flux │ GH Act │
│  GitLab CI │ Jenkins │ PagerDuty │ OpsGenie │
└─────────────────────────────────────────────┘
```

**Pattern:** Hexagonal (ports & adapters) with plugin-based provider system

**Module:** `github.com/deepnative/engine`

**Dependencies:** Only `gopkg.in/yaml.v3`

## Quick Start

```bash
# Build
make build

# Run CLI
./bin/dne version
./bin/dne providers

# Run API server
./bin/dne serve configs/examples/multi-cloud.yaml

# Or use the dedicated server binary
DNE_CONFIG=configs/examples/multi-cloud.yaml ./bin/dne-server
```

## Providers

### Cloud (CloudProvider interface)
| Provider | Resources | Status |
|----------|-----------|--------|
| AWS | EC2, S3, IAM, Cost Explorer | Scaffolded |
| Azure | VMs, Storage, Identity | Scaffolded |
| GCP | Compute, Storage, IAM | Scaffolded |

### Pipeline (PipelineProvider interface)
| Provider | Features | Status |
|----------|----------|--------|
| ArgoCD | Applications, sync, history | Scaffolded |
| Flux | Kustomizations, reconciliation | Scaffolded |
| GitHub Actions | Workflows, dispatch | Scaffolded |
| GitLab CI | Pipelines, triggers | Scaffolded |
| Jenkins | Jobs, builds | Scaffolded |

### SRE (SREProvider interface)
| Provider | Features | Status |
|----------|----------|--------|
| PagerDuty | Incidents, on-call | Scaffolded |
| OpsGenie | Incidents, on-call | Scaffolded |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness probe |
| GET | `/readyz` | Readiness probe |
| GET | `/api/v1/providers` | List providers |
| GET | `/api/v1/resources` | List resources |
| GET | `/api/v1/costs` | Get cost data |
| GET | `/api/v1/pipelines` | List pipelines |
| POST | `/api/v1/pipelines/{id}/sync` | Trigger sync |
| GET | `/api/v1/insights` | Run full analysis |
| GET | `/api/v1/insights/cost` | Cost insights |

## Configuration

Config uses YAML with environment variable expansion (`${VAR_NAME}`):

```yaml
providers:
  - name: aws-prod
    provider: aws
    kind: cloud
    config:
      region: us-east-1
      access_key_id_env: AWS_ACCESS_KEY_ID
      secret_access_key_env: AWS_SECRET_ACCESS_KEY
```

See [configs/examples/multi-cloud.yaml](configs/examples/multi-cloud.yaml) for a complete example.

## Architecture

Full design documentation with diagrams:

- [Architecture Overview](docs/ARCHITECTURE.md) - System overview, hexagonal design, provider hierarchy
- [Design Flows](docs/DESIGN-FLOWS.md) - Request flows, event system, insights pipeline
- [Deployment & Packages](docs/DESIGN-DEPLOYMENT.md) - Deployment, dependencies, API reference

## Development

```bash
make test       # Run tests
make lint       # Run linter
make coverage   # Generate coverage report
make fmt        # Format code
make clean      # Clean build artifacts
```

## Docker

```bash
make docker-build
docker run -p 8080:8080 -v $(pwd)/configs:/etc/dne deepnative/engine:latest
```

## License

Apache 2.0 - See [LICENSE](LICENSE)
