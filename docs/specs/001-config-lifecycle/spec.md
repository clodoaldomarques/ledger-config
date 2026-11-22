# Especificação da Funcionalidade: Ciclo de Vida da Configuração de Ledger

**Branch da funcionalidade**: `001-config-lifecycle`

**Criado**: 2026-11-21

**Status**: Rascunho (derivado do código em `60095c7`)

**Entrada**: engenharia reversa de `internal/application/ledger`, `internal/domain/ledger`, `internal/infra/db/dynamodb` e `internal/infra/rest/`

## Cenários de Usuário e Testes *(obrigatório)*

### História de Usuário 1 - Criar configuração (Prioridade: P1)

Um serviço consumidor cadastra uma regra de ledger para um tenant ou programa, informando código de processamento e scripts.

**Por que esta prioridade**: é a operação que cria a entidade central e alimenta todas as demais.

**Teste independente**: enviar `POST /v1/ledger/config` com `x-tenant` e payload válido contra DynamoDB LocalStack e SNS mockado; verificar `201`, `config_id` UUID e `version = 1`.

**Cenários de aceitação**:

1. **Dado** um payload com nível `program`, `processing_code`, `program_id` e ao menos um script válido, **Quando** a requisição é enviada com `x-tenant: LEDGER`, **Então** o sistema responde `201` com a configuração criada, `enable = true` e `version = 1`.
2. **Dado** uma configuração já existente com o mesmo nível, tenant, programa e código de processamento, **Quando** uma nova criação é enviada, **Então** o sistema responde `400` com `Config was created with id: <id existente>`.
3. **Dado** um payload com nível `platform`, sem `processing_code` ou sem scripts, **Quando** a criação é enviada, **Então** o sistema responde `400` com o erro de validação correspondente.
4. **Dado** que a gravação no DynamoDB sucedeu mas a publicação no SNS falhou, **Quando** a criação é processada, **Então** o sistema responde `400` e o item permanece gravado (não há rollback).

---

### História de Usuário 2 - Consultar configuração por tipo de evento (Prioridade: P1)

Um consumidor pede a configuração aplicável a um código de processamento e programa.

**Por que esta prioridade**: é a principal operação de leitura em tempo de processamento.

**Teste independente**: semear configurações de program e tenant para o mesmo código e consultar `GET /v1/ledger/config/:event_type_id/:program_id`; verificar que a de program prevalece.

**Cenários de aceitação**:

1. **Dado** configurações habilitadas nos níveis `program` e `tenant` para o mesmo `event_type_id`, **Quando** a consulta é feita, **Então** a configuração de `program` é retornada com `200`.
2. **Dado** apenas uma configuração de tenant habilitada, **Quando** a consulta é feita, **Então** ela é retornada com `200` (fallback).
3. **Dado** que nenhuma configuração habilitada corresponde, **Quando** a consulta é feita, **Então** o sistema responde `404` com `ledger config not found`.
4. **Dado** um `program_id` não numérico, **Quando** a consulta é feita, **Então** o sistema responde `404`.
5. **Dado** um `event_type_id` em minúsculas, **Quando** a consulta é feita, **Então** o sistema o converte para maiúsculas antes de buscar.

---

### História de Usuário 3 - Listar configurações (Prioridade: P2)

Um operador lista as configurações do tenant, opcionalmente filtrando por programa.

**Por que esta prioridade**: suporta auditoria e depuração sem exigir o código de processamento.

**Teste independente**: semear duas configurações para o tenant e uma para outro tenant, chamar `GET /v1/ledger/configs` e verificar que apenas as do tenant retornam.

**Cenários de aceitação**:

1. **Dado** configurações do tenant informado em `x-tenant`, **Quando** `GET /v1/ledger/configs` é chamado, **Então** o sistema responde `200` com a lista.
2. **Dado** o parâmetro `program_id`, **Quando** a listagem é chamada, **Então** apenas configurações cujo `filters` começa com o prefixo do programa são retornadas.
3. **Dado** que não há configurações para o tenant, **Quando** a listagem é chamada, **Então** o sistema responde `400` (erro `ledger config not found` do repositório).

---

### História de Usuário 4 - Atualizar e desativar configuração (Prioridade: P2)

Um operador altera descrição/scripts de uma configuração ou a desativa sem excluí-la.

**Por que esta prioridade**: permite corrigir regras sem perder o histórico representado por `version`.

**Teste independente**: criar uma configuração, fazer `PATCH` e `DELETE .../disable` e verificar `version` incrementado e `enable = false`.

**Cenários de aceitação**:

1. **Dado** uma configuração habilitada, **Quando** `PATCH /v1/ledger/config/:config_id` é enviado com descrição e scripts válidos, **Então** o sistema responde `200` com `version` incrementado e `updated_at` atualizado.
2. **Dado** uma configuração desabilitada, **Quando** `PATCH` é enviado, **Então** o sistema responde `400` com `ledger config not found`.
3. **Dado** uma configuração habilitada, **Quando** `DELETE /v1/ledger/config/:config_id/disable` é chamado, **Então** o sistema responde `200` com `enable = false` e `version` incrementado.
4. **Dado** a rota de desativação atual, **Quando** o handler lê o parâmetro, **Então** ele obtém `script_id` em vez de `config_id` — comportamento conhecido a corrigir.

---

### Casos Limite

- Ausência de `x-tenant`: `org_id` fica vazio e a validação de nível falha; `x-cid` ausente é gerado como UUID.
- `enable` enviado no corpo de criação é ignorado: a criação sempre começa habilitada.
- Uma configuração de nível `tenant` com `program_id` diferente de zero é rejeitada; uma de nível `program` sem `program_id` também.
- Scripts duplicados pela chave `flow#script_id` são rejeitados.
- Falha de publicação SNS após gravar não desfaz a gravação.
- A rota de desativação usa `:config_id`, mas o handler lê `script_id`.

## Requisitos *(obrigatório)*

### Requisitos Funcionais

- **FR-001**: O sistema DEVE aceitar criação somente nos níveis `tenant` e `program`.
- **FR-002**: O sistema DEVE normalizar `processing_code` para maiúsculas.
- **FR-003**: O sistema DEVE gerar `config_id` como UUID e iniciar `version` em 1.
- **FR-004**: O sistema DEVE bloquear criação duplicada para a mesma chave lógica nível + tenant + programa + código.
- **FR-005**: O sistema DEVE validar a entidade no domínio antes de gravar.
- **FR-006**: O sistema DEVE gravar em `ConfigTable` e publicar o evento após a gravação.
- **FR-007**: O sistema DEVE resolver a precedência `program` antes de `tenant` e considerar apenas configurações habilitadas.
- **FR-008**: O sistema DEVE listar por `org_id`, com filtro opcional por prefixo de programa na GSI.
- **FR-009**: O sistema DEVE incrementar `version` e atualizar `updated_at` em alterações.
- **FR-010**: O sistema DEVE rejeitar alterações de configurações desabilitadas.
- **FR-011**: O sistema DEVE representar erros como `{"message": "..."}`, com status `400` para mutações e listagem sem itens, e `404` para busca por tipo de evento sem resultado.
- **FR-012**: O sistema DEVE corrigir a leitura do parâmetro `config_id` no endpoint de desativação.

### Entidades Principais

- **Config**: `config_id`, `level`, `processing_code`, `org_id`, `program_id`, `description`, `scripts`, `enable`, `created_at`, `updated_at`, `version`.
- **Script**: `script_id`, `flow` (`regular` ou `migration`), `description`, `expression`, contas de débito/crédito opcionais; chave de unicidade `flow#script_id`.
- **Account**: `number`, `description`, `cosif`.
- **Item de ConfigTable**: entidade Config + campo derivado `filters`, usado pela GSI `GSI-Index`.

## Critérios de Sucesso *(obrigatório)*

### Resultados Mensuráveis

- **SC-001**: 100% das criações válidas retornam `201` com `config_id` no formato UUID.
- **SC-002**: 100% das duplicatas são bloqueadas antes da gravação.
- **SC-003**: Em 100% das consultas com configuração de program e tenant, a de program é a retornada.
- **SC-004**: Toda alteração bem-sucedida incrementa `version` exatamente em 1.
- **SC-005**: 100% das consultas por tipo de evento ignoram configurações desabilitadas.

## Suposições

- `x-tenant` é confiável: não há autenticação ou autorização no serviço.
- O seed `ledger-101` de nível `platform` existe para origem das configurações de tenant.
- Consumidores aceitam o corpo de erro único `{"message": "..."}`.
- LocalStack com `ConfigTable` e tópico SNS está disponível para testes de integração.
