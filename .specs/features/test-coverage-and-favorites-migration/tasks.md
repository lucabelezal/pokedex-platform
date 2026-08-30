# Test Coverage & Favorites Migration - Tasks

## Execution Protocol (MANDATORY -- do not skip)

Implement these tasks with the `tlc-spec-driven` skill: **activate it by name and follow its Execute flow and Critical Rules.** Do not search for skill files by filesystem path. The skill is the source of truth for the full flow (per-task cycle, sub-agent delegation, adequacy review, Verifier, discrimination sensor).

**If the skill cannot be activated, STOP and tell the user — do not proceed without it.**

---

**Design**: `.specs/features/test-coverage-and-favorites-migration/design.md`
**Status**: Approved

---

## Test Coverage Matrix

> Generated from codebase sampling: `core/app/pokemon-catalog-service/go.mod`, `core/app/auth-service/go.mod`, CI workflow `.github/workflows/go-ci.yml`. Guidelines: `AGENTS.md` (cobertura 75% min, 90% ideal, table-driven), `Makefile` (go test -v -race).

| Code Layer | Required Test Type | Coverage Expectation | Location Pattern | Run Command |
|------------|-------------------|---------------------|-----------------|-------------|
| Domain / service (use case) | unit | Todas as branches; 1:1 com ACs da spec | `internal/**/*_test.go` | `go test -v -race ./internal/...` |
| HTTP Handler / adapter inbound | unit | Happy path + edge cases + error paths | `internal/**/*_test.go` | `go test -v -race ./internal/...` |
| Repository / adapter outbound (Postgres) | unit (mock pgxpool) | Key query paths + error handling | `internal/**/*_test.go` | `go test -v -race ./internal/...` |
| Integration (real DB) | integration | Add/Remove flow end-to-end | `tests/integration/*_test.go` | `DATABASE_URL=... go test -v -race ./tests/integration` |
| Config / DTO / entity | none | — (build gate only) | — | `go build ./...` + `go vet ./...` |

## Gate Check Commands

> Generated from project `Makefile` e `go-ci.yml`.

| Gate Level | When to Use | Command |
|------------|-------------|---------|
| Quick | Após tasks com apenas unit tests | `go test -v -race ./internal/...` (catalog/auth) ou `go test -v -race ./tests/unit` (mobile-bff) |
| Full | Após tasks com integration tests | `DATABASE_URL=... go test -v -race ./internal/... -timeout 60s` |
| Build | Após conclusão de phase ou tasks de config/entity | `go build ./... && go vet ./... && go fmt ./...` |
| Coverage | Verificação de cobertura | `go test -v -race -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out \| grep total` |

---

## Execution Plan

Phases são ordenadas e executadas sequencialmente — tasks dentro de cada phase em ordem.

```
Phase 1: Catalog Tests (T1 - T4)              ~4 tasks
Phase 2: Auth Tests   (T5 - T7)               ~3 tasks
Phase 3: Favorites Write (T8 - T15)            ~8 tasks
```

**Batching:** 15 tasks → 3 batches paralelos (TC, TA, FW). Cada batch pode rodar em sub-agent dedicado. Verifier roda após Batch 3.

---

## Task Breakdown

### T1: pokemon_repository_test.go com mock pgxpool ✅ done (0a00688)

**What**: Criar `core/app/pokemon-catalog-service/internal/repository/pokemon_repository_test.go` com mock do `pgxpool` para `PokemonRepository` (GetByID, List, Search, GetByIDs).

**Where**: `core/app/pokemon-catalog-service/internal/repository/pokemon_repository_test.go`

**Depends on**: None

**Requirements**: T1

**Tools**:
- Skill: `golang-testing`

**Done when**:
- [ ] Mock `PokemonRepository` via `mock.Mock` (testify) ou `pashagolub/pgxmock`
- [ ] Testes table-driven: GetByID (found/not found), List (paginado), Search (query válida/vazia), GetByIDs (batch)
- [ ] `go build ./...` passa

**Tests**: unit — `go test -v -race ./internal/repository/...`
**Gate**: quick

---

### T2: handlers_test.go table-driven ✅ done (5d3e88f)

**What**: Criar `core/app/pokemon-catalog-service/internal/http/handlers_test.go` com testes table-driven para handlers.

**Where**: `core/app/pokemon-catalog-service/internal/http/handlers_test.go`

**Depends on**: T1

**Requirements**: T2

**Tools**:
- Skill: `golang-testing`

**Done when**:
- [ ] `TestListPokemons`: happy path, filtros, paginação, empty
- [ ] `TestGetPokemonByID`: found (200), not found (404)
- [ ] `TestGetFavoritesBatch`: valid IDs (200), partial invalid, empty list (200 + `[]`), >100 IDs (400)
- [ ] `TestHealthHandler`: 200 OK com `{"status":"ok"}`
- [ ] Table-driven com `t.Run`, `assert.New(t)` dentro do subtest (não no parent)
- [ ] `go test -v -race ./internal/http/...` passa

**Tests**: unit
**Gate**: quick

---

### T3: Adicionar go test -coverprofile ✅ done (5d3e88f)

**What**: Adicionar target `coverage` no `core/app/pokemon-catalog-service/Makefile`.

**Where**: `core/app/pokemon-catalog-service/Makefile`

**Depends on**: T2

**Requirements**: T3

**Done when**:
- [ ] Target `coverage:` executa `go test -v -race -coverprofile=coverage.out ./internal/...`
- [ ] `go tool cover -func=coverage.out | grep total` mostra percentual
- [ ] `make coverage` funciona

**Tests**: none
**Gate**: build

---

### T4: Verificar cobertura catalog-service >75% ✅ done (5d3e88f — 75.6%)

**What**: Rodar coverage e confirmar >75% (ideal 90%+).

**Where**: `core/app/pokemon-catalog-service/`

**Depends on**: T3

**Requirements**: T4

**Done when**:
- [ ] `go test -v -race -coverprofile=coverage.out ./internal/...` passa
- [ ] `go tool cover -func=coverage.out | grep total` ≥ 75%
- [ ] Ideal: ≥ 90%

**Tests**: full — coverage
**Gate**: full

---

### T5: handlers_test.go Signup ✅ done (748219d)

**What**: Adicionar testes table-driven para `Signup` handler no `core/app/auth-service/internal/http/handlers_test.go`.

**Where**: `core/app/auth-service/internal/http/handlers_test.go`

**Depends on**: None

**Requirements**: T5

**Tools**:
- Skill: `golang-testing`

**Done when**:
- [ ] Casos: success (201 + body com token), email duplicado (409), senha curta (400), email inválido (400)
- [ ] Table-driven com `t.Run`, `assert.New(t)` dentro do subtest
- [ ] Mock `TokenRepository` (testify/mock) ou stub interface
- [ ] `go test -v -race ./internal/http/...` passa

**Tests**: unit
**Gate**: quick

---

### T6: handlers_test.go Login ✅ done (748219d)

**What**: Adicionar testes table-driven para `Login` handler.

**Where**: `core/app/auth-service/internal/http/handlers_test.go`

**Depends on**: T5

**Requirements**: T6

**Tools**:
- Skill: `golang-testing`

**Done when**:
- [ ] Casos: success (200 + token + refresh), senha incorreta (401), usuário não existe (401)
- [ ] Table-driven com `t.Run`
- [ ] Mock `TokenRepository`
- [ ] `go test -v -race ./internal/http/...` passa

**Tests**: unit
**Gate**: quick

---

### T7: Verificar cobertura auth-service >75% ✅ done (748219d — 79.3%)

**What**: Rodar coverage no auth-service.

**Where**: `core/app/auth-service/`

**Depends on**: T6

**Requirements**: T7

**Done when**:
- [ ] `go test -v -race -coverprofile=coverage.out ./internal/...` passa
- [ ] `go tool cover -func=coverage.out | grep total` ≥ 75%
- [ ] Ideal: ≥ 90%

**Tests**: full — coverage
**Gate**: full

---

### T8: Interface FavoriteWriteProvider em ports/outbound/ ✅ done (40e7957)

**What**: Criar `core/bff/mobile-bff/internal/ports/outbound/favorite_write_provider.go` com interface `FavoriteWriteProvider` (AddFavorite, RemoveFavorite).

**Where**: `core/bff/mobile-bff/internal/ports/outbound/favorite_write_provider.go`

**Depends on**: None

**Requirements**: T8

**Tools**:
- Skill: `golang-database`, `golang-error-handling`

**Done when**:
- [ ] Interface `FavoriteWriteProvider` com `AddFavorite(ctx, userID, pokemonID) error` e `RemoveFavorite(ctx, userID, pokemonID) error`
- [ ] `var _ outbound.FavoriteWriteProvider = (*FavoriteWriteClient)(nil)` compile check (quando cliente existir)
- [ ] `go build ./...` passa

**Tests**: none — interface pura
**Gate**: build

---

### T9: PostgresFavoriteRepository.Add/Remove no catalog-service ✅ done (0a00688/748219d)

**What**: Criar `core/app/pokemon-catalog-service/internal/repository/favorite_write_repository.go` com `AddFavorite` (INSERT ON CONFLICT) e `RemoveFavorite`.

**Where**: `core/app/pokemon-catalog-service/internal/repository/favorite_write_repository.go`

**Depends on**: T8

**Requirements**: T9

**Tools**:
- Skill: `golang-database`, `golang-error-handling`

**Done when**:
- [ ] `AddFavorite`: `INSERT INTO user_favorites (user_id, pokemon_id, created_at) VALUES ($1, $2, $3) ON CONFLICT (user_id, pokemon_id) DO NOTHING` + `RowsAffected() == 0` → `ErrFavoriteAlreadyExists`
- [ ] `RemoveFavorite`: `DELETE FROM user_favorites WHERE user_id = $1 AND pokemon_id = $2` + `RowsAffected() == 0` → `ErrFavoriteNotFound`
- [ ] Queries parametrizadas, `ctx` propagado, `rows.Err()` verificado
- [ ] `go build ./...` passa

**Tests**: unit — mock pgxpool, table-driven
**Gate**: quick

---

### T10: FavoriteWriteHandler (POST/DELETE) ✅ done (ja existia em handlers.go)

**What**: Criar `core/app/pokemon-catalog-service/internal/http/favorite_write_handler.go` com handlers POST e DELETE.

**Where**: `core/app/pokemon-catalog-service/internal/http/favorite_write_handler.go`

**Depends on**: T9

**Requirements**: T10

**Tools**:
- Skill: `golang-error-handling`

**Done when**:
- [ ] `POST /v1/pokemons/{id}/favorite` — valida ID, chama `AddFavorite`, retorna 201/409/400
- [ ] `DELETE /v1/pokemons/{id}/favorite` — valida ID, chama `RemoveFavorite`, retorna 204/404
- [ ] Extrai `userID` do contexto JWT (middleware existente)
- [ ] Valida que Pokémon existe antes de add (opcional: `GetByID` check)
- [ ] `go build ./...` passa

**Tests**: unit — mock repository, table-driven para cada status code
**Gate**: quick

---

### T11: Registrar rotas no main.go catalog-service ✅ done (ja existia em NewMux)

**What**: Adicionar rotas de favorite write no `core/app/pokemon-catalog-service/cmd/server/main.go`.

**Where**: `core/app/pokemon-catalog-service/cmd/server/main.go`

**Depends on**: T10

**Requirements**: T11

**Done when**:
- [ ] Rotas registradas com middleware de auth (extrai userID do JWT)
- [ ] `go build ./...` passa

**Tests**: none — wiring
**Gate**: build

---

### T12: FavoriteWriteClient no mobile-bff ✅ done (FavoriteCatalogClient + 40e7957)

**What**: Criar `core/bff/mobile-bff/internal/adapters/outbound/http/favorite_write_client.go` implementando `ports/outbound.FavoriteWriteProvider`.

**Where**: `core/bff/mobile-bff/internal/adapters/outbound/http/favorite_write_client.go`

**Depends on**: T11

**Requirements**: T12

**Tools**:
- Skill: `golang-error-handling`

**Done when**:
- [ ] Struct `FavoriteWriteClient` com `AddFavorite` + `RemoveFavorite`
- [ ] `POST /v1/pokemons/{id}/favorite` e `DELETE /v1/pokemons/{id}/favorite` com 5s timeout
- [ ] Usa `POKEMON_CATALOG_SERVICE_URL` existente (sem nova env)
- [ ] Propaga JWT/or userID via header/cookie
- [ ] `var _ outbound.FavoriteWriteProvider = (*FavoriteWriteClient)(nil)` compile check
- [ ] `go build ./...` passa

**Tests**: unit — mock HTTP server (`httptest.NewServer`), table-driven para 201/409/404/500
**Gate**: quick

---

### T13: Migrar FavoriteService + handlers BFF ✅ done (40e7957)

**What**: Atualizar `core/bff/mobile-bff/internal/service/favorite_service.go` e `favorite_handler.go` para usar `FavoriteWriteProvider` (via feature flag).

**Where**: 
- `core/bff/mobile-bff/internal/service/favorite_service.go`
- `core/bff/mobile-bff/internal/adapters/inbound/http/favorite_handler.go`

**Depends on**: T12

**Requirements**: T13

**Done when**:
- [ ] `FavoriteService` recebe `FavoriteWriteProvider` + flag `FAVORITES_VIA_CATALOG` (bool)
- [ ] `AddFavorite`/`RemoveFavorite` delegam para provider se flag on, senão fallback para `FavoriteRepository`
- [ ] Handlers usam novo service (sem mudança na rota registrada)
- [ ] `go build ./...` passa

**Tests**: unit — mock provider, table-driven para ambos caminhos (flag on/off)
**Gate**: quick

---

### T14: Remover PostgresFavoriteRepository do BFF ✅ done (ja usava FavoriteCatalogClient p/ escrita)

**What**: Remover wiring de `PostgresFavoriteRepository` do `cmd/server/main.go` do mobile-bff (quando flag estiver on por padrão).

**Where**: `core/bff/mobile-bff/cmd/server/main.go`

**Depends on**: T13

**Requirements**: T14

**Done when**:
- [ ] `main.go` não importa mais `adapters/outbound/postgres` para favoritos quando `FAVORITES_VIA_CATALOG=true`
- [ ] Fallback Postgres mantido apenas quando flag `false`
- [ ] `go build ./...` passa com flag on e off

**Tests**: build — verifica ambos caminhos compilam
**Gate**: build

---

### T15: Testes integração add/remove via catalog-service ✅ done (40e7957)

**What**: Testes de integração end-to-end: BFF → catalog-service → DB.

**Where**: `core/bff/mobile-bff/tests/integration/favorite_write_test.go` (novo)

**Depends on**: T14

**Requirements**: T15

**Tools**:
- Skill: `golang-testing`, `golang-database`

**Done when**:
- [ ] `TestAddFavoriteViaCatalog`: add → verify via GetUserFavorites
- [ ] `TestRemoveFavoriteViaCatalog`: add → remove → verify removed
- [ ] `TestConcurrentAddFavorite`: goroutines concorrentes, apenas 1 inserido (`INSERT ON CONFLICT`)
- [ ] `go test -tags=integration -v -race ./tests/integration` passa (com `DATABASE_URL`)

**Tests**: full — integration (requer Postgres em `core/docker-compose.test.yml`, porta 5433)
**Gate**: full

---

## Granularity Check

| Task | Atômico? | Verificação |
|------|---------|-------------|
| TC-01 | Sim | Um arquivo test, mock pgxpool |
| TC-02 | Sim | Um arquivo test, 4 handlers |
| TC-03 | Sim | Makefile target |
| TC-04 | Sim | Coverage check |
| TA-01 | Sim | Signup table-driven |
| TA-02 | Sim | Login table-driven |
| TA-03 | Sim | Coverage check |
| FW-01 | Sim | Uma interface |
| FW-02 | Sim | Um repository file |
| FW-03 | Sim | Um handler file |
| FW-04 | Sim | main.go routes |
| FW-05 | Sim | Um client file |
| FW-06 | Sim | Service + handler migration |
| FW-07 | Sim | main.go wiring |
| FW-08 | Sim | Integration test file |

## Dependencies & Batching

```
T1 → T2 → T3 → T4
T5 → T6 → T7
T8 → T9 → T10 → T11 → T12 → T13 → T14 → T15

Batches paralelos:
  Batch 1 (Sub-agent 1): T1..T4 (catalog)
  Batch 2 (Sub-agent 2): T5..T7 (auth)
  Batch 3 (Sub-agent 3): T8..T15 (favorites, sequencial)
```

## Gate Commands

| Level | Command |
|-------|---------|
| Quick | `go test -v -race ./internal/...` (catalog/auth) ou `go test -v -race ./tests/unit` (mobile-bff) |
| Full | `go test -v -race ./internal/... -timeout 60s` com DATABASE_URL |
| Build | `go build ./... && go vet ./... && go fmt ./...` |
| Coverage | `go test -v -race -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out \| grep total` |

## Verifier

Após todas tasks:
```bash
python3 .claude/skills/tlc-spec-driven/scripts/validate_state.py test-coverage-and-favorites-migration
```
Deve dar **PASS** com evidências de coverage >75% e 3 gaps resolvidos.
