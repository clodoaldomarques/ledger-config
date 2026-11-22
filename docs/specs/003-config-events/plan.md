# Plano de Implementação: Eventos de Configuração

**Branch**: `003-config-events` | **Data**: 2026-11-21 | **Spec**: [spec.md](spec.md)

**Entrada**: especificação derivada de `internal/infra/message/` e dos casos de uso em `internal/application/ledger/`.

**Nota**: decisões já implementadas no código em `60095c7`.

## Resumo

Publicar eventos no SNS usando o `core-sdk`, com uma porta `Topic` injetada nos casos de uso. O contrato de mensagem é um espelho da entidade `Config`, convertido explicitamente antes da publicação.

## Contexto Técnico

**Linguagem/Versão**: Go 1.26.1

**Dependências principais**: `core-sdk/pkg/aws/sns`, AWS SDK v2, uuid, gomock

**Armazenamento**: não há persistência de eventos; tópico `config-sns-topic` e fila `script-sqs-queue` no LocalStack

**Testes**: mocks `NewMockTopic` em `service_test.go` e `bench_test.go`

**Plataforma alvo**: serviço HTTP em container

**Tipo de projeto**: adapter de infraestrutura + contrato de mensagem

**Metas de desempenho**: não declaradas

**Restrições**: sem transação com banco, sem retry, sem outbox

**Escopo**: 1 publicador, 1 contrato de mensagem, 1 tópico

## Verificação da Constituição

| Princípio | Situação | Evidência |
|---|---|---|
| I. Separação em camadas | Passa | Casos de uso dependem da interface `Topic`; adapter em `infra/message` |
| II. Domínio valida | Passa | Evento é construído a partir de entidade já validada |
| III. Configuração singleton | Passa | `sns.NewPublisher(ctx, config.New())` |
| IV. Observabilidade no request | Passa | Span `Topic::Emit` com `cid` e `config` |
| V. Consultas no modelo de dados | N/A | Funcionalidade não consulta DynamoDB |
| VI. Testes unitários determinísticos | Passa | `MockTopic` com expectativas nos testes de serviço |
| VII. Infraestrutura como código | Passa | `scripts/terraform/sns.tf`, `sqs.tf` |

## Estrutura do Projeto

### Documentação (esta funcionalidade)

```text
docs/specs/003-config-events/
├── spec.md
├── plan.md
└── tasks.md
```

### Código-fonte (raiz do repositório)

```text
internal/
├── application/ledger/interfaces.go        # Topic
├── application/ledger/mock.go              # MockTopic gerado
├── application/ledger/service.go           # chamadas a Emit
├── infra/message/service.go                # Topic + New + Health
├── infra/message/types.go                  # Config/Script/Account de mensagem
scripts/terraform/{sns,sqs}.tf
```

**Decisão de estrutura**: nenhum pacote novo; contrato de mensagem vive junto do publicador.

## Complexidade Conhecida

| Desvio | Por que existe | Alternativa mais simples considerada |
|---|---|---|
| `Level` serializado como `event_id` | Tag JSON provavelmente copiada de outro evento | Usar `json:"level"` |
| `ToScriptMessage` com `make(len)` + `append` | Crescimento do slice sem reset do comprimento | `make([]Script, 0, len(scs))` |
| Falha de publicação sem rollback | Casos de uso gravam antes de emitir | Outbox ou transação compensatória |
| `Close()` apenas loga | Cliente do core-sdk não exige fechamento explícito | Remover o `defer` dos handlers |
