# Tarefas: Ciclo de Vida da Configuração de Ledger

**Entrada**: documentos em `/specs/001-config-lifecycle/`

**Pré-requisitos**: `plan.md` (obrigatório), `spec.md` (obrigatório), `../platform.md`

**Convenção**: `[x]` = já implementado em `60095c7`; `[ ]` = lacuna identificada.

**Formato**: `[ID] [P?] [US] Descrição com caminho exato do arquivo`

## Fase 1: Setup (Infraestrutura Comum)

**Objetivo**: projeto Go executável com ambiente local.

- [x] T001 Criar módulo `github.com/clodoaldomarques/ledger-config` em `go.mod`
- [x] T002 [P] Declarar dependências Echo, AWS SDK v2, core-sdk, validator e uuid em `go.mod`
- [x] T003 [P] Definir `.env` com endpoint AWS, tópico SNS e variáveis OTEL
- [x] T004 [P] Criar `scripts/terraform/dynamo.tf` com `ConfigTable` e `GSI-Index`
- [x] T005 [P] Criar `scripts/terraform/configs.tf` com o seed `ledger-101`

**Checkpoint**: `go mod download` e `make up` disponíveis.

---

## Fase 2: Fundacional (Pré-requisitos Bloqueantes)

**Objetivo**: camadas compartilhadas por todos os casos de uso.

- [x] T006 [P] Implementar singleton de ambiente em `config/config.go`
- [x] T007 [P] Implementar entidades e validação em `internal/domain/ledger/{config,script,account}.go`
- [x] T008 [P] Implementar `buildFilters`/`buildAllQuery` em `internal/infra/db/dynamodb/types.go`
- [x] T009 Implementar repositório DynamoDB em `internal/infra/db/dynamodb/repository.go`
- [x] T010 Declarar portas `Repository`, `Topic`, `ProgramAPI` em `internal/application/ledger/interfaces.go`
- [x] T011 [P] Gerar mocks com `go generate ./internal/application/ledger` em `internal/application/ledger/mock.go`
- [x] T012 [P] Registrar rotas e interceptadores em `internal/infra/rest/server/server.go`
- [x] T013 Implementar leitura de `x-cid`/`x-tenant` em `internal/infra/rest/ledger/handler.go`
- [x] T014 [P] Implementar DTOs de request/response em `internal/infra/rest/ledger/{request,response}.go`
- [x] T015 Implementar publicador SNS em `internal/infra/message/{service,types}.go`

**Checkpoint**: base pronta; stories podem ser implementadas em paralelo.

---

## Fase 3: US1 - Criar configuração (Prioridade: P1) 🎯 MVP

**Objetivo**: cadastro validado e publicado.

**Teste independente**: `POST /v1/ledger/config` com `x-tenant` e payload válido → `201`.

- [x] T016 [US1] Implementar `Service.CreateConfig` em `internal/application/ledger/service.go`
- [x] T017 [US1] Implementar handler `CreateConfig` em `internal/infra/rest/ledger/handler.go`
- [x] T018 [US1] Cobrir sucesso, duplicata, falha de validação e falha de repositório em `internal/application/ledger/service_test.go`
- [ ] T019 [US1] Adicionar teste de integração com LocalStack para gravação + publicação
- [ ] T020 [US1] Substituir erro de duplicata em texto por erro tipado em `internal/application/ledger/errors.go`

**Checkpoint**: criação funcionando de forma independente.

---

## Fase 4: US2 - Consultar por tipo de evento (Prioridade: P1)

**Objetivo**: resolver configuração aplicável com precedência program → tenant.

**Teste independente**: semear configurações nos dois níveis e consultar a rota de busca.

- [x] T021 [US2] Implementar `FindConfigByLevel` em `internal/application/ledger/service.go`
- [x] T022 [US2] Implementar consulta GSI em `internal/infra/db/dynamodb/repository.go`
- [x] T023 [US2] Implementar handler `FindLedgerConfig` em `internal/infra/rest/ledger/handler.go`
- [x] T024 [US2] Cobrir precedência, fallback e `404` em `internal/application/ledger/service_test.go`

**Checkpoint**: consulta utilizável isoladamente.

---

## Fase 5: US3 - Listar configurações (Prioridade: P2)

**Objetivo**: listar por tenant com filtro opcional de programa.

**Teste independente**: semear configurações de dois tenants e listar com um deles no header.

- [x] T025 [US3] Implementar `FindAllConfigs` em `internal/application/ledger/service.go`
- [x] T026 [US3] Implementar queries com/sem prefixo de programa em `internal/infra/db/dynamodb/repository.go`
- [x] T027 [US3] Implementar handler `FindAllLedgerConfig` em `internal/infra/rest/ledger/handler.go`
- [ ] T028 [US3] Retornar lista vazia com `200` em vez de `ErrConfigNotFound`

**Checkpoint**: listagem validada.

---

## Fase 6: US4 - Atualizar e desativar (Prioridade: P2)

**Objetivo**: manter configurações existentes sem recriá-las.

**Teste independente**: criar, atualizar e desativar uma configuração; conferir `version`.

- [x] T029 [US4] Implementar `UpdateConfig` e `DisableConfig` em `internal/application/ledger/service.go`
- [x] T030 [US4] Implementar handlers `UpdateConfig`/`DisableConfig` em `internal/infra/rest/ledger/handler.go`
- [x] T031 [US4] Implementar `UpdateItem` com expressão dinâmica em `internal/infra/db/dynamodb/repository.go`
- [x] T032 [US4] Cobrir atualização e desativação em `internal/application/ledger/service_test.go`
- [ ] T033 [US4] Corrigir `c.Param("script_id")` para `config_id` em `internal/infra/rest/ledger/handler.go`
- [ ] T034 [US4] Publicar a entidade `saved` no update em `internal/application/ledger/service.go`
- [ ] T035 [US4] Expor `EnableConfig` por rota HTTP em `internal/infra/rest/server/server.go`

**Checkpoint**: ciclo de vida completo; lacunas marcadas permanecem como exercício.

---

## Fase N: Polimento

- [ ] T036 [P] Revisar `README.md` e `docs/specs/` após correções
- [ ] T037 [P] Adicionar testes de contrato HTTP
- [ ] T038 Rodar `go build ./...`, `go vet ./...` e `go test ./...`

## Dependências e Ordem de Execução

- Fase 2 bloqueia todas as stories.
- US1 e US2 são independentes entre si após a Fase 2; US3 e US4 dependem da base de dados criada em US1 para testes manuais.
- Correções T020, T028, T033-T035 devem ser feitas após a story correspondente, com testes de regressão.
