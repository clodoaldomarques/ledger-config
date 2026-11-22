# Especificação da Funcionalidade: Eventos de Configuração

**Branch da funcionalidade**: `003-config-events`

**Criado**: 2026-11-21

**Status**: Rascunho (derivado do código em `60095c7`)

**Entrada**: engenharia reversa de `internal/infra/message/`, `internal/application/ledger/interfaces.go` e `scripts/terraform/sns.tf`

## Cenários de Usuário e Testes *(obrigatório)*

### História de Usuário 1 - Publicar mudanças de configuração (Prioridade: P1)

Serviços consumidores precisam ser avisados quando uma configuração é criada, alterada, desativada ou copiada para um tenant.

**Por que esta prioridade**: a API é a fonte de configuração; sem eventos, os consumidores ficam com estado defasado.

**Teste independente**: com publicador mockado, executar cada mutação e verificar que `Emit` é chamado uma vez com o `cid` da requisição.

**Cenários de aceitação**:

1. **Dado** uma criação bem-sucedida, **Quando** o caso de uso termina, **Então** um evento é publicado no tópico `CONFIG_SNS_TOPIC`.
2. **Dado** uma atualização ou desativação bem-sucedida, **Quando** o caso de uso termina, **Então** um evento é publicado.
3. **Dado** uma ativação de tenant que cria N configurações, **Quando** a operação termina, **Então** N eventos são publicados.
4. **Dado** que a publicação falha, **Quando** a mutação é processada, **Então** a requisição recebe `400` e a gravação já realizada não é desfeita.
5. **Dado** uma requisição sem `x-cid`, **Quando** o evento é publicado, **Então** o `cid` usado é um UUID gerado pelo handler.

---

### História de Usuário 2 - Consumir o contrato do evento (Prioridade: P2)

Um consumidor precisa interpretar o conteúdo publicado sem conhecer o modelo interno do serviço.

**Por que esta prioridade**: define o contrato entre Ledger Config e os demais serviços.

**Teste independente**: chamar `ToConfigMessage` com uma configuração de exemplo e validar o JSON resultante.

**Cenários de aceitação**:

1. **Dado** uma configuração de domínio, **Quando** ela é convertida para mensagem, **Então** o JSON contém `config_id`, `processing_code`, `org_id`, `program_id`, `description`, `scripts`, `enable`, `created_at`, `updated_at` e `version`.
2. **Dado** o envelope do evento, **Quando** ele é publicado, **Então** contém `EventID` UUID, `EventType` igual a `ledger` e `EventDate` preenchido.
3. **Dado** uma configuração com nível `program`, **Quando** a mensagem é gerada, **Então** o campo de nível é serializado como `event_id` — nome atual a corrigir.
4. **Dado** uma configuração com scripts, **Quando** a lista de scripts é convertida, **Então** a mensagem inclui entradas zeradas no início devido a `make` seguido de `append` — comportamento atual a corrigir.

---

### Casos Limite

- `message.New(ctx)` é criado por request e fechado com `defer`; o fechamento apenas registra log.
- `CONFIG_SNS_TOPIC` vazio faz com que o publicador não tenha ARN válida; não há validação explícita na inicialização.
- Erros de publicação são propagados para o handler como erro de caso de uso, não como erro 5xx.
- Não há outbox, retry ou idempotency key própria; o `EventID` é gerado a cada publicação.

## Requisitos *(obrigatório)*

### Requisitos Funcionais

- **FR-001**: O sistema DEVE publicar um evento para cada mutação bem-sucedida de configuração.
- **FR-002**: O sistema DEVE usar o tópico configurado em `CONFIG_SNS_TOPIC`.
- **FR-003**: O sistema DEVE incluir `cid`, identificador UUID, tipo `ledger` e data no envelope.
- **FR-004**: O sistema DEVE converter a entidade de domínio para o contrato de mensagem antes de publicar.
- **FR-005**: O sistema DEVE propagar falhas de publicação ao chamador.
- **FR-006**: O sistema DEVE expor a publicação como porta `Topic` para permitir mocks nos testes.
- **FR-007**: O sistema DEVE corrigir a tag JSON do nível para `level`.
- **FR-008**: O sistema DEVE gerar a lista de scripts sem entradas zeradas.

### Entidades Principais

- **Event**: envelope com `event_id`, `event_type`, `event_data`, `event_date`.
- **ConfigMessage**: representação serializável de `Config`.
- **ScriptMessage / AccountMessage**: representações serializáveis de `Script` e `Account`.
- **Topic**: porta `Emit(ctx, cid, config)` implementada por `message.Topic`.

## Critérios de Sucesso *(obrigatório)*

### Resultados Mensuráveis

- **SC-001**: 100% das mutações bem-sucedidas geram exatamente um evento por configuração afetada.
- **SC-002**: 100% dos eventos possuem `event_type = ledger` e UUID válido no identificador.
- **SC-003**: Nenhuma mensagem contém campos com nomes divergentes do contrato acordado.
- **SC-004**: Falhas de publicação são sempre visíveis ao chamador.

## Suposições

- Consumidores leem mensagens diretamente do SNS/SQS sem transformação adicional.
- O `core-sdk` é responsável pela conexão e pelo `Health()` do publicador.
- O tópico e a inscrição SNS→SQS são criados pelo Terraform.
