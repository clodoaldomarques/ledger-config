# Tarefas: Observabilidade e Health Checks

**Entrada**: documentos em `/specs/004-observability-health/`

**Pré-requisitos**: `plan.md`, `spec.md`, `../platform.md`

**Convenção**: `[x]` = já implementado em `60095c7`; `[ ]` = lacuna identificada.

## Fase 1: Setup

- [x] T001 Declarar variáveis OTEL em `.env` e `scripts/k8s/deployment.yaml`
- [x] T002 [P] Configurar `otel-collector`, `prometheus` e `zipkin` em `docker-compose.yaml`
- [x] T003 [P] Adicionar configurações em `scripts/docker/otel/config.yaml` e `scripts/docker/prometheus/prometheus.yml`

## Fase 2: Fundacional

- [x] T004 Inicializar OTEL e tratar sinais em `cmd/main.go`
- [x] T005 Registrar interceptadores de tracer, logger e métricas em `internal/infra/rest/server/server.go`
- [x] T006 [P] Implementar `Health()` em `internal/infra/db/dynamodb/repository.go`
- [x] T007 [P] Implementar `Health()` em `internal/infra/message/service.go`
- [x] T008 Instalar tracer noop nos testes em `internal/application/ledger/service_test.go`

**Checkpoint**: instrumentação disponível para todas as rotas.

---

## Fase 3: US1 - Liveness (Prioridade: P1) 🎯 MVP

**Objetivo**: responder ao orquestrador mesmo sem dependências.

**Teste independente**: `GET /health/live` → `200` sem corpo.

- [x] T009 [US1] Registrar rota `GET /health/live` em `internal/infra/rest/server/server.go`
- [x] T010 [US1] Implementar handler `Liveness` em `internal/infra/rest/server/server.go`
- [ ] T011 [US1] Adicionar teste HTTP para o endpoint de liveness

---

## Fase 4: US2 - Readiness (Prioridade: P1)

**Objetivo**: só aceitar tráfego com DynamoDB e SNS disponíveis.

**Teste independente**: parar LocalStack e verificar `503` em `/health/ready`.

- [x] T012 [US2] Registrar rota `GET /health/ready` em `internal/infra/rest/server/server.go`
- [x] T013 [US2] Implementar verificação com contexto de 500 ms em `internal/infra/rest/server/server.go`
- [x] T014 [US2] Configurar probes no deployment em `scripts/k8s/deployment.yaml`
- [ ] T015 [US2] Adicionar testes com dependências simuladas falhando

**Checkpoint**: prontidão validada.

---

## Fase 5: US3 - Correlação de requisição (Prioridade: P2)

**Objetivo**: rastrear uma chamada do handler ao publicador.

**Teste independente**: enviar `x-cid` e `x-tenant` e conferir os atributos do span.

- [x] T016 [US3] Implementar `NewSpanFromContext` em `internal/infra/rest/ledger/handler.go`
- [x] T017 [US3] Propagar `cid` para service, repositório e tópico em `internal/application/ledger/service.go`, `internal/infra/db/dynamodb/repository.go` e `internal/infra/message/service.go`
- [x] T018 [US3] Implementar shutdown gracefully em `cmd/main.go`
- [ ] T019 [US3] Validar formato dos headers antes de propagá-los

---

## Fase 6: US4 - Métricas e logs (Prioridade: P2)

**Objetivo**: telemetria centralizada sem código por rota.

**Teste independente**: chamar uma rota com coletor ativo e verificar dados no Prometheus/Zipkin.

- [x] T020 [US4] Aplicar `otelecho.Interceptor()` a todas as rotas em `internal/infra/rest/server/server.go`
- [x] T021 [US4] Registrar erros com `span.SetError` nos casos de uso em `internal/application/ledger/service.go`
- [ ] T022 [US4] Configurar fallback noop quando o coletor não estiver disponível
- [ ] T023 [US4] Adicionar workflow de CI com `go test ./...` e `go vet ./...`

**Checkpoint**: observabilidade completa.

---

## Fase N: Polimento

- [ ] T024 Tornar o timeout de readiness configurável em `internal/infra/rest/server/server.go`
- [ ] T025 [P] Atualizar `docs/specs/004-observability-health/` após as mudanças

## Dependências e Ordem de Execução

- Fase 2 bloqueia todas as stories.
- US1 e US2 são independentes; US3 depende dos interceptadores da Fase 2.
- T022 e T023 são melhorias de ambiente, sem impacto no contrato HTTP.
