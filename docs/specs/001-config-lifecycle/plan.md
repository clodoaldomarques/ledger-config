# Plano de Implementação: Ciclo de Vida da Configuração de Ledger

**Branch**: `001-config-lifecycle` | **Data**: 2026-11-21 | **Spec**: [spec.md](spec.md)

**Entrada**: especificação derivada de `internal/`, `config/` e `scripts/terraform/`.

**Nota**: este plano descreve decisões **já presentes** no código; não é um plano a executar do zero.

## Resumo

Expor um CRUD de configurações de ledger com regras no domínio, persistência em DynamoDB com consultas via GSI e publicação de eventos. A implementação atual usa Echo nos handlers, `sync.Once` para configuração e mocks gerados para isolar o caso de uso.

## Contexto Técnico

**Linguagem/Versão**: Go 1.26.1 (`go.mod`)

**Dependências principais**: Echo v4, AWS SDK v2 (DynamoDB), core-sdk (logger/tracer/métricas/SNS), validator, uuid, testify, gomock

**Armazenamento**: DynamoDB `ConfigTable`; `org_id` + `config_id`, GSI `GSI-Index` por `filters`

**Testes**: `go test ./...` (unitários, mocks gerados); integração manual com LocalStack

**Plataforma alvo**: serviço Linux em container (imagem Alpine), porta 5000

**Tipo de projeto**: serviço HTTP

**Metas de desempenho**: não declaradas no repositório; benchmarks disponíveis para caminhos de criação/atualização/consulta

**Restrições**: sem CI, sem autenticação, sem rollback de SNS; `config.New()` captura ambiente apenas uma vez

**Escopo**: 5 rotas HTTP, 1 tabela DynamoDB, 1 tópico SNS

## Verificação da Constituição

*GATE: aplicado antes de Phase 0 e reavaliado após o desenho.*

| Princípio | Situação | Evidência |
|---|---|---|
| I. Separação em camadas | Passa | Handler → Service → Domain; adapters em `internal/infra` |
| II. Domínio valida | Passa | `Config.Validate()`, `validateScript()` e `Script.Validate()` |
| III. Configuração singleton | Passa | `config/config.go` |
| IV. Observabilidade no request | Passa | Interceptadores em `server.routes()` e spans em todas as camadas |
| V. Consultas no modelo de dados | Passa | `buildFilters()`, `buildAllQuery()`, GSI `GSI-Index` |
| VI. Testes unitários determinísticos | Passa | `service_test.go`, `bench_test.go`, `TestMain` com tracer noop |
| VII. Infraestrutura como código | Passa | Terraform, Compose, k8s, Dockerfile |

## Estrutura do Projeto

### Documentação (esta funcionalidade)

```text
docs/specs/001-config-lifecycle/
├── spec.md
├── plan.md
└── tasks.md
```

### Código-fonte (raiz do repositório)

```text
cmd/main.go
config/config.go
internal/
├── application/ledger/{service,interfaces,mock,errors}.go
├── domain/ledger/{config,script,account,errors}.go
├── infra/
│   ├── db/dynamodb/{config,repository,types,errors}.go
│   ├── message/{service,types}.go
│   ├── rest/ledger/{handler,request,response}.go
│   └── rest/server/server.go
scripts/terraform/{dynamo,configs,sns,sqs,provider}.tf
```

**Decisão de estrutura**: monorepo Go único; casos de uso em `internal/application`, regras em `internal/domain`, transporte e adapters em `internal/infra`. Não há diretório `pkg` em uso.

## Complexidade Conhecida

| Desvio | Por que existe | Alternativa mais simples considerada |
|---|---|---|
| Parâmetro `script_id` na rota de disable | Herdado de versão anterior da API; corrigir exige alinhar rota e handler | Renomear o parâmetro da rota para `script_id` — mudaria o contrato público |
| Evento do `PATCH` usa o DTO da request | `UpdateConfig` chama `Emit(..., scr)` em vez de `saved` | Publicar `saved`; exigir teste de contrato do evento |
| Criação duplicada usa erro como mensagem de sucesso | Detecta existência via `FindConfigByLevel` e retorna texto com o id | Usar erro tipado `ErrDuplicatedConfig` |
| `FindAllConfigs` retorna `ErrConfigNotFound` para lista vazia | Reaproveitamento do erro do repositório | Retornar lista vazia com `200` |
