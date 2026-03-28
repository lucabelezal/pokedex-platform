# BFF

## Purpose

The `mobile-bff` is the client-facing application of the platform. Its job is not to be the source of truth for the whole domain, but to provide responses shaped for mobile and web experience.

## Current Responsibilities

- Expose frontend-oriented HTTP endpoints.
- Aggregate data from `pokemon-catalog-service`.
- Manage user favorites.
- Delegate authentication flows to `auth-service`.
- Return responses already shaped for UI consumption.

## What The BFF Should Own

- Request orchestration.
- Experience-oriented response composition.
- Session and identity propagation.
- Client-specific contracts.

## What The BFF Should Not Own

- Canonical Pokemon catalog rules that belong to `pokemon-catalog-service`.
- Authentication core logic that belongs to `auth-service`.
- Infrastructure-specific decisions leaking into use case code.

## Hexagonal Architecture Status

The current BFF is partially well aligned with hexagonal architecture:

- `internal/domain` keeps core models independent from transport concerns.
- `internal/ports` defines inbound and outbound contracts.
- `internal/service` acts as an application layer implementing use cases.
- `internal/adapters/http` and `internal/adapters/repository` work as entry and exit adapters.

The main issue is not the folder layout. The main issue is dependency direction in a few places.

## Main Improvement Points

### 1. Remove test mocks from production bootstrap

`cmd/server/main.go` imports `tests/mocks` to assemble runtime fallbacks. This makes production composition depend on a test package.

Recommended direction:

- Move fallback implementations to `internal/adapters/repository/mock` or `internal/adapters/repository/memory`.
- Keep `tests/` only for test-only utilities.

### 2. Replace concrete auth client dependency with a port

`internal/adapters/http/handlers.go` depends on `*repository.AuthServiceClient`.

Recommended direction:

- Create an outbound port such as `AuthProvider`.
- Let the HTTP handler depend on the port, not on the concrete client.
- Keep `AuthServiceClient` as one adapter implementing that port.

### 3. Avoid bypassing the use case layer

The handler currently receives both `favoriteUseCase` and `favoriteRepo`. When an entry adapter talks directly to a repository, the application layer becomes easier to bypass.

Recommended direction:

- Keep handlers calling use cases only.
- If a handler needs extra favorite data, move that behavior into an application service.

### 4. Remove duplicated mapping rules

Type color mapping appears in more than one place. That creates drift risk.

Recommended direction:

- Centralize this mapping in a domain policy or dedicated mapper package owned by one layer.

### 5. Review port design

`PokemonRepository` includes `GetFavorites`, which does not appear to match the main responsibility of the Pokemon catalog port.

Recommended direction:

- Keep ports focused by capability.
- Move favorite-related operations to favorite-specific ports only.

## Practical Assessment

The BFF is not badly structured. It already has the right architectural vocabulary and most of the right boundaries. The main next step is to make the dependency direction stricter so that the hexagonal structure is enforced by code, not only by folder names.

