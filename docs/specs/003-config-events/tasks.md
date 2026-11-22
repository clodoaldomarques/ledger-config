# Tarefas: Eventos de Configuração

**Entrada**: documentos em `/specs/003-config-events/`

**Pré-requisitos**: `plan.md`, `spec.md`, `../platform.md`

**Convenção**: `[x]` = já implementado em `60095c7`; `[ ]` = lacuna identificada.

## Fase 1: Setup

- [x] T001 Declarar `CONFIG_SNS_TOPIC` em `.env` e no manifest de deploy `scripts/k8s/deployment.yaml`
- [x] T002 [P] Criar tópico e inscrição em `scripts/terraform/{sns,sqs}.tf`

## Fase 2: Fundacional

- [x] T003 Definir porta `Topic` em `internal/application/ledger/interfaces.go`
- [x] T004 [P] Gerar `MockTopic` via `go generate ./internal/application/ledger`
- [x] T005 [P] Implementar contrato de mensagem em `internal/infra/message/types.go`
- [x] T006 Implementar publicador `message.New` em `internal/infra/message/service.go`

**Checkpoint**: publicador disponível para os casos de uso.

---

## Fase 3: US1 - Publicar mudanças (Prioridade: P1) 🎯 MVP

**Objetivo**: toda mutação gera evento.

**Teste independente**: mockar `Topic` e executar criação, atualização, desativação e ativação.

- [x] T007 [US1] Emitir evento em `Service.CreateConfig` em `internal/application/ledger/service.go`
- [x] T008 [US1] Emitir evento em `Service.UpdateConfig`, `DisableConfig`, `EnableConfig` e `ActivateTenant` em `internal/application/ledger/service.go`
- [x] T009 [US1] Propagar erro de `Emit` para o handler em `internal/application/ledger/service.go`
- [x] T010 [US1] Cobrir publicação com `MockTopic` em `internal/application/ledger/service_test.go`
- [ ] T011 [US1] Definir comportamento de retry/outbox para falha de SNS em `internal/application/ledger/service.go`

**Checkpoint**: eventos sendo publicados.

---

## Fase 4: US2 - Contrato do evento (Prioridade: P2)

**Objetivo**: consumidores receberem JSON estável e correto.

**Teste independente**: serializar `ConfigMessage` e comparar com o contrato esperado.

- [x] T012 [US2] Construir envelope com `EventID`, `EventType` e `EventDate` em `internal/infra/message/service.go`
- [x] T013 [US2] Converter Config, Script e Account em `internal/infra/message/types.go`
- [ ] T014 [US2] Corrigir tag JSON do nível de `event_id` para `level` em `internal/infra/message/types.go`
- [ ] T015 [US2] Corrigir `make`/`append` em `ToScriptMessage` em `internal/infra/message/types.go`
- [ ] T016 [US2] Adicionar teste de serialização em `internal/infra/message/types_test.go`

**Checkpoint**: contrato validado.

---

## Fase N: Polimento

- [ ] T017 [P] Documentar o contrato do evento em `docs/specs/003-config-events/spec.md` após correções
- [ ] T018 Rodar `go test ./...`, `go build ./...` e `go vet ./...`

## Dependências e Ordem de Execução

- Fase 2 bloqueia US1 e US2.
- US1 depende dos casos de uso da funcionalidade `001-config-lifecycle`.
- T011 depende de decisão de produto sobre consistência; não é necessária para o MVP.
- T014/T015 mudam o contrato publicado: exigem coordenação com consumidores.
