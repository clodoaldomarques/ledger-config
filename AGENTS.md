# AGENTS.md

## Run and verify

- `go test ./...` is the full suite: unit tests only (gomock), no Docker/AWS/LocalStack needed; `TestMain` installs a noop tracer. Verified passing offline.
- `make test` writes `cover.out` then runs `go tool cover -html=cover.out` (opens a browser / blocks). For non-interactive checks use `go test ./...` or `go test ./... -coverprofile=cover.out`.
- Focused test: `go test ./internal/application/ledger -run TestService_CreateScript -count=1`.
- Mocks: `go generate ./internal/application/ledger` (needs `mockgen` from `go.uber.org/mock`). `internal/application/ledger/mock.go` is generated — edit `interfaces.go` and regenerate, do not hand-edit.
- Benchmarks live in `internal/application/ledger/bench_test.go`; real names are `BenchmarkCreateConfig_Success`, `BenchmarkUpdateConfig_Success`, `BenchmarkFindConfigByLevel_ProgramHit`, `BenchmarkFindAllConfigs`, `BenchmarkCreateConfig_Parallel`. Makefile examples/profile targets still say `BenchmarkCreateEvent_*` and are stale. Use `make bench-one name=BenchmarkCreateConfig_Success` (Makefile correctly adds `-run=^$`).
- No CI workflow, lint, or formatter config exists. `go build ./...` / `go vet ./...` are the de facto checks.
- `make run` exports `.env` then runs `go run cmd/main.go`. Plain `go run ./cmd/...` (as in README) does **not** load `.env`, so AWS/OTEL settings fall back to empty defaults.

## Wiring and gotchas

- Flow: `cmd/main.go` → `internal/infra/rest/server.New()` (Echo routes + tracer/logger/metrics middleware) → handlers in `internal/infra/rest/ledger` → `internal/application/ledger.Service` → `internal/domain/ledger`; adapters live in `internal/infra/{db/dynamodb,message,gringotts}`. Handlers construct a fresh DynamoDB repo and SNS publisher per request; there is no DI container.
- `config.New()` is a `sync.Once` singleton: env values are captured on first call, while `WithX` options mutate that same instance on later calls. `APP_PORT` defaults to 5000.
- `x-cid` and `x-tenant` headers supply the correlation id and org id; `x-cid` is auto-generated when absent. There is no auth middleware.
- Route `DELETE /v1/ledger/config/:config_id/disable` declares `:config_id`, but `DisableConfig` reads `c.Param("script_id")` — that value is empty today. Keep parameter names aligned when touching this endpoint.
- `internal/infra/gringotts` (Ollama) and the `ProgramAPI` interface are unwired; `Config.GriphookIAUrl` has no env loader. They are not part of the live request path.

## Local infra and deploy

- `make up` = `docker compose up -d` + `make terraform`; `make down` = `docker compose down -v` (deletes LocalStack data). `restart`/`reload` are destructive.
- Compose provides LocalStack 4566 (SQS/SNS/DynamoDB), mockserver 1080, OTEL collector 4317/4318, Redis, stackport 8080, Zipkin and Prometheus. `.env` points the app at `localhost:4566`, SNS topic `arn:aws:sns:us-east-1:000000000000:config-sns-topic`, OTEL at `localhost:4317`.
- Terraform is hardcoded to LocalStack at `192.168.67.2:30002` (not compose's `localhost:4566`); `make terraform` blocks on `nc -z 192.168.67.2 30002` before plan/apply. Expect it to hang on a fresh machine.
- `make apply`/`destroy` use `kubectl` on `scripts/k8s/`; the deployment requires an `aws-credentials` secret and image `clodoaldomarques/ledger-config:latest`. `make build`/`push` need `version=` (no default). The Docker image builds `./cmd` with Go 1.27.1 while `go.mod` says 1.26.1.

## Repo quirks

- `.gitignore` ignores `*.md`, `.env`, `.vscode/`, `docker-compose*.yml`, Terraform state/lock files and `cover.out`. New Markdown files (including this one) need `git add -f` to be tracked; already-tracked files are unaffected.
