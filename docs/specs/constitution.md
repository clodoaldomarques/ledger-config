# Constituição — Ledger Config

> Derivada da implementação em `60095c7`. Esta é uma constituição **descritiva**: registra o que o projeto já faz hoje.

## Princípios Fundamentais

### I. Separação em camadas (NÃO NEGOCIÁVEL)

O fluxo é `cmd/main.go` → `internal/infra/rest/server` → `internal/infra/rest/ledger` (handlers) → `internal/application/ledger` (casos de uso) → `internal/domain/ledger` (regras). Adapters de infraestrutura ficam em `internal/infra/{db,message,rest,gringotts}`. Handlers não contêm regras de negócio; o domínio não conhece HTTP, AWS ou banco.

### II. O domínio valida, a interface transporta

`Config.Validate()` e `Script.Validate()` são a fonte das regras (`processing_code`, níveis `tenant`/`program`, duplicidade por `flow#script_id`). DTOs de request (`internal/infra/rest/ledger/request.go`) apenas convertem JSON para entidades; a validação de formato (`validate` tag) é complementar, não substituta.

### III. Configuração por ambiente, com singleton explícito

`config.New()` usa `sync.Once`: variáveis de ambiente são capturadas uma vez e opções `WithX` alteram a mesma instância. O binário não lê `.env`; `make run` e o `launch.json` do VS Code fazem isso. Nenhum segredo é versionado.

### IV. Observabilidade é parte do request

Todo request passa pelos interceptadores de tracer, logger e métricas do `core-sdk`. `x-cid` identifica a correlação (gerado quando ausente) e `x-tenant` identifica a organização. Spans são criados em handler, service e repositório; erros são registrados com `span.SetError`.

### V. Consultas são desenhadas no modelo de dados

A tabela `ConfigTable` usa `org_id` como hash key, `config_id` como range key e a GSI `GSI-Index` com `filters`. O campo `filters` (`TENANT#...#PROGRAMID#...#EVENTTYPEID#...`) é construído no domínio da infraestrutura e é a base da precedência program → tenant e dos filtros de listagem.

### VI. Testes unitários determinísticos

A suíte usa `go test ./...`, `testify` e mocks `gomock` gerados por `mockgen` a partir de `interfaces.go`. `TestMain` instala um tracer noop, eliminando dependência de coletor OTEL. Não há testes de integração na suíte padrão.

### VII. Infraestrutura como código, com ambientes locais explícitos

`docker-compose.yaml`, `scripts/terraform/`, `scripts/k8s/` e `scripts/docker/` descrevem o ambiente. O deploy local depende de LocalStack e de um IP fixo no provider Terraform; mudanças nesses arquivos afetam o fluxo `make up`/`make apply`.

## Restrições e Convenções

- Go conforme `go.mod`; a imagem de deploy usa `golang:1.27.1-alpine` e binário estático (`CGO_ENABLED=0`).
- Handler cria repositório e publicador SNS por request e os fecha com `defer`.
- Erros de caso de uso são convertidos em `400` (ou `404` em consultas); o corpo de erro é `{"message": "..."}`.
- `Level` aceito na API é `tenant` ou `program`; `platform` só pode ser criado pela equipe.
- Nomes de benchmark e perfis devem ser verificados em `internal/application/ledger/bench_test.go` antes de usar o Makefile.

## Fluxo de Desenvolvimento

1. Especificar a mudança em `docs/specs/`.
2. Ajustar entidades e validação em `internal/domain/` quando a regra mudar.
3. Atualizar `interfaces.go` e rodar `go generate ./internal/application/ledger` se uma porta mudar.
4. Implementar caso de uso, handler e DTOs.
5. Rodar `go test ./...`; `make test` apenas quando desejar abrir o relatório HTML.
6. Verificar `go build ./...` e `go vet ./...` (não há CI configurada).

## Governança

Esta constituição é um registro do estado atual. Alterações de arquitetura ou de regras de domínio devem atualizar este arquivo, a spec afetada e o `plan.md` correspondente na mesma mudança. Complexidade adicional (por exemplo, novo adapter ou novo nível de configuração) precisa ser justificada na seção *Complexity Tracking* do plano.

**Versão**: 1.0 | **Ratificada**: 2026-11-21 | **Última alteração**: 2026-11-21
