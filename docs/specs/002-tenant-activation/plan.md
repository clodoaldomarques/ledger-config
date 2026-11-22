# Plano de Implementação: Ativação de Tenant

**Branch**: `002-tenant-activation` | **Data**: 2026-11-21 | **Spec**: [spec.md](spec.md)

**Entrada**: especificação derivada de `internal/application/ledger/service.go`, `internal/infra/rest/ledger/handler.go` e `scripts/terraform/`.

**Nota**: decisões já implementadas no código em `60095c7`.

## Resumo

Reutilizar as configurações `platform` da organização `LEDGER` como template para um novo tenant. O caso de uso combina leitura em lote, verificação de existência por nível, gravação individual e publicação de eventos.

## Contexto Técnico

**Linguagem/Versão**: Go 1.26.1

**Dependências principais**: AWS SDK v2 (DynamoDB), core-sdk (tracer/logger/SNS), uuid, validator

**Armazenamento**: `ConfigTable`, GSI `GSI-Index`; `filters` do tipo `TENANT#<org>#PROGRAMID#...#EVENTTYPEID#...`

**Testes**: unitários com gomock em `service_test.go`; não há teste específico para ativação

**Plataforma alvo**: serviço HTTP em container

**Tipo de projeto**: caso de uso em serviço HTTP existente

**Metas de desempenho**: não declaradas

**Restrições**: sem transação no DynamoDB; origem fixa `LEDGER`; sem autorização

**Escopo**: 1 endpoint, 1 caso de uso, nenhuma nova tabela

## Verificação da Constituição

| Princípio | Situação | Evidência |
|---|---|---|
| I. Separação em camadas | Passa | DTO em `rest/ledger`, regra em `application/ledger`, persistência em `infra/db` |
| II. Domínio valida | Passa parcial | `org_id` validado no DTO; cópia não chama `Validate()` na entidade |
| III. Configuração singleton | Passa | `config.New()` usado por repositório e publicador |
| IV. Observabilidade no request | Passa | Span `Service::ActivateTenant` e cid do header |
| V. Consultas no modelo de dados | Passa | `FindAllConfigs` e `FindConfigByLevel` usam a GSI |
| VI. Testes unitários determinísticos | Atenção | Não há testes para `ActivateTenant` |
| VII. Infraestrutura como código | Passa | Seed `ledger-101` em Terraform |

## Estrutura do Projeto

### Documentação (esta funcionalidade)

```text
docs/specs/002-tenant-activation/
├── spec.md
├── plan.md
└── tasks.md
```

### Código-fonte (raiz do repositório)

```text
internal/
├── application/ledger/service.go          # ActivateTenant
├── application/ledger/errors.go           # ErrOrgActivated
├── infra/rest/ledger/handler.go           # ActivateTenant
├── infra/rest/ledger/request.go           # PostOrgActivate
├── infra/db/dynamodb/repository.go        # FindAllConfigs / SaveConfig
├── infra/db/dynamodb/types.go             # buildFilters
└── infra/message/service.go               # Emit
scripts/terraform/configs.tf               # configuração modelo
```

**Decisão de estrutura**: reaproveitar camadas existentes; nenhum pacote novo.

## Complexidade Conhecida

| Desvio | Por que existe | Alternativa mais simples considerada |
|---|---|---|
| `ErrOrgActivated` é retornado quando nada é criado | Reaproveitamento de erro tipado existente | Retornar `200` com lista vazia |
| Timestamps locais na cópia | Uso de `time.Now()` em vez de `time.Now().UTC()` | Padronizar para UTC como nos outros casos de uso |
| Sem transação | DynamoDB não oferece transação no código atual | `TransactWriteItems` para todas as cópias |
| Configurações modelo desabilitadas copiadas | `FindAllConfigs` não filtra `enable` | Filtrar na aplicação antes de copiar |
