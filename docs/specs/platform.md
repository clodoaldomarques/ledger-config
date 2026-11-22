# Contexto Técnico — Ledger Config

Documento de apoio às specs, derivado de `go.mod`, `config/config.go`, `docker-compose.yaml`, `scripts/` e `Makefile`.

## Stack

| Item | Versão/configuração |
|---|---|
| Linguagem | Go 1.26.1 (`go.mod`); imagem Go 1.27.1 (`scripts/docker/api/Dockerfile`) |
| HTTP | Echo v4.16.0 |
| Persistência | AWS SDK v2 + DynamoDB (`ConfigTable`) |
| Eventos | SNS via `github.com/clodoaldomarques/core-sdk` v0.8.8 |
| Observabilidade | OpenTelemetry 1.47 (traces/métricas) e logger Zap do core-sdk |
| Validação | `go-playground/validator` + validação manual no domínio |
| Testes | `stretchr/testify`, `go.uber.org/mock` (mockgen), `bench_test.go` |
| Outros | `google/uuid`, `ollama/api` (adapter não ligado) |

## Configuração

| Variável | Origem | Uso |
|---|---|---|
| `APP_PORT` | `config.New()` | Porta HTTP, padrão `5000` |
| `AWS_ADDRESS`, `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` | `.env` / k8s | Cliente AWS e LocalStack |
| `CONFIG_SNS_TOPIC` | `.env` / k8s | ARN do tópico de configurações |
| `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME` | `.env` / k8s | Exportação OTLP |

`config.New()` é singleton (`sync.Once`): valores do ambiente são lidos na primeira chamada. O binário não carrega `.env` automaticamente.

## Modelo de dados

- Tabela `ConfigTable`: `org_id` (hash), `config_id` (range), GSI `GSI-Index` com `filters` como range key e projeção `ALL`.
- `filters`: `TENANT#<org>#PROGRAMID#<0000|id formatado>%04d#EVENTTYPEID#<código>`. Para nível `platform`, o tenant é `LEDGER`.
- Item semeado por Terraform: `ledger-101`, org `LEDGER`, nível `platform`, código `101`.
- Evento SNS: tópico `config-sns-topic`, inscrito na fila `script-sqs-queue`.

## Ambiente local

`docker-compose.yaml` publica:

| Serviço | Porta | Papel |
|---|---|---|
| localstack | 4566 | SQS, SNS, DynamoDB |
| mockserver | 1080 | Expectativas HTTP (`scripts/docker/mockserver/expectations.json`) |
| otel-collector | 4317, 4318, 8888 | OTLP gRPC/HTTP e métricas |
| redis | 6379 | Suporte local |
| stackport | 8080 | Visualização de recursos AWS no LocalStack |
| zipkin | 9411 | UI de traces |
| prometheus | 9090 | Coleta de métricas |

`make up` executa `docker compose up -d` e `make terraform`. O provider Terraform aponta para `192.168.67.2:30002` (valor fixo em `scripts/terraform/provider.tf`), não para o `localhost:4566` do Compose; `make terraform` espera esse endereço com `nc` antes de planejar/aplicar.

## Deploy

`make build version=<tag>` constrói a imagem com `scripts/docker/api/Dockerfile`; `make push` envia as tags `<tag>` e `latest`; `make apply` aplica `scripts/k8s/`. Os manifests usam a imagem `clodoaldomarques/ledger-config:latest`, Service `LoadBalancer` na porta 5000 (nodePort 31000) e exigem o secret `aws-credentials`. `make destroy` remove os recursos.

## Verificação

```bash
go test ./...                                     # suíte unitária, sem serviços externos
go test ./internal/application/ledger -run TestService_CreateScript -count=1
make test                                         # grava cover.out e abre HTML no navegador
go generate ./internal/application/ledger         # requer mockgen de go.uber.org/mock
go build ./... && go vet ./...
make bench-one name=BenchmarkCreateConfig_Success
```

Não há workflow de CI, lint ou formatador configurado no repositório.

## Desvios conhecidos

- `DELETE /v1/ledger/config/:config_id/disable`: handler lê `c.Param("script_id")`; o id chega vazio.
- `UpdateConfig` publica no SNS o DTO da request (sem `config_id` e com `version` zerado) em vez da entidade atualizada.
- `message.Config.Level` usa a tag JSON `event_id`; `ToScriptMessage` cria um slice com `len` e depois faz `append`, produzindo entradas zeradas no início.
- `Service.EnableConfig` existe, mas não há rota HTTP que a exponha.
- `internal/infra/gringotts` e a interface `ProgramAPI` não estão ligados ao fluxo de request.
- Makefile e alvos de perfilagem citam `BenchmarkCreateEvent_*`; os benchmarks reais usam `BenchmarkCreateConfig_*`.
