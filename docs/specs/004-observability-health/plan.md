# Plano de Implementação: Observabilidade e Health Checks

**Branch**: `004-observability-health` | **Data**: 2026-11-21 | **Spec**: [spec.md](spec.md)

**Entrada**: especificação derivada de `cmd/main.go`, `internal/infra/rest/server/server.go` e `scripts/docker/`.

**Nota**: decisões já implementadas no código em `60095c7`.

## Resumo

Instrumentar todas as rotas com os interceptadores do `core-sdk`, expor probes de liveness e readiness e desligar telemetria no encerramento do processo. A configuração de exportação vem de variáveis de ambiente.

## Contexto Técnico

**Linguagem/Versão**: Go 1.26.1

**Dependências principais**: `core-sdk/pkg/otel`, `core-sdk/pkg/otel/echo`, `core-sdk/pkg/zap/logger`, Echo v4

**Armazenamento**: N/A (readiness usa `DescribeTable` na tabela existente)

**Testes**: `TestMain` em `service_test.go` instala tracer noop; não há testes de health

**Plataforma alvo**: container Linux, probes no deployment Kubernetes

**Tipo de projeto**: Cross-cutting em serviço HTTP

**Metas de desempenho**: readiness com timeout fixo de 500 ms

**Restrições**: sem CI; comportamento offline do OTEL depende do core-sdk

**Escopo**: 2 endpoints, 3 interceptadores, configuração de shutdown

## Verificação da Constituição

| Princípio | Situação | Evidência |
|---|---|---|
| I. Separação em camadas | Passa | Server monta middleware; handlers extraem headers |
| II. Domínio valida | N/A | Não há regras de domínio nesta funcionalidade |
| III. Configuração singleton | Passa | Variáveis OTEL lidas pelo SDK; `APP_PORT` via `config.New()` |
| IV. Observabilidade no request | Passa | Interceptadores e `NewSpanFromContext` |
| V. Consultas no modelo de dados | Passa parcial | Readiness usa `DescribeTable`, não a GSI |
| VI. Testes unitários determinísticos | Atenção | Tracer noop já configurado, mas sem testes dos probes |
| VII. Infraestrutura como código | Passa | Compose otel-collector/zipkin/prometheus, probes no k8s |

## Estrutura do Projeto

### Documentação (esta funcionalidade)

```text
docs/specs/004-observability-health/
├── spec.md
├── plan.md
└── tasks.md
```

### Código-fonte (raiz do repositório)

```text
cmd/main.go                                  # otel.Start, sinais, shutdown
config/config.go                             # APP_PORT
internal/infra/rest/server/server.go         # rotas, interceptadores, probes
internal/infra/rest/ledger/handler.go        # NewSpanFromContext
internal/infra/db/dynamodb/repository.go     # Health()
internal/infra/message/service.go            # Health()
scripts/docker/otel/config.yaml
scripts/docker/prometheus/prometheus.yml
scripts/k8s/deployment.yaml                  # livenessProbe/readinessProbe
.env                                         # OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_SERVICE_NAME
```

**Decisão de estrutura**: nenhum pacote de observabilidade próprio; tudo delegado ao `core-sdk`.

## Complexidade Conhecida

| Desvio | Por que existe | Alternativa mais simples considerada |
|---|---|---|
| Timeout de readiness fixo em 500 ms | Valor embutido em `server.go` | Tornar configurável por variável de ambiente |
| Sem testes dos probes | Suíte cobre apenas o caso de uso de ledger | Adicionar testes com clientes HTTP fake |
| Sem CI | Repositório de estudo/portfólio | Workflow com `go test` e `go vet` |
| OTEL sem fallback explícito | Inicialização delegada ao core-sdk | Configurar exporter noop quando não houver coletor |
