# AGENTS.md

## Idioma

- Escreva respostas, explicações, mensagens de commit, código (comentários, logs, strings visíveis ao usuário) e artefatos em **português do Brasil**.
- Mensagens de commit seguem o padrão Conventional Commits em pt-BR.

## Estrutura do Repositório

- `core/` — implementação executável da plataforma.
  - `core/app/` — serviços internos (`auth-service`, `pokemon-catalog-service`).
  - `core/bff/` — BFFs orientados ao cliente (`mobile-bff`).
  - `core/gateway/` — configuração do Kong API Gateway.
  - `core/infra/` — infraestrutura compartilhada (PostgreSQL, Redis, seeds, migrations).
  - `core/docker-compose.yml` — compose principal da plataforma.
- `doc/` — documentação de arquitetura, decisões e padrões.
- `bruno/` — coleções de teste de API (Bruno).
- `.github/` — CI/CD, agents, skills e templates.
- `.specs/` — especificações de features (spec-driven development).

## Serviços

| Serviço | Módulo | Responsabilidade |
|---------|--------|-----------------|
| `mobile-bff` | `core/bff/mobile-bff` | Orquestração voltada ao cliente, composição de respostas, autenticação, favoritos |
| `pokemon-catalog-service` | `core/app/pokemon-catalog-service` | Fonte canônica do catálogo de Pokémon |
| `auth-service` | `core/app/auth-service` | Autenticação, ciclo de vida de tokens JWT, refresh, revogação |

## Arquitetura — mobile-bff (Hexagonal)

O `mobile-bff` segue Ports & Adapters (Hexagonal):

```
internal/
  domain/          ← entidades, value objects, erros de domínio
  ports/
    inbound/       ← interfaces de use cases
    outbound/      ← interfaces de repositórios e clientes externos
  service/         ← implementação dos use cases
  adapters/
    inbound/http/  ← handlers HTTP, DTOs, middleware
    outbound/http/ ← clients HTTP para serviços internos
    outbound/postgres/ ← repositórios PostgreSQL
  config/          ← configuração por variáveis de ambiente
  infrastructure/  ← logger, utilities
cmd/server/main.go ← entry point, wiring de dependências
```

**Regras arquiteturais:**
- Handlers HTTP dependem de `ports/inbound` (use cases), nunca de clients concretos.
- Use cases dependem de `ports/outbound` (interfaces), implementados pelos adapters.
- Erros externos são normalizados no adapter ou na camada de serviço, não no handler.
- Novas entidades de domínio vivem em `domain/`, não em `ports/`.
- O package `tests/` nunca é importado por código de produção.

## Padrão de Desenvolvimento

Features seguem o fluxo **spec-driven** (tlc-spec-driven):

```
Specify → Design → Tasks → Execute
```

Artefatos em `.specs/features/[feature]/`:
- `spec.md` — requisitos com IDs rastreáveis
- `design.md` — arquitetura e componentes
- `tasks.md` — tarefas atômicas com verificação
- `validation.md` — relatório do Verifier

Estado do projeto em `.specs/STATE.md` (Decisions + Handoff).

## Convenções Go

- **Receiver**: 1-2 letras. `c` para Client, `s` para Service. Nunca `this` ou `self`.
- **Sem prefixo Get**: `Count()` não `GetCount()`. Exceção: `http.ResponseWriter`.
- **Initialisms**: `ID`, `DB`, `URL`, `HTTP`, `JSON`. Nunca `Id`, `Db`, `Url`.
- **Error strings**: minúsculas, sem ponto final.
- **Indent error flow**: retorne o erro imediatamente.
- **Wrapping**: `%w` para erros inspecionáveis (`errors.Is`/`errors.As`), `%v` para anotação.
- **Slices via JSON**: use `make([]T, 0)` ao retornar slices via API/JSON.
- **Interface compliance**: `var _ Interface = (*Impl)(nil)` em cada adapter e service.
- **Constantes de tempo**: `const timeout = 30 * time.Second`.
- **Testes**: table-driven com `t.Run`. Cobertura mínima 75%, ideal 90%.

Para revisão de código Go, carregue a skill `go-style-combined` (`.github/skills/go-style-combined/SKILL.md`).

## Commits

- Formato: `tipo(escopo-opcional): descrição curta em português`
- Tipos: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `build`, `perf`, `revert`
- Um commit atômico por task concluída.
- Exemplo: `feat(bff): adicionar circuit breaker nos clients HTTP`

## Agents Disponíveis

Agentes especializados em `.github/agents/`:

| Agent | Uso |
|-------|-----|
| `project-planner.agent.md` | Planejar features, decompor tarefas, estrutura de arquivos |
| `backend-architect.agent.md` | Decisões de arquitetura, padrões, contratos |
| `security-auditor.agent.md` | Auditoria de segurança, revisão de vulnerabilidades |
| `database-architect.agent.md` | Schema, queries, migrations, performance |
| `test-engineer.agent.md` | Estratégia de testes, cobertura, mocks |
| `debugger.agent.md` | Diagnóstico de bugs, análise de causa raiz |
| `devops-engineer.agent.md` | Docker, CI/CD, infraestrutura |
| `documentation-writer.agent.md` | Documentação, README, docs de arquitetura |

## Skills Disponíveis

Skills Go em `.github/skills/`:

| Skill | Descrição |
|-------|-----------|
| `golang-testing` | Testes unitários, integração, mocks, HTTP testing |
| `golang-error-handling` | Criação, wrapping, tratamento de erros |
| `golang-database` | PostgreSQL, transações, performance, scanning |
| `golang-security` | Auditoria de segurança, secrets, injeção, cookies |
| `golang-linter` | golangci-lint, configuração, diretivas |
| `golang-documentation` | Templates de docs, comentários, README |
| `go-style-combined` | Uber + Google Go Style Guide |
| `go-api-design` | Design de APIs REST |
| `go-security-audit` | Auditoria de segurança |
| `go-test-quality` | Qualidade de testes |
| `go-test-table-driven` | Testes table-driven |
| `go-error-handling` | Tratamento de erros |
| `go-architecture-review` | Revisão de arquitetura |

## CI/CD

- `.github/workflows/go-ci.yml` — build, test, vet, lint
- `.github/workflows/conventional-commits.yml` — validação de commits
- `.github/workflows/pr-title-conventional.yml` — validação de títulos de PR
