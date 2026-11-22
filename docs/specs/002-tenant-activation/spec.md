# Especificação da Funcionalidade: Ativação de Tenant

**Branch da funcionalidade**: `002-tenant-activation`

**Criado**: 2026-11-21

**Status**: Rascunho (derivado do código em `60095c7`)

**Entrada**: engenharia reversa de `Service.ActivateTenant`, handler `ActivateTenant` e seed `scripts/terraform/configs.tf`

## Cenários de Usuário e Testes *(obrigatório)*

### História de Usuário 1 - Ativar um tenant (Prioridade: P1)

Um operador cadastra uma nova organização e precisa que ela receba as configurações padrão mantidas pela equipe no nível `platform`.

**Por que esta prioridade**: sem a cópia, o tenant recém-criado não tem regras de ledger aplicáveis.

**Teste independente**: com o seed `ledger-101` presente, enviar `POST /v1/ledger/org/activate` com `{"org_id":"TENANT_A"}` e verificar o retorno com uma configuração de nível `tenant`.

**Cenários de aceitação**:

1. **Dado** que existe ao menos uma configuração `platform` da organização `LEDGER`, **Quando** a ativação é solicitada para `TENANT_A`, **Então** o sistema responde `200` com uma lista de configurações copiadas, com novos `config_id`, nível `tenant`, `enable = true` e `version = 1`.
2. **Dado** que a cópia preserva `processing_code`, `program_id`, `description` e `scripts` da configuração original, **Quando** a ativação é concluída, **Então** a configuração do tenant contém os mesmos valores, exceto `org_id`, `config_id` e timestamps.
3. **Dado** que a publicação do evento falha durante a cópia, **Quando** a ativação é processada, **Então** o sistema responde `400` e as configurações já gravadas permanecem.
4. **Dado** um corpo sem `org_id`, **Quando** a ativação é solicitada, **Então** o sistema responde `400` com `OrgID is required`.

---

### História de Usuário 2 - Evitar duplicação em reativações (Prioridade: P2)

Um operador repete a ativação e não deve gerar configurações duplicadas.

**Por que esta prioridade**: a operação pode ser reexecutada após falha parcial.

**Teste independente**: ativar o mesmo tenant duas vezes e conferir que a quantidade de configurações não aumenta.

**Cenários de aceitação**:

1. **Dado** que o tenant já possui a configuração equivalente, **Quando** a ativação é repetida, **Então** nenhuma nova configuração é gravada e nenhum evento é publicado.
2. **Dado** que nenhuma configuração nova foi criada, **Quando** a ativação termina, **Então** o sistema responde `400` com `tenant <org_id> was activated` — comportamento atual a revisar.
3. **Dado** que existem várias configurações `platform`, **Quando** algumas já foram copiadas, **Então** apenas as faltantes são criadas.

---

### Casos Limite

- A origem é fixa: `FindAllConfigs(ctx, cid, "LEDGER", nil)`; não existe parâmetro para escolher outra organização.
- `created_at` e `updated_at` da cópia usam `time.Now()` (sem `.UTC()`), diferente dos demais casos de uso.
- Configurações `platform` desabilitadas também servem de origem: a listagem não filtra `enable`.
- Não há autorização: qualquer chamador com `x-tenant` pode informar outro `org_id` no corpo.
- erros de publicação não desfazem gravações anteriores.

## Requisitos *(obrigatório)*

### Requisitos Funcionais

- **FR-001**: O sistema DEVE obter todas as configurações da organização `LEDGER` como modelo.
- **FR-002**: O sistema DEVE exigir `org_id` no corpo da requisição.
- **FR-003**: O sistema DEVE criar uma cópia por configuração modelo com novo UUID e nível `tenant`.
- **FR-004**: O sistema DEVE preservar código de processamento, programa, descrição e scripts.
- **FR-005**: O sistema DEVE iniciar cada cópia com `enable = true` e `version = 1`.
- **FR-006**: O sistema DEVE verificar a existência da configuração no tenant antes de gravar.
- **FR-007**: O sistema DEVE publicar um evento para cada cópia gravada.
- **FR-008**: O sistema DEVE sinalizar quando nenhuma cópia foi criada.

### Entidades Principais

- **PostOrgActivate**: DTO de entrada com `org_id`.
- **Configuração modelo**: item `platform` pertencente a `LEDGER`.
- **Configuração copiada**: item `tenant` pertencente ao `org_id` informado, com `filters` recalculado.

## Critérios de Sucesso *(obrigatório)*

### Resultados Mensuráveis

- **SC-001**: Após uma ativação com sucesso, 100% das configurações modelo têm uma cópia correspondente no tenant.
- **SC-002**: Nenhuma cópia duplicada é criada em duas ativações consecutivas.
- **SC-003**: 100% das cópias retornadas têm nível `tenant`, `enable = true` e `version = 1`.
- **SC-004**: Um evento é publicado para cada cópia gravada, sem eventos para configurações ignoradas.

## Suposições

- O seed `ledger-101` (org `LEDGER`, nível `platform`) existe.
- `x-cid` é usado apenas para correlação; a identidade do tenant vem do corpo.
- A lista de configurações modelo cabe em uma única página de consulta DynamoDB.
