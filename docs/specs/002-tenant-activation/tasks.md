# Tarefas: Ativação de Tenant

**Entrada**: documentos em `/specs/002-tenant-activation/`

**Pré-requisitos**: `plan.md`, `spec.md`, `../platform.md`

**Convenção**: `[x]` = já implementado em `60095c7`; `[ ]` = lacuna identificada.

## Fase 1: Setup

- [x] T001 Garantir seed `ledger-101` em `scripts/terraform/configs.tf`
- [x] T002 [P] Registrar `POST /v1/ledger/org/activate` em `internal/infra/rest/server/server.go`

## Fase 2: Fundacional

- [x] T003 Reutilizar `FindAllConfigs`, `FindConfigByLevel` e `SaveConfig` em `internal/infra/db/dynamodb/repository.go`
- [x] T004 [P] Definir `ErrOrgActivated` em `internal/application/ledger/errors.go`
- [x] T005 [P] Definir DTO `PostOrgActivate` e validação em `internal/infra/rest/ledger/request.go`

**Checkpoint**: dependências disponíveis para o caso de uso.

---

## Fase 3: US1 - Ativar tenant (Prioridade: P1) 🎯 MVP

**Objetivo**: copiar as configurações modelo para o tenant informado.

**Teste independente**: `POST /v1/ledger/org/activate` com `{"org_id":"TENANT_A"}` → `200` com cópias.

- [x] T006 [US1] Implementar `Service.ActivateTenant` em `internal/application/ledger/service.go`
- [x] T007 [US1] Implementar handler `ActivateTenant` em `internal/infra/rest/ledger/handler.go`
- [x] T008 [US1] Publicar evento para cada cópia em `internal/application/ledger/service.go`
- [ ] T009 [US1] Adicionar teste unitário de sucesso com mock de repositório e tópico em `internal/application/ledger/service_test.go`
- [ ] T010 [US1] Adicionar teste de integração com LocalStack cobrindo o seed

**Checkpoint**: ativação utilizável.

---

## Fase 4: US2 - Reativação sem duplicatas (Prioridade: P2)

**Objetivo**: repetir a operação sem criar novas configurações.

**Teste independente**: ativar duas vezes e comparar a quantidade de itens no DynamoDB.

- [x] T011 [US2] Verificar existência com `FindConfigByLevel` antes de salvar em `internal/application/ledger/service.go`
- [ ] T012 [US2] Retornar `200` com lista vazia quando nada é criado em `internal/application/ledger/service.go`
- [ ] T013 [US2] Adicionar teste de reativação em `internal/application/ledger/service_test.go`

**Checkpoint**: operação idempotente do ponto de vista de dados.

---

## Fase N: Polimento

- [ ] T014 Padronizar `time.Now()` para UTC na cópia em `internal/application/ledger/service.go`
- [ ] T015 Ignorar configurações modelo desabilitadas em `internal/application/ledger/service.go`
- [ ] T016 [P] Revisar `docs/specs/002-tenant-activation/` após as correções
- [ ] T017 Rodar `go test ./...`, `go build ./...` e `go vet ./...`

## Dependências e Ordem de Execução

- Fase 2 bloqueia as stories.
- US2 depende da verificação de existência implementada em US1.
- T014 e T015 são melhorias independentes, mas devem ter testes de regressão.
