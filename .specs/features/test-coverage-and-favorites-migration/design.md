# Test Coverage & Favorites Migration — Design

## Context

Gaps identified in `production-readiness/validation.md`:
- **G1 (PR-25)**: pokemon-catalog-service has 0% unit-test coverage
- **G2 (PR-26)**: auth-service login/signup handlers untested (critical auth flow)
- **G3 (PR-11 partial)**: BFF favorite writes still hit PostgreSQL directly via `PostgresFavoriteRepository`, violating separation of concerns

This feature closes the 3 gaps with tests + favorite-write isolation.

## Architecture Decisions

### G1: Catalog-service Test Strategy
- Mock `pgxpool` via the existing `PokemonRepository` interface (`ports/outbound`)
- Table-driven tests for every handler + repository method
- `goleak.VerifyTestMain` in `TestMain` for goroutine leak detection
- Coverage via `go test -coverprofile=coverage.out ./internal/...` — target >75% (ideal 90%+)

### G2: Auth-service Test Strategy
- Table-driven tests for Signup/Login handlers
- Mock `TokenRepository` interface (sql/mock or testify/mock)
- `goleak.VerifyTestMain` in `TestMain`
- Coverage via `go test -coverprofile=coverage.out ./internal/...`

### G3: Favorite Write Migration
- **Endpoint design**: RESTful in catalog-service
  - `POST   /v1/pokemons/{id}/favorite` — add favorite
  - `DELETE /v1/pokemons/{id}/favorite` — remove favorite
- **Auth**: JWT via existing `TokenValidator` middleware
- **DB**: use existing `user_favorites` table, `INSERT ON CONFLICT DO NOTHING`
- **BFF**: new `FavoriteWriteClient` implements `ports/outbound/FavoriteWriteProvider`
  - Remove `PostgresFavoriteRepository` from BFF
  - Feature flag `FAVORITES_VIA_CATALOG=true` for instant rollback

## Component Design

### Catalog-service

```
internal/
  http/
    handlers.go              # + favoriteWriteHandler (POST/DELETE)
    handlers_test.go         # table-driven: add/remove favorite, existing handlers
  repository/
    favorite_write_repository.go       # Add/Remove favoritos (new)
    favorite_write_repository_test.go  # mock pgxpool (new)
  domain/
    errors.go                # ErrFavoriteAlreadyExists, ErrFavoriteNotFound (reuse)
```

### Mobile-bff

```
internal/
  adapters/outbound/http/
    favorite_write_client.go        # implements FavoriteWriteProvider (new)
    favorite_write_client_test.go   # mock HTTP server
  service/
    favorite_service.go             # uses FavoriteWriteProvider (updated)
  ports/outbound/
    favorite_write_provider.go      # interface Add/Remove (new)
  config/
    config.go                       # FAVORITES_VIA_CATALOG flag (new env)
  cmd/server/main.go                # wiring: FavoriteWriteClient vs Postgres repo
```

## Data Flow (G3)

```
BFF AddFavorite
  │
  ├─ (flag off) → PostgresFavoriteRepository (legacy, to be removed after FW-07)
  │
  └─ (flag on)  → FavoriteWriteClient
                    │
                    POST /v1/pokemons/{id}/favorite
                    │
                    catalog-service: FavoriteWriteHandler.AddFavorite
                    │
                    FavoriteWriteRepository.AddFavorite (INSERT ON CONFLICT)
                    │
                    user_favorites table
                    │
                    201 Created / 409 Conflict
```

Same flow for `RemoveFavorite` with `DELETE`.

## Security

- Rate limit on write endpoints (same `AUTH_RATE_LIMIT` middleware)
- Validation: Pokémon must exist (`GetByID`) before add — returns 404 if not found
- JWT validation via existing middleware (`TokenValidator`)

## Observability

- OTel spans on new handlers (`FavoriteWriteHandler.AddFavorite`, `RemoveFavorite`)
- Metrics: `favorite_write_duration_seconds`, `favorite_write_errors_total`
- Structured logs `slog` with `pokemon_id`, `user_id`, `op`

## Rollout

- Feature flag `FAVORITES_VIA_CATALOG=true` (default on) in `mobile-bff/internal/config`
- Rollback: set `FAVORITES_VIA_CATALOG=false` → BFF fallback to `PostgresFavoriteRepository`
- No database migration needed (table `user_favorites` already exists)
- Docker compose: no change (catalog-service already owns DB)

## Alternatives Considered

- **gRPC for favorite writes** — rejected; HTTP is sufficient for current platform scope (per spec Out of Scope)
- **Event-driven (Kafka) for favorites** — rejected; sync REST is simpler and favorite writes are not high-throughput
- **Separate favorites microservice** — rejected; catalog-service already owns Pokémon data, keeping favorites co-located reduces latency
