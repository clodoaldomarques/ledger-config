# Especificação da Funcionalidade: Observabilidade e Health Checks

**Branch da funcionalidade**: `004-observability-health`

**Criado**: 2026-11-21

**Status**: Rascunho (derivado do código em `60095c7`)

**Entrada**: engenharia reversa de `cmd/main.go`, `internal/infra/rest/server/server.go`, `config/config.go` e `scripts/docker/`

## Cenários de Usuário e Testes *(obrigatório)*

### História de Usuário 1 - Verificar se o processo está vivo (Prioridade: P1)

Orquestradores precisam de um endpoint leve para reiniciar o container quando o processo HTTP parou.

**Por que esta prioridade**: é o menor sinal de vida e não depende de infraestrutura externa.

**Teste independente**: iniciar o servidor sem DynamoDB e chamar `GET /health/live`.

**Cenários de aceitação**:

1. **Dado** que o servidor HTTP está escutando, **Quando** `GET /health/live` é chamado, **Então** o sistema responde `200` sem corpo.
2. **Dado** que o DynamoDB está indisponível, **Quando** `GET /health/live` é chamado, **Então** o sistema continua respondendo `200`.

---

### História de Usuário 2 - Verificar prontidão das dependências (Prioridade: P1)

O orquestrador só deve enviar tráfego quando o serviço consegue acessar DynamoDB e SNS.

**Por que esta prioridade**: evitar que a API aceite requisições que falharão por dependências indisponíveis.

**Teste independente**: com LocalStack parado, chamar `GET /health/ready` e verificar `503`; com LocalStack ativo, verificar `200`.

**Cenários de aceitação**:

1. **Dado** que `DescribeTable` e o `Health()` do publicador SNS são bem-sucedidos, **Quando** `GET /health/ready` é chamado, **Então** o sistema responde `200`.
2. **Dado** que a consulta ao DynamoDB falha, **Quando** `GET /health/ready` é chamado, **Então** o sistema responde `503`.
3. **Dado** que o `Health()` do SNS falha, **Quando** `GET /health/ready` é chamado, **Então** o sistema responde `503`.
4. **Dado** que as dependências demoram mais que 500 ms, **Quando** o contexto expira, **Então** o sistema responde `503`.

---

### História de Usuário 3 - Rastrear uma requisição (Prioridade: P2)

Um operador precisa correlacionar logs, spans e eventos de uma mesma requisição.

**Por que esta prioridade**: reduz o tempo de diagnóstico em chamadas que atravessam handler, serviço, repositório e SNS.

**Teste independente**: enviar `x-cid: abc-123` e verificar que o valor aparece nos atributos dos spans e nos logs da requisição.

**Cenários de aceitação**:

1. **Dado** um `x-cid` informado, **Quando** a requisição é processada, **Então** o mesmo valor é usado nos spans de handler, serviço, repositório e publicador.
2. **Dado** uma requisição sem `x-cid`, **Quando** ela é processada, **Então** um UUID é gerado e usado como identificador de correlação.
3. **Dado** um `x-tenant`, **Quando** a requisição é processada, **Então** o valor é registrado como `org_id` no span do handler.
4. **Dado** o recebimento de `SIGTERM`, **Quando** o processo encerra, **Então** o servidor HTTP faz shutdown gracefully em até 25 s e o OpenTelemetry é desligado.

---

### História de Usuário 4 - Coletar métricas e logs (Prioridade: P2)

A equipe precisa observar latência, erros e volume de requisições sem instrumentar cada handler manualmente.

**Por que esta prioridade**: os interceptadores centrais cobrem todas as rotas automaticamente.

**Teste independente**: com o coletor OTEL do Compose ativo, chamar uma rota e verificar que métricas chegam ao Prometheus e spans ao Zipkin.

**Cenários de aceitação**:

1. **Dado** o coletor OTEL configurado, **Quando** uma rota é chamada, **Então** spans e métricas são exportados.
2. **Dado** um erro de validação, **Quando** o handler responde, **Então** o erro é registrado no span correspondente.
3. **Dado** `OTEL_SERVICE_NAME` definido, **Quando** a telemetria é exportada, **Então** o serviço aparece com esse nome.

---

### Casos Limite

- `otel.Start(ctx)` é chamado antes de `server.New()`; o comportamento sem coletor depende do `core-sdk` (não há configuração explícita no repositório).
- O timeout de 500 ms do readiness é fixo em código, não configurável.
- `x-cid` e `x-tenant` não são validados quanto a formato.
- Não há testes automatizados para os endpoints de health nem para os interceptadores.

## Requisitos *(obrigatório)*

### Requisitos Funcionais

- **FR-001**: O sistema DEVE expor `GET /health/live` sem depender de serviços externos.
- **FR-002**: O sistema DEVE expor `GET /health/ready` verificando DynamoDB e SNS.
- **FR-003**: O sistema DEVE responder `503` quando qualquer verificação de prontidão falhar.
- **FR-004**: O sistema DEVE limitar a verificação de prontidão a 500 ms.
- **FR-005**: O sistema DEVE aplicar interceptadores de tracer, logger e métricas a todas as rotas.
- **FR-006**: O sistema DEVE usar `x-cid` como correlação e gerar UUID quando ausente.
- **FR-007**: O sistema DEVE registrar `org_id` a partir de `x-tenant` nos spans.
- **FR-008**: O sistema DEVE encerrar HTTP e OpenTelemetry de forma graceful em até 25 s.

### Entidades Principais

- **Health endpoint**: resposta HTTP sem corpo, com status `200` ou `503`.
- **Correlation context**: par `(cid, org_id)` extraído dos headers e propagado por spans.
- **Telemetry exporter**: coletor OTLP configurado por variáveis de ambiente.

## Critérios de Sucesso *(obrigatório)*

### Resultados Mensuráveis

- **SC-001**: 100% das chamadas a `/health/live` são respondidas sem consultar DynamoDB ou SNS.
- **SC-002**: `/health/ready` responde em no máximo 500 ms em qualquer cenário.
- **SC-003**: 100% das requisições possuem `cid` encontrável nos spans e logs.
- **SC-004**: Nenhuma rota deixa de ser instrumentada quando uma nova é adicionada ao roteador.

## Suposições

- O `core-sdk` fornece os interceptadores e a inicialização OTEL.
- Prometheus e Zipkin são usados apenas no ambiente local.
- O orquestrador de deploy usa os probes declarados em `scripts/k8s/deployment.yaml`.
