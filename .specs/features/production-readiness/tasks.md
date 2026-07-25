# Production Readiness - Tasks

## Execution Protocol (MANDATORY -- do not skip)

Implement these tasks with the `tlc-spec-driven` skill: **activate it by name and follow its Execute flow and Critical Rules.** Do not search for skill files by filesystem path. The skill is the source of truth for the full flow (per-task cycle, sub-agent delegation, adequacy review, Verifier, discrimination sensor).

**If the skill cannot be activated, STOP and tell the user — do not proceed without it.**

---

**Design**: `.specs/features/production-readiness/design.md`
**Status**: Approved

---

## Test Coverage Matrix

> Generated from codebase sampling: `core/bff/mobile-bff/Makefile`, CI workflow `.github/workflows/go-ci.yml`, existing tests in `tests/unit/` e `tests/integration/`. Guidelines found: `AGENTS.md` (cobertura 75% min, 90% ideal, table-driven), `Makefile` (go test -v -race), `.github/copilot-instructions.md` (convenções Go).

| Code Layer | Required Test Type | Coverage Expectation | Location Pattern | Run Command |
|------------|-------------------|---------------------|-----------------|-------------|
| Domain / service (use case) | unit | Todas as branches; 1:1 com ACs da spec; todos os edge cases listados | `tests/unit/*_test.go` ou `internal/**/*_test.go` | `go test -v -race ./tests/unit` ou `go test -v -race ./internal/...` |
| HTTP Handler / adapter inbound | unit | Happy path + edge cases + error paths para cada endpoint modificado | `tests/unit/handlers_test.go` ou `internal/**/*_test.go` | `go test -v -race ./tests/unit` |
| Repository / adapter outbound (Postgres) | integration | Key query paths + error handling; usar real PostgreSQL | `tests/integration/*_test.go` | `DATABASE_URL=... go test -v -race ./tests/integration` |
| Outbound HTTP client | unit | Mock do servidor remoto; happy path + timeouts + erros 4xx/5xx | `tests/unit/*_test.go` | `go test -v -race ./tests/unit` |
| Middleware | unit | Todos os headers/cookies/status codes para cada cenário | `tests/unit/middleware_test.go` | `go test -v -race ./tests/unit` |
| Config / DTO / entity | none | — (build gate only) | — | `go build ./...` + `go vet ./...` |

## Gate Check Commands

> Generated from `Makefile` e `go-ci.yml`.

| Gate Level | When to Use | Command |
|------------|-------------|---------|
| Quick | Após tasks com apenas unit tests | `go test -v -race ./tests/unit` (ou `./internal/...` para serviços sem tests/) |
| Full | Após tasks com integration tests | `go test -v -race ./tests/unit ./tests/integration -timeout 30s` |
| Build | Após conclusão de phase ou tasks de config/entity | `go build ./... && go vet ./... && go fmt ./...` |
| Lint | Validação final de cada serviço | `go fmt ./... && go vet ./...` |
| Coverage | Verificação de cobertura | `go test -v -race -coverprofile=coverage.out ./tests/unit ./tests/mocks ./internal/...` |

---

## Execution Plan

Phases são ordenadas e executadas sequencialmente — tasks dentro de cada phase em ordem. Tasks P1 primeiro, P2 em seguida, P3 por último.

```
Phase 1: Foundation (T1 → T2 → T3 → T4 → T5 → T6)      ~6 tasks
Phase 2: Security (T7 → T8 → T9)                           ~3 tasks
Phase 3: Favorites Isolation (T10 → T11 → T12)             ~3 tasks
Phase 4: Observability (T13 → T14 → T15 → T16 → T17 → T18) ~6 tasks
Phase 5: Bug Fixes (T19 → T20 → T21 → T22 → T23)          ~5 tasks
Phase 6: Tests (T24 → T25 → T26)                           ~3 tasks
Phase 7: Infra + Code Quality (T27 → T28 → T29 → T30)      ~4 tasks
```

**Batching:** ~30 tasks → 3 batches de ~10 tasks cada (Batch 1 = Phases 1-2, Batch 2 = Phases 3-5, Batch 3 = Phases 6-7). Verifier roda após Batch 3.

---

## Task Breakdown

### T1: Centralizar Type Colors em domain/type_colors.go

**What**: Criar arquivo único `domain/type_colors.go` com o mapa de type colors (português + inglês), e substituir todas as 5 cópias existentes por imports deste arquivo.
**Where**: `core/bff/mobile-bff/internal/domain/type_colors.go` (novo)
**Depends on**: None
**Reuses**: valores do mapa existente em `response_builder.go` (mais completo com aliases PT)
**Requirements**: PR-13

**Tools**:
- MCP: none
- Skill: `go-style-combined`

**Done when**:
- [ ] Arquivo `type_colors.go` criado com `func TypeColor(typeName string) string` e mapa completo (PT + EN)
- [ ] `getTypeColor()` em `service/pokemon_service.go` substituído por chamada ao `domain.TypeColor()`
- [ ] `getTypeColor()` em `adapters/inbound/http/response_builder.go` substituído por chamada ao `domain.TypeColor()`
- [ ] `postgresTypeColor()` em `adapters/outbound/postgres/pokemon_repository.go` substituído
- [ ] `mockTypeColor()` em `tests/mocks/mock_repositories.go` substituído
- [ ] `typeColor()` em `catalog-service/repository/pokemon_repository.go` substituído (copiar type_colors.go para o catalog-service)
- [ ] `grep -r "getTypeColor\|postgresTypeColor\|mockTypeColor\|typeColor" --include="*.go" core/` retorna apenas definições em `domain/type_colors.go`
- [ ] Build gate passes: `go build ./...`

**Tests**: unit — testes para `TypeColor()` com tipos válidos, inválidos, case-insensitive
**Gate**: build

---

### T2: Criar adapters/outbound/memory/ e remover import de tests/mocks do main.go

**What**: Mover `MockPokemonRepository` e `MockFavoriteRepository` de `tests/mocks/` para `internal/adapters/outbound/memory/`. Atualizar `cmd/server/main.go` para importar do novo package.
**Where**: `core/bff/mobile-bff/internal/adapters/outbound/memory/` (novo)
**Depends on**: T1
**Reuses**: código existente de `tests/mocks/mock_repositories.go`
**Requirements**: PR-14

**Tools**:
- MCP: none
- Skill: `go-style-combined`

**Done when**:
- [ ] `memory/pokemon_repository.go` criado com implementação de `ports/outbound/PokemonRepository`
- [ ] `memory/favorite_repository.go` criado com implementação de `ports/outbound/FavoriteRepository`
- [ ] `memory/mock_repository.go` mantém ambos os mocks juntos
- [ ] `cmd/server/main.go` importa `adapters/outbound/memory` em vez de `tests/mocks`
- [ ] `tests/unit/*_test.go` e `tests/mocks/` usam type alias ou re-export do package memory
- [ ] `grep -r "tests/mocks" --include="*.go" core/bff/mobile-bff/cmd/` não retorna resultado
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — atualizar imports nos testes existentes; testes continuam passando
**Gate**: build

---

### T3: Criar CircuitBreakerClient wrapper com gobreaker

**What**: Criar wrapper `CircuitBreakerClient` que encapsula `http.Client` com circuit breaker (sony/gobreaker) + retries com backoff exponencial.
**Where**: `core/bff/mobile-bff/internal/adapters/outbound/http/circuit_breaker.go` (novo)
**Depends on**: T2
**Reuses**: `http.Client` padrão do Go; assinatura `Do(req *http.Request) (*http.Response, error)`
**Requirements**: PR-01, PR-02, PR-03

**Tools**:
- MCP: `context7` para `sony/gobreaker`
- Skill: `go-error-handling`

**Done when**:
- [ ] `CircuitBreakerClient` struct criada com `Do()` método
- [ ] Configuração: `MaxRequests=1`, `Interval=60s`, `Timeout=30s`, `FailureCount=5`
- [ ] Retry com backoff: 3 tentativas, backoff [1s, 3s, 10s]
- [ ] Jitter adicionado ao backoff (±10%)
- [ ] Logging via slog em cada transição de estado do circuit breaker
- [ ] `var _ http.RoundTripper = (*CircuitBreakerClient)(nil)` compile-time check (se aplicável)
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — teste com mock HTTP server simulando falhas transitórias e abertura/fechamento do circuito
**Gate**: quick

---

### T4: Integrar CircuitBreakerClient nos outbound clients (PokemonCatalogClient + AuthServiceClient)

**What**: Substituir `http.Client` direto por `CircuitBreakerClient` em `PokemonCatalogClient` e `AuthServiceClient`.
**Where**: `core/bff/mobile-bff/internal/adapters/outbound/http/pokemon_catalog_client.go`, `auth_service_client.go`
**Depends on**: T3
**Reuses**: `CircuitBreakerClient` de T3
**Requirements**: PR-01, PR-02

**Tools**:
- MCP: none
- Skill: `go-error-handling`

**Done when**:
- [ ] `PokemonCatalogClient` usa `CircuitBreakerClient` em vez de `http.Client` direto
- [ ] `AuthServiceClient` usa `CircuitBreakerClient` em vez de `http.Client` direto
- [ ] Configuração de timeout mantida (5s catalog, 10s auth)
- [ ] Circuit breaker nomeado para métricas: `"pokemon-catalog-service"` e `"auth-service"`
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — atualizar testes existentes de `auth_client_test.go` para mockar o CircuitBreakerClient; adicionar cenários de circuit breaker aberto
**Gate**: quick

---

### T5: Adicionar resposta de degradação controlada (503) no BFF

**What**: Quando circuit breaker está aberto (serviço indisponível), handlers retornam 503 com JSON `{"error": "...", "degraded": true}` em vez de erro genérico.
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/handler.go` (extensão)
**Depends on**: T4
**Reuses**: `dto.ErrorResponse` existente
**Requirements**: PR-04, PR-05

**Tools**:
- MCP: none
- Skill: `go-error-handling`

**Done when**:
- [ ] Função helper `writeDegradedResponse(w, "serviço de catálogo temporariamente indisponível")` criada
- [ ] Handlers de Pokémon verificam erro de circuit breaker e retornam 503
- [ ] Handlers de auth verificam erro de circuit breaker e retornam 503
- [ ] Response inclui campo `"degraded": true` no JSON
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — testes nos handlers existentes com mock de `PokemonUseCase` retornando erro de circuit breaker
**Gate**: quick

---

### T6: Remover fallback silencioso de 2 Pokémon no catalog-service

**What**: Remover `InMemoryPokemonRepository` como fallback no `cmd/server/main.go` do catalog-service. Quando DB falha, serviço retorna erro explícito.
**Where**: `core/app/pokemon-catalog-service/cmd/server/main.go`
**Depends on**: None (independente)
**Reuses**: N/A (remoção)
**Requirements**: PR-15

**Tools**:
- MCP: none
- Skill: `go-error-handling`

**Done when**:
- [ ] `InMemoryPokemonRepository` removido do `main.go` do catalog-service
- [ ] `main.go` retorna `log.Fatal` se não conseguir conectar ao banco
- [ ] Build gate passes: `go build ./...`

**Tests**: none — remoção de código; build gate basta
**Gate**: build

---

### T7: Dockerfiles seguros — USER 1000, SHA pinning, HEALTHCHECK

**What**: Atualizar todos os 3 Dockerfiles com práticas de segurança: `USER 1000` no runtime, imagens pinadas por SHA digest, HEALTHCHECK nos que faltam.
**Where**: 
- `core/bff/mobile-bff/Dockerfile`
- `core/app/auth-service/Dockerfile`
- `core/app/pokemon-catalog-service/Dockerfile`
**Depends on**: None (independente)
**Requirements**: PR-06, PR-07, PR-10

**Tools**:
- MCP: none
- Skill: none

**Done when**:
- [ ] mobile-bff Dockerfile: `FROM alpine:latest` substituído por SHA digest, `USER 1000` adicionado, HEALTHCHECK já existe — verificar
- [ ] auth-service Dockerfile: `FROM golang:alpine` → SHA digest, `USER 1000`, HEALTHCHECK adicionado (`wget http://localhost:8082/health`)
- [ ] catalog-service Dockerfile: `FROM golang:alpine` → SHA digest, `USER 1000`, HEALTHCHECK adicionado (`wget http://localhost:8081/health`)
- [ ] Docker compose build test: `docker compose -f core/docker-compose.yml build`

**Tests**: none — verificação manual via `docker inspect` ou build gate
**Gate**: build

---

### T8: Remover credenciais hardcoded do database.go

**What**: Remover string de conexão hardcoded `"postgres://user:password@localhost:5432/pokedex"` do `database.go`. Se `DATABASE_URL` não estiver configurada, retornar nil sem panic.
**Where**: `core/bff/mobile-bff/internal/adapters/outbound/postgres/database.go`
**Depends on**: None (independente)
**Requirements**: PR-08

**Tools**:
- MCP: none
- Skill: `go-error-handling`

**Done when**:
- [ ] `NewConnection()` retorna `nil, nil` se `DATABASE_URL` vazia (sem credenciais hardcoded)
- [ ] Código que chama `NewConnection()` já trata nil corretamente (fallback para memory/)
- [ ] `grep -r "postgres://user:password\|user:password@" --include="*.go" core/` não retorna resultado
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — teste de `NewConnection()` com DATABASE_URL vazia e com URL válida
**Gate**: quick

---

### T9: Adicionar SecurityHeadersMiddleware no BFF

**What**: Middleware que adiciona `Strict-Transport-Security`, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY` em todas as respostas.
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/middleware.go` (extensão)
**Depends on**: None (independente)
**Requirements**: PR-09

**Tools**:
- MCP: none
- Skill: `go-security-audit`

**Done when**:
- [ ] Função `SecureHeadersMiddleware(next http.Handler) http.Handler` criada
- [ ] Headers: `Strict-Transport-Security: max-age=31536000; includeSubDomains`, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`
- [ ] Middleware registrado no `handler.go` antes dos demais middlewares
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — teste HTTP que verifica presença dos 3 headers em resposta do BFF
**Gate**: quick

---

### T10: Criar endpoint GET /v1/pokemons/favorites no catalog-service

**What**: Novo handler no catalog-service que aceita query param `ids` (lista separada por vírgula, máx 100) e retorna array de Pokémon.
**Where**: `core/app/pokemon-catalog-service/internal/http/handlers.go` (extensão)
**Depends on**: None (independente)
**Requirements**: PR-12

**Tools**:
- MCP: none
- Skill: `go-api-design`

**Done when**:
- [ ] Handler `GetFavoritesBatch` implementado
- [ ] Query param `ids` parseado com `strings.Split`, validado (máx 100)
- [ ] Busca em batch via `GetByIDs(ctx, ids)` no repository
- [ ] IDs não encontrados omitidos (não causam erro)
- [ ] Lista vazia retorna `[]` (não `null`)
- [ ] Rota registrada no `main.go` do catalog-service
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — teste do handler com mock repository, cobrindo: IDs válidos, IDs parcialmente inválidos, lista vazia, >100 IDs
**Gate**: quick

---

### T11: Criar FavoriteServiceClient no BFF como adapter HTTP

**What**: Novo outbound adapter que chama `GET /v1/pokemons/favorites?ids=...` no catalog-service para buscar detalhes de Pokémon favoritos.
**Where**: `core/bff/mobile-bff/internal/adapters/outbound/http/favorite_service_client.go` (novo)
**Depends on**: T10
**Reuses**: padrão de `pokemon_catalog_client.go` (timeout, config, error wrapping)
**Requirements**: PR-11, PR-22

**Tools**:
- MCP: none
- Skill: `go-api-design`

**Done when**:
- [ ] `FavoriteServiceClient` struct criada implementando `ports/outbound/FavoriteProvider`
- [ ] Método `GetFavoriteDetails(ctx, ids []string) ([]domain.Pokemon, error)`
- [ ] Config via env `POKEMON_CATALOG_SERVICE_URL` (já existe)
- [ ] Timeout de 5s
- [ ] `var _ outbound.FavoriteProvider = (*FavoriteServiceClient)(nil)`
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — mock HTTP server simulando catalog-service, testes para happy path, erro 4xx, erro 5xx, timeout
**Gate**: quick

---

### T12: Migrar handlers de favoritos do BFF para usar FavoriteServiceClient (remover acesso direto ao DB)

**What**: Atualizar `favorite_handler.go` e `favorite_service.go` para usar `FavoriteServiceClient` em vez de `PostgresFavoriteRepository`. Remover dependência de `ports/outbound/FavoriteRepository` do BFF.
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/favorite_handler.go`, `internal/service/favorite_service.go`
**Depends on**: T11
**Requirements**: PR-11

**Tools**:
- MCP: none
- Skill: `go-architecture-review`

**Done when**:
- [ ] `GetUserFavorites` handler chama `FavoriteServiceClient.GetFavoriteDetails()` em batch
- [ ] N+1 query eliminada (uma única chamada HTTP para todos os favoritos)
- [ ] `PostgresFavoriteRepository` não é mais usado pelo BFF
- [ ] `cmd/server/main.go` não injeta mais `FavoriteRepository` do postgres
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — atualizar testes de `favorite_service_test.go` e `handlers_test.go` para usar mock do `FavoriteProvider` em vez de `FavoriteRepository`
**Gate**: quick

---

### T13: Adicionar OpenTelemetry SDK + tracing no mobile-bff

**What**: Configurar OTel SDK com OTLP HTTP exporter para Jaeger. Adicionar middleware de tracing e propagação W3C TraceContext.
**Where**: `core/bff/mobile-bff/internal/infrastructure/observability/tracing.go` (novo)
**Depends on**: T2
**Requirements**: PR-16

**Tools**:
- MCP: `context7` para `go.opentelemetry.io/otel`
- Skill: none

**Done when**:
- [ ] `InitTracerProvider(endpoint string) (*sdktrace.TracerProvider, error)` criado
- [ ] Middleware `TracingMiddleware` que extrai trace context dos headers de entrada e cria span
- [ ] Propagação W3C TraceContext nos outbound HTTP clients (adicionar headers `traceparent`)
- [ ] Config via env `OTEL_EXPORTER_OTLP_ENDPOINT` (default: `http://localhost:4318`)
- [ ] Shutdown do TracerProvider no `main.go` via `defer`
- [ ] Gate check passes: `go build ./...`

**Tests**: none — tracing é infraestrutura observável, não testável unitariamente
**Gate**: build

---

### T14: Adicionar OpenTelemetry SDK + tracing no catalog-service

**What**: Mesma configuração OTel do T13 aplicada ao pokemon-catalog-service.
**Where**: `core/app/pokemon-catalog-service/internal/observability/tracing.go` (novo)
**Depends on**: T10
**Requirements**: PR-16

**Tools**:
- MCP: `context7` para `go.opentelemetry.io/otel`
- Skill: none

**Done when**:
- [ ] `InitTracerProvider` criado
- [ ] Middleware de tracing nos handlers HTTP
- [ ] Propagação W3C (extrai `traceparent` dos headers de entrada)
- [ ] Shutdown no `main.go`
- [ ] Gate check passes: `go build ./...`

**Tests**: none
**Gate**: build

---

### T15: Adicionar OpenTelemetry SDK + tracing no auth-service

**What**: Mesma configuração OTel aplicada ao auth-service.
**Where**: `core/app/auth-service/internal/observability/tracing.go` (novo)
**Depends on**: None (independente)
**Requirements**: PR-16

**Tools**:
- MCP: `context7` para `go.opentelemetry.io/otel`
- Skill: none

**Done when**:
- [ ] `InitTracerProvider` criado
- [ ] Middleware de tracing nos handlers HTTP
- [ ] Propagação W3C (extrai `traceparent` dos headers de entrada)
- [ ] Shutdown no `main.go`
- [ ] Gate check passes: `go build ./...`

**Tests**: none
**Gate**: build

---

### T16: Adicionar métricas Prometheus nos 3 serviços

**What**: Expor endpoint `GET /metrics` com métricas HTTP (request count, latency histogram, errors) em todos os serviços. Middleware que coleta métricas por método, path e status.
**Where**:
- `core/bff/mobile-bff/internal/adapters/inbound/http/middleware.go` (extensão)
- `core/app/pokemon-catalog-service/internal/http/` (novo middleware.go)
- `core/app/auth-service/internal/http/` (novo middleware.go)
**Depends on**: T13, T14, T15
**Requirements**: PR-17

**Tools**:
- MCP: `context7` para `prometheus/client_golang`
- Skill: none

**Done when**:
- [ ] Middleware `MetricsMiddleware` criado em cada serviço
- [ ] Métricas: `http_requests_total{method, path, status}`, `http_request_duration_seconds{method, path}`, `http_errors_total{method, path, error_type}`
- [ ] Endpoint `/metrics` registrado (não roteado via Kong — interno)
- [ ] Gate check passes: `go build ./...`

**Tests**: none — métricas verificáveis via curl em `/metrics`
**Gate**: build

---

### T17: Unificar logging no auth-service (log.Printf → log/slog)

**What**: Substituir todas as chamadas `log.Printf` no auth-service por `log/slog` estruturado, seguindo o padrão já usado no mobile-bff.
**Where**: `core/app/auth-service/internal/` (todos os arquivos com `log.Printf`)
**Depends on**: None (independente)
**Reuses**: padrão `logger.New()` do `mobile-bff/internal/infrastructure/logger/`
**Requirements**: PR-18

**Tools**:
- MCP: none
- Skill: `go-style-combined`

**Done when**:
- [ ] Logger slog configurado no `main.go` com `LOG_LEVEL` e `LOG_FORMAT` (json/text)
- [ ] Todas as chamadas `log.Printf` substituídas por `slog.InfoContext`, `slog.ErrorContext`, etc.
- [ ] Logs de auditoria auth (`auth_audit`) mantêm campos estruturados
- [ ] `grep -r "log.Printf" --include="*.go" core/app/auth-service/` não retorna resultado
- [ ] Gate check passes: `go build ./...`

**Tests**: none — logging é infraestrutura
**Gate**: build

---

### T18: Adicionar /ready endpoint nos serviços (auth + catalog)

**What**: Separar health probes em `/health` (liveness — sempre 200) e `/ready` (readiness — verifica DB). Aplicar no auth-service e catalog-service.
**Where**:
- `core/app/auth-service/internal/http/handlers.go` (extensão)
- `core/app/pokemon-catalog-service/internal/http/handlers.go` (extensão)
**Depends on**: None (independente)
**Requirements**: PR-19

**Tools**:
- MCP: none
- Skill: `go-api-design`

**Done when**:
- [ ] `/health` retorna `{"status":"ok","service":"..."}` (sempre 200)
- [ ] `/ready` faz `db.Ping()` e retorna 200 se OK, 503 se falhar
- [ ] Endpoints registrados no `main.go` de cada serviço
- [ ] Kong NÃO roteia `/ready` e `/metrics` (apenas `/health` é público)
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — teste do `/ready` com DB mock (Ping sucesso vs falha)
**Gate**: quick

---

### T19: Corrigir ordenação numérica na home page

**What**: Corrigir `sortHomePokemonPage` para comparar `Number` como inteiro, não como string.
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/home_handler.go:176-178`
**Depends on**: T1
**Requirements**: PR-20

**Tools**:
- MCP: none
- Skill: none

**Done when**:
- [ ] Conversão de `Number` para `int` via `strconv.Atoi()` antes da comparação
- [ ] Tratamento de erro se `Number` não for número válido (fallback para comparação string)
- [ ] Teste confirma: Pokémon #25 (Pikachu) aparece antes de #100 (Voltorb) no sort "Menor número"
- [ ] Gate check passes: `go test -v -race ./tests/unit -run TestHomePokemonPage`

**Tests**: unit — table-driven test com pares de Pokémon (25 vs 100, 1 vs 999, "invalid" vs 5)
**Gate**: quick

---

### T20: Corrigir race condition TOCTOU no AddFavorite

**What**: Substituir check-then-insert por `INSERT ... ON CONFLICT DO NOTHING` atômico.
**Where**: `core/bff/mobile-bff/internal/adapters/outbound/postgres/favorite_repository.go:26-41`
**Depends on**: None (independente)
**Requirements**: PR-21

**Tools**:
- MCP: none
- Skill: `golang-database`

**Done when**:
- [ ] `IsFavorite` check removido
- [ ] `INSERT ... ON CONFLICT (user_id, pokemon_id) DO NOTHING` executado diretamente
- [ ] `RowsAffected()` verificado para determinar se inseriu ou já existia
- [ ] Retorno correto: nil se inseriu; `domain.ErrFavoriteAlreadyExists` se já existia
- [ ] Gate check passes: `go test -v -race ./tests/integration`

**Tests**: integration — teste com goroutines concorrentes tentando adicionar mesmo favorito
**Gate**: full

---

### T21: Corrigir erros de DB mascarados como "pokemon não encontrado"

**What**: No `PostgresPokemonRepository.GetByID`, distinguir `pgx.ErrNoRows` (pokémon não existe) de outros erros (timeout, conexão recusada).
**Where**: `core/bff/mobile-bff/internal/adapters/outbound/postgres/pokemon_repository.go:48-53`
**Depends on**: None (independente)
**Requirements**: PR-23

**Tools**:
- MCP: none
- Skill: `golang-error-handling`

**Done when**:
- [ ] `errors.Is(err, pgx.ErrNoRows)` → retorna `domain.ErrPokemonNotFound`
- [ ] Outros erros → wrappados com `fmt.Errorf("erro ao buscar pokemon %s: %w", id, err)`
- [ ] Handler HTTP converte erros de DB para 500 (não 404)
- [ ] Gate check passes: `go test -v -race ./tests/unit`

**Tests**: unit — teste do repository com mock pgx pool: ErrNoRows → 404, timeout → erro propagado
**Gate**: quick

---

### T22: Delegar filtro/ordenação ao catalog-service no loadHomePokemonPage

**What**: Modificar `loadHomePokemonPage` para enviar query params de filtro (tipo, região, busca, ordenação) ao catalog-service, em vez de carregar todos os Pokémon em memória.
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/home_handler.go:97-142`
**Depends on**: T6
**Requirements**: PR-24

**Tools**:
- MCP: none
- Skill: `go-api-design`

**Done when**:
- [ ] `PokemonCatalogClient.ListPokemons()` aceita query params: `type`, `region`, `search`, `sort`, `page`, `pageSize`
- [ ] Catalog-service handler processa os query params no SQL (WHERE, ORDER BY, LIMIT/OFFSET)
- [ ] `loadHomePokemonPage` faz uma única chamada HTTP com filtros, não carrega todos
- [ ] Paginação correta com `HasNext` baseado no resultado do catalog-service
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — teste do handler BFF com mock do PokemonUseCase retornando página filtrada; teste do handler catalog com mock repository
**Gate**: quick

---

### T23: Corrigir loadHomePokemonPage para usar chamada única com filtros (continuação do T22)

**What**: Finalizar a migração do `loadHomePokemonPage` removendo o loop `for` de fetch-all.
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/home_handler.go`
**Depends on**: T22
**Requirements**: PR-24

**Done when**:
- [ ] Loop `for page := 1; hasNext; page++` removido
- [ ] Chamada única ao `ListPokemons()` com página e filtros
- [ ] Ordenação delegada ao catalog-service (parâmetro `sort`)
- [ ] Gate check passes: `go test -v -race ./tests/unit`

**Tests**: unit — atualizar testes existentes de home handler
**Gate**: quick

---

### T24: Adicionar testes unitários no pokemon-catalog-service

**What**: Criar testes unitários para handlers, repository (com mock) e lógica de domínio do catalog-service. Meta: >75% coverage.
**Where**: `core/app/pokemon-catalog-service/internal/http/handlers_test.go` (novo)
**Depends on**: T10
**Requirements**: PR-25

**Tools**:
- MCP: none
- Skill: `golang-testing`

**Done when**:
- [ ] Testes para `ListPokemons` handler (com filtros, paginação, sem resultados)
- [ ] Testes para `GetPokemonByID` (existente, não encontrado)
- [ ] Testes para `GetFavoriteBatch` (novo endpoint do T10)
- [ ] Testes para `PokemonRepository` com mock (GetByID, List, Search, GetByIDs)
- [ ] Cobertura >75%: `go test -v -race -coverprofile=coverage.out ./internal/...`
- [ ] Gate check passes

**Tests**: unit
**Gate**: quick

---

### T25: Adicionar testes de login e signup no auth-service

**What**: Criar testes unitários para os handlers `Signup` e `Login` que atualmente não têm cobertura.
**Where**: `core/app/auth-service/internal/http/handlers_test.go` (extensão)
**Depends on**: T17
**Requirements**: PR-26

**Tools**:
- MCP: none
- Skill: `golang-testing`

**Done when**:
- [ ] Teste `TestSignup`: sucesso (201), email duplicado (409), senha curta (400), email inválido (400)
- [ ] Teste `TestLogin`: sucesso (200 + token), senha incorreta (401), usuário não encontrado (401)
- [ ] Testes table-driven com `t.Run`
- [ ] Gate check passes: `go test -v -race ./internal/...`

**Tests**: unit
**Gate**: quick

---

### T26: Corrigir testes de integração — t.Skipf → t.Fatal

**What**: Testes de integração não devem pular silenciosamente quando PostgreSQL não está disponível. Se `DATABASE_URL` está configurada mas conexão falha, deve falhar com `t.Fatal`.
**Where**: `core/bff/mobile-bff/tests/integration/postgres_repository_test.go:34`
**Depends on**: None (independente)
**Requirements**: PR-27

**Tools**:
- MCP: none
- Skill: `golang-testing`

**Done when**:
- [ ] `t.Skipf` substituído por `t.Fatalf` quando `DATABASE_URL` está configurada mas conexão falha
- [ ] `t.Skipf` mantido apenas quando `DATABASE_URL` não está configurada (ambiente sem DB)
- [ ] Test runner setup explícito: `if os.Getenv("DATABASE_URL") == "" { t.Skip(...) }`
- [ ] Gate check passes: `make integration-test`

**Tests**: integration — ajuste nos próprios testes
**Gate**: full

---

### T27: Rate limiter distribuído com Redis (sliding window)

**What**: Substituir rate limiter in-memory do middleware por implementação Redis com sliding window (sorted sets), seguindo padrão do `go-ratelimiter-sliding-window`.
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/middleware.go`
**Depends on**: None (independente)
**Requirements**: PR-28

**Tools**:
- MCP: `context7` para `go-redis/redis`
- Skill: none

**Done when**:
- [ ] Redis client configurado via `REDIS_URL` (já existe no config)
- [ ] Implementação sliding window: `ZADD` + `ZREMRANGEBYSCORE` + `ZCARD`
- [ ] Fallback in-memory se Redis indisponível (com log de warning)
- [ ] Config: `AUTH_RATE_LIMIT_REQUESTS` e `AUTH_RATE_LIMIT_WINDOW_SECONDS` mantidos
- [ ] Cleanup de entradas expiradas em goroutine separada (cada 60s)
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — mock do Redis client, testes de janela deslizante
**Gate**: quick

---

### T28: Configurar pool sizes explícitos no pgxpool

**What**: Adicionar `MaxConns`, `MinConns`, `MaxConnLifetime`, `HealthCheckPeriod` explícitos na configuração do `pgxpool` em todos os serviços.
**Where**:
- `core/bff/mobile-bff/internal/adapters/outbound/postgres/database.go`
- `core/app/auth-service/internal/repository/user_repository.go`
- `core/app/pokemon-catalog-service/internal/repository/pokemon_repository.go`
**Depends on**: None (independente)
**Requirements**: PR-29

**Tools**:
- MCP: `context7` para `jackc/pgx`
- Skill: `golang-database`

**Done when**:
- [ ] `pgxpool.Config` com: `MaxConns: 20`, `MinConns: 5`, `MaxConnLifetime: 30m`, `HealthCheckPeriod: 1m`
- [ ] Config aplicada em todos os 3 serviços
- [ ] Valores configuráveis via env (com defaults)
- [ ] Gate check passes: `go build ./...`

**Tests**: none — config de infraestrutura
**Gate**: build

---

### T29: Separar response_builder.go (>500 linhas) em múltiplos arquivos

**What**: Dividir `response_builder.go` (580 linhas) em: `response_builder.go` (builders principais), `response_format.go` (format helpers como weight/height), `response_cookie.go` (cookie building).
**Where**: `core/bff/mobile-bff/internal/adapters/inbound/http/`
**Depends on**: T1
**Requirements**: PR-30

**Tools**:
- MCP: none
- Skill: `go-style-combined`

**Done when**:
- [ ] `response_format.go`: `normalizeWeight`, `normalizeHeight`, `normalizeHexColor`, `formatNumber`, `sanitizeDescription`
- [ ] `response_cookie.go`: `buildAuthCookie`, `buildClearCookie`
- [ ] `response_builder.go`: mantém `Build*` functions (Home, PokemonScreen, PokemonDetail, etc.)
- [ ] Nenhum arquivo >500 linhas no `adapters/inbound/http/`
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — testes existentes de DTO continuam passando
**Gate**: build

---

### T30: Remover código morto e limpar stubs

**What**: Remover: `Validate()` methods não usados, `ErrUserNotFound`, `ErrInvalidPagination`, `GetFavorites()` stub vazio no `PokemonCatalogClient`.
**Where**:
- `core/bff/mobile-bff/internal/domain/pokemon.go`
- `core/bff/mobile-bff/internal/domain/errors.go`
- `core/bff/mobile-bff/internal/adapters/outbound/http/pokemon_catalog_client.go`
- `core/bff/mobile-bff/internal/adapters/outbound/postgres/pokemon_repository.go`
**Depends on**: T12
**Requirements**: PR-31

**Tools**:
- MCP: none
- Skill: none

**Done when**:
- [ ] `Validate()` methods removidos se não usados
- [ ] `ErrUserNotFound`, `ErrInvalidPagination` removidos se não referenciados
- [ ] `GetFavorites()` stub no `PokemonCatalogClient` removido
- [ ] `GetFavorites()` no `PostgresPokemonRepository` removido
- [ ] `grep` confirma que símbolos removidos não têm referências
- [ ] Gate check passes: `go build ./...`

**Tests**: unit — verificar que testes existentes não quebram
**Gate**: build

---

## Granularity Check

| Task | Atômico? | Verificação |
|------|---------|-------------|
| T1 | Sim | Um arquivo criado + substituições (um conceito: DRY type colors) |
| T2 | Sim | Move mocks para novo package + atualiza imports |
| T3 | Sim | Um componente: CircuitBreakerClient |
| T4 | Sim | Integração de um componente em 2 clients existentes |
| T5 | Sim | Helper + alteração em handlers existentes |
| T6 | Sim | Remoção de fallback em um arquivo |
| T7 | Sim | 3 Dockerfiles atualizados (mesmo padrão aplicado) |
| T8 | Sim | Remoção de credencial + ajuste de fallback |
| T9 | Sim | Um middleware novo |
| T10 | Sim | Um endpoint novo no catalog-service |
| T11 | Sim | Um client novo no BFF |
| T12 | Sim | Atualização de handlers + service para novo client |
| T13-T15 | Sim | Um conceito (OTel) aplicado a 3 serviços separados |
| T16 | Sim | Métricas aplicadas a 3 serviços (mesmo padrão) |
| T17 | Sim | Substituição de logging em um serviço |
| T18 | Sim | 2 endpoints novos (health/ready) em 2 serviços |
| T19 | Sim | Correção de uma função de sort |
| T20 | Sim | Substituição de query SQL |
| T21 | Sim | Correção de error wrapping em um método |
| T22-T23 | Sim | Refatoração de handler do BFF + endpoint catalog |
| T24 | Sim | Testes em um serviço |
| T25 | Sim | Testes em handlers específicos |
| T26 | Sim | Ajuste de skip → fail |
| T27 | Sim | Substituição de backend de rate limiter |
| T28 | Sim | Config de pool em 3 arquivos |
| T29 | Sim | Split de arquivo monolítico |
| T30 | Sim | Remoção de código não usado |

## Diagram-Definition Cross-Check

| Phase | Diagrama no design.md | Tasks cobrem? |
|-------|----------------------|---------------|
| Foundation | CircuitBreakerClient, TypeColors | T1-T6 ✓ |
| Security | SecureHeadersMiddleware, Dockerfiles | T7-T9 ✓ |
| Favorites | BatchFavoritesHandler, FavoriteServiceClient | T10-T12 ✓ |
| Observability | OtelInstrumentation, PrometheusMetrics, HealthHandlers | T13-T18 ✓ |
| Bug Fixes | N/A (correções pontuais) | T19-T23 ✓ |
| Tests | N/A (infra de qualidade) | T24-T26 ✓ |
| Infra + Quality | Rate limiter Redis, pool config, file split | T27-T30 ✓ |

## Test Co-location Validation

| Task | Layer | Required | Co-located? |
|------|-------|----------|-------------|
| T1 | domain | unit | ✓ (testa TypeColor()) |
| T2 | adapter outbound | unit | ✓ (atualiza imports nos testes) |
| T3 | adapter outbound | unit | ✓ (mock HTTP server) |
| T4 | adapter outbound | unit | ✓ (atualiza testes existentes) |
| T5 | adapter inbound | unit | ✓ (mock use case) |
| T6 | config | none | ✓ (build gate) |
| T7 | infra | none | ✓ (build gate) |
| T8 | adapter outbound | unit | ✓ (testa NewConnection) |
| T9 | middleware | unit | ✓ (testa headers) |
| T10 | adapter inbound | unit | ✓ (mock repository) |
| T11 | adapter outbound | unit | ✓ (mock HTTP server) |
| T12 | service + handler | unit | ✓ (mock FavoriteProvider) |
| T13-T15 | infra | none | ✓ (build gate) |
| T16 | middleware | none | ✓ (verificável via curl) |
| T17 | infra | none | ✓ (build gate) |
| T18 | adapter inbound | unit | ✓ (mock DB ping) |
| T19 | adapter inbound | unit | ✓ (table-driven sort) |
| T20 | adapter outbound | integration | ✓ (goroutines concorrentes) |
| T21 | adapter outbound | unit | ✓ (mock pgx pool) |
| T22-T23 | adapter inbound | unit | ✓ (mock use case) |
| T24 | adapter inbound | unit | ✓ (testes no catalog-service) |
| T25 | adapter inbound | unit | ✓ (testes signup/login) |
| T26 | integration | integration | ✓ (ajuste nos próprios testes) |
| T27 | middleware | unit | ✓ (mock Redis) |
| T28 | infra | none | ✓ (build gate) |
| T29 | adapter inbound | none | ✓ (build gate + testes existentes) |
| T30 | domain + adapter | none | ✓ (build gate + grep) |

Todos ✓ — nenhuma reestruturação necessária.
