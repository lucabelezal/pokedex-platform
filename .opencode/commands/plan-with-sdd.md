---
description: Planeja e executa feature via Spec-Driven Development (Specify → Design → Tasks → Execute) com validacao deterministica
---

Feature: $ARGUMENTS

> Se $ARGUMENTS estiver vazio, pergunte qual feature planejar antes de prosseguir.

## Skill principal

Carregue **tlc-spec-driven** (`.claude/skills/tlc-spec-driven/SKILL.md` ou `.agents/skills/tlc-spec-driven/`).

Regras criticas dessa skill — nao negocie:

- Artefatos em `.specs/features/[feature]/` (spec.md, design.md, tasks.md, validation.md) + estado em `.specs/STATE.md`
- Execute `python3 <skill-dir>/scripts/validate_spec.py` antes de confirmar spec, `validate_tasks.py` antes de apresentar tasks, `validate_state.py` ao fechar
- Uma task = um commit atomico Conventional Commit pt-BR. Marque task como done em `tasks.md` antes do commit
- Verifier independente (author != verifier) roda apos ultima task — spec-anchored check + discrimination sensor

## Fluxo

### 0) Pre-voo

1. Leia `.specs/STATE.md` (Decisions + Handoff). Reconcilie com `git status --porcelain`, branch atual e ultimos commits.
2. Carregue lessons confirmadas: `python3 <skill-dir>/scripts/lessons.py list --status confirmed`
3. Aplique `harness/eval` como referencia de qualidade — skills Go do projeto devem ser respeitadas na implementacao:
   - `go-style-combined` (canonica, Uber+Google + decisoes locais)
   - `go-architecture-review` (hexagonal: handlers→ports/inbound, service→ports/outbound, domain sem deps)
   - `go-error-handling`, `go-api-design`, `go-security-audit`, `go-test-quality`, `go-test-table-driven`
   - Convencoes `AGENTS.md` validadas por `harness/evals/agents-md.json` (`python3 harness/eval/run.py schema`)

### 1) Specify — WHAT (sempre)

- Siga `references/specify.md` da skill (EARS, dimensions sweep, closure gate)
- Pergunte conversacionalmente (problema, quem e a dor, sucesso, constraints, out-of-scope)
- Escreva ACs em EARS: Ubiquitous / WHEN / WHILE / WHERE / IF-THEN — cada AC com SHALL e outcome preciso
- Atribua Requirement IDs (`[PREFIX]-NN`) e trace em `.specs/features/[feature]/spec.md`
- Rode `python3 <skill-dir>/scripts/validate_spec.py .specs/features/[feature]` e corrija ate passar
- Apresente spec para confirmacao antes de avancar

### 2) Design — HOW (auto-size: pule se Small/Medium)

- Se Large/Complex: siga `references/design.md` — arquitetura, componentes, contratos, decisoes
- Consulte `go-architecture-review` e `go-api-design` para decisoes de layout Go / API REST
- Escreva `.specs/features/[feature]/design.md`

### 3) Tasks — BREAKDOWN (auto-size: pule se Small/Medium)

- Siga `references/tasks.md` — tasks atomicas com `Tests` e `Gate` cada
- Rode `python3 <skill-dir>/scripts/validate_tasks.py .specs/features/[feature]`
- Se > ~8 tasks, ofereca execucao com sub-agents (batches de ~7 tasks, fases intactas)
- Apresente `tasks.md` para aprovacao

### 4) Execute — BUILD

- Siga `references/implement.md` — leia completamente antes de codar
- Para cada task em ordem:
  1. Derive testes dos ACs (assert outcome do spec, nao da implementacao)
  2. Implemente
  3. Gate: `go test ./...` / `make test` (harness) deve passar — runner decide, nao auto-avaliacao
  4. Marque task done em `tasks.md` + trace Requirement ID
  5. Commit atomico: `python3 <skill-dir>/scripts/check_commit.py --message "<msg>"` deve passar; formato `tipo(escopo): descricao em pt-BR`
- Apos ultima task: dispare Verifier (author != verifier) → `validation.md` (PASS/FAIL, evidencia `file:line`, sensor, diff range)
- Rode `python3 <skill-dir>/scripts/validate_state.py [feature]` — deve ser PASS com evidencia
- Distile lessons: `python3 <skill-dir>/scripts/lessons.py` para falhas grounded

## Harness do projeto

- Testes: `make test` (unit, integracao auto-skip), `make test-integration` (sobe `core/docker-compose.test.yml`), `make lint` (`gofmt`+`go vet`), `make ci`
- Eval: `python3 harness/eval/run.py list|schema|run --dry-run` ou `make eval`; para LLM: `EVAL_API_KEY=... python3 harness/eval/run.py run`
- CI deve permanecer verde: `core/app/pokemon-catalog-service`, `core/app/auth-service`, `core/bff/mobile-bff`, `core/infra/postgres/json2sql`

## Ao retomar trabalho interrompido

1. Leia `.specs/STATE.md` e reconcilie com git (branch, porcelain, commits)
2. Proponha proximo passo reconciliado antes de codar — evidencia vence snapshot stale
