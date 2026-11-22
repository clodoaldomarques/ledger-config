# Especificações do Ledger Config

Specs geradas por **engenharia reversa** do código em `main` (commit `60095c7`, 2026-11-21). Elas descrevem o comportamento **atual**, não um estado desejado.

## Estrutura

| Arquivo | Conteúdo |
|---|---|
| `constitution.md` | Princípios e restrições extraídos da forma como o projeto é implementado |
| `platform.md` | Contexto técnico: stack, configuração, DynamoDB, infra local, deploy e testes |
| `001-config-lifecycle/` | Criação, consulta, listagem, atualização e desativação de configurações |
| `002-tenant-activation/` | Cópia das configurações `LEDGER` para um novo tenant |
| `003-config-events/` | Publicação e contrato dos eventos SNS |
| `004-observability-health/` | Liveness, readiness, tracing, métricas e logs |

Cada funcionalidade segue o formato Spec Kit:

- `spec.md` — cenários de usuário, requisitos funcionais (FR), entidades e critérios de sucesso (SC).
- `plan.md` — decisões técnicas já presentes no código, estrutura de arquivos e complexidade conhecida.
- `tasks.md` — sequência de implementação/reconstrução.

## Convenção das tarefas

- `[x]` = comportamento **já implementado** no commit de referência.
- `[ ]` = lacuna ou melhoria identificada na engenharia reversa (bug, teste ausente, código não ligado).

Os identificadores `US`, `FR`, `SC` e `T` seguem a nomenclatura do Spec Kit para permitir rastreabilidade entre spec, plano e tarefas.

## Uso para estudo

1. Leia `constitution.md` e `platform.md`.
2. Escolha uma funcionalidade e leia `spec.md` antes de olhar o código.
3. Compare `plan.md` com os arquivos citados.
4. Use `tasks.md` para identificar o que foi implementado e o que ficou como exercício.

As specs são um espelho do código: se o código mudar, atualize também estes arquivos.
