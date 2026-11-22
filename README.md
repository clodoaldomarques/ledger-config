# Ledger Config

> Backend service responsible for managing accounting configuration for the Ledger platform.

## Overview

Ledger Config is a Go HTTP service that centralizes the ledger rules (scripts) used by other Ledger services. Configurations are stored in DynamoDB and every mutation publishes an event to an SNS topic. The service is built with Echo and instrumented with OpenTelemetry (traces, metrics and logs).

## Architecture

```text
cmd/                         entry point, graceful shutdown
config/                      env-backed configuration singleton
internal/
  application/ledger/        business service, port interfaces, generated mocks
  domain/ledger/             config/script entities and validation
  domain/program/            program entity
  infra/db/dynamodb/         ConfigTable repository
  infra/message/             SNS publisher
  infra/rest/server/         Echo router, middleware, health checks
  infra/rest/ledger/         handlers and request/response DTOs
  infra/gringotts/           Ollama client (currently unwired)
scripts/                     Dockerfile, Terraform, k8s manifests, mockserver/OTEL/Prometheus config
docs/                        Postman collection and C4 diagrams
```

Request flow: `cmd/main.go` → `server.New()` → `internal/infra/rest/ledger` handler → `internal/application/ledger.Service` → `internal/domain/ledger`. Adapters live under `internal/infra`. Handlers build a new DynamoDB repository and SNS publisher per request; there is no DI container.

Configurations have three levels: `platform` (team-owned), `tenant` (requires `org_id`) and `program` (requires `org_id` and `program_id`). A `processing_code` and at least one script are always required; each script needs a `flow` (`regular` or `migration`) and an `expression`.

## HTTP API

The service listens on `APP_PORT` (default `5000`). `x-tenant` supplies the organization id and `x-cid` the correlation id (generated when missing); neither is sent in the request body.

| Method | Path | Description |
|---|---|---|
| `GET` | `/health/live` | Liveness probe |
| `GET` | `/health/ready` | Checks DynamoDB and SNS (500 ms timeout) |
| `POST` | `/v1/ledger/config` | Create a config → `201` |
| `PATCH` | `/v1/ledger/config/:config_id` | Update description/scripts → `200` |
| `DELETE` | `/v1/ledger/config/:config_id/disable` | Disable a config → `200` (see Known limitations) |
| `GET` | `/v1/ledger/config/:event_type_id/:program_id` | Find by processing code (program level first, then tenant) → `200`/`404` |
| `GET` | `/v1/ledger/configs?program_id=...` | List configs for the tenant → `200` |
| `POST` | `/v1/ledger/org/activate` | Copy `LEDGER` platform configs into a tenant (`{"org_id":"..."}`) → `200` |

Example create body:

```json
{
  "level": "program",
  "processing_code": "101",
  "program_id": 123,
  "description": "Compra a vista",
  "scripts": [
    {
      "script_id": 101,
      "flow": "regular",
      "description": "Compra a vista - Cartão",
      "expression": "Amounts.TRANSACTION_AMOUNT + Fees.IOF",
      "debit_account": { "number": "101", "description": "Conta de débito" },
      "credit_account": { "number": "201", "description": "Conta de crédito" }
    }
  ]
}
```

Ready-to-use requests are in `docs/collections/Accounting-Scripts.postman_collection.json`.

## Configuration

The binary does not load `.env` by itself. `make run` exports it before starting the app; the VS Code launch config uses it as `envFile`.

| Variable | Default | Purpose |
|---|---|---|
| `APP_PORT` | `5000` | HTTP port |
| `AWS_ADDRESS` | `""` | AWS endpoint override (LocalStack in development) |
| `AWS_REGION` | `""` | AWS region |
| `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` | `""` | AWS credentials |
| `CONFIG_SNS_TOPIC` | `""` | Topic ARN for config events |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | — | OTLP collector endpoint |
| `OTEL_SERVICE_NAME` | — | Service name in traces/metrics |

`config.New()` is a `sync.Once` singleton: environment values are captured on the first call.

## Local development

Requirements: Go (see `go.mod`; the Docker image uses Go 1.27.1), Docker with Compose, Make. Terraform and `kubectl` are only needed for infrastructure and deploy; `netcat` (`nc`) is used by `make terraform`.

```bash
go mod download
make run                         # exports .env, then: go run cmd/main.go
```

`go run ./cmd/...` also starts the service, but without loading `.env`.

Local dependencies are declared in `docker-compose.yaml`: LocalStack 4566 (SQS/SNS/DynamoDB), mockserver 1080, OTEL collector 4317/4318, Redis, stackport 8080, Zipkin 9411 and Prometheus 9090.

```bash
make up                          # docker compose up -d, then make terraform
make logs
make down                        # docker compose down -v (removes volumes)
```

Terraform (`scripts/terraform/`) creates the `ConfigTable` DynamoDB table (`org_id` hash key, `config_id` range key, `GSI-Index` on `filters`), the `config-sns-topic` topic, a `script-sqs-queue` subscription and a seed config `ledger-101`.

## Tests and benchmarks

```bash
go test ./...                                        # full suite; unit tests only, no external services needed
go test ./internal/application/ledger -run TestService_CreateScript -count=1
make test                                            # coverage in cover.out, then opens HTML in a browser
go generate ./internal/application/ledger            # regenerates mock.go (requires mockgen from go.uber.org/mock)
make bench-one name=BenchmarkCreateConfig_Success
```

Other Make targets: `bench`, `bench-count`, `bench-baseline`, `bench-compare` (needs `benchstat`), `bench-parallel`, `profile-cpu`, `profile-mem`, `profile-trace`, `profile-clean`. There is no CI, lint or formatter configuration in the repository; use `go build ./...` / `go vet ./...` as additional checks.

## Deploy

```bash
make build version=<tag>         # builds scripts/docker/api/Dockerfile
make push version=<tag>
make apply                       # kubectl apply -f scripts/k8s/
make destroy
```

`make publish version=<tag>` runs build + push. The manifests deploy image `clodoaldomarques/ledger-config:latest`, expose port `5000` (nodePort `31000`), and require a Kubernetes secret named `aws-credentials`. `make restart` and `make reload` are destructive.

## Known limitations

- Terraform is hardcoded to LocalStack at `192.168.67.2:30002`, while Compose publishes LocalStack on `localhost:4566`; `make terraform` waits for that IP before plan/apply. Adjust `scripts/terraform/provider.tf` for other environments.
- The disable route declares `:config_id`, but `DisableConfig` reads `c.Param("script_id")`, so the config id is currently empty.
- Makefile benchmark examples and profile targets still reference old `BenchmarkCreateEvent_*` names; see `internal/application/ledger/bench_test.go` for the current ones.
- `internal/infra/gringotts` is not part of the request path yet.

## Author

**Clodoaldo Marques** — [GitHub](https://github.com/clodoaldomarques) · [LinkedIn](https://www.linkedin.com/in/clodoaldomarques/)
