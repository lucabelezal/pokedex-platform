# Hexagonal Architecture

## Purpose

This document explains how hexagonal architecture is currently applied inside `core/bff/mobile-bff`.

## Layer Model

```text
Inbound adapters
  -> inbound ports
    -> application services
      -> outbound ports
        -> outbound adapters
```

## Current Mapping In The Project

### Domain

Path: `core/bff/mobile-bff/internal/domain`

Responsibility:

- Define core business models.
- Hold domain-level errors and validation rules.
- Stay independent from HTTP, PostgreSQL, and external services.

Examples:

- `pokemon.go`
- `errors.go`

### Inbound Ports

Path: `core/bff/mobile-bff/internal/ports`

Responsibility:

- Define what the application can do from the perspective of callers.
- Describe use cases without tying them to HTTP.

Examples:

- `PokemonUseCase`
- `FavoriteUseCase`

### Application Services

Path: `core/bff/mobile-bff/internal/service`

Responsibility:

- Implement use cases.
- Coordinate domain objects and outbound ports.
- Apply application rules such as validation, orchestration, and pagination defaults.

Examples:

- `PokemonService`
- `FavoriteService`

### Outbound Ports

Path: `core/bff/mobile-bff/internal/ports`

Responsibility:

- Define what infrastructure capabilities the application needs.
- Hide persistence and remote communication details behind interfaces.

Examples:

- `PokemonRepository`
- `FavoriteRepository`

### Inbound Adapters

Path: `core/bff/mobile-bff/internal/adapters/http`

Responsibility:

- Receive HTTP requests.
- Parse transport-specific data.
- Call use cases.
- Convert application results into DTOs and HTTP responses.

Examples:

- `handlers.go`
- `middleware.go`
- `response_builder.go`
- `dto/`

### Outbound Adapters

Path: `core/bff/mobile-bff/internal/adapters/repository`

Responsibility:

- Implement ports using PostgreSQL or remote HTTP services.
- Translate infrastructure details into application-facing behavior.

Examples:

- `favorite_repository.go`
- `pokemon_repository.go`
- `pokemon_catalog_service_repository.go`
- `auth_service_client.go`

## What Is Working Well

- The code already separates domain, ports, services, and adapters.
- Use cases are represented as interfaces.
- Repositories are abstractions rather than direct database calls in handlers.
- HTTP DTOs stay inside the HTTP adapter package.

## Where The Hexagon Is Weaker

### Concrete dependency in the entry adapter

`handlers.go` depends on `*repository.AuthServiceClient`, which is a concrete outbound adapter.

Why this matters:

- An inbound adapter should depend on a port, not on a concrete infrastructure implementation.

### Repository access leaking into the handler

The handler also receives `favoriteRepo` directly.

Why this matters:

- It weakens the use case boundary.
- It makes it easier for transport code to coordinate persistence concerns.

### Production code importing test package

`cmd/server/main.go` imports `tests/mocks`.

Why this matters:

- Test code becomes part of runtime composition.
- The dependency graph becomes harder to reason about.

### Business mapping duplicated across layers

Pokemon type color mapping exists in more than one package.

Why this matters:

- Domain behavior can drift between adapter and service layers.

## Recommended Next Refactor

### Step 1

Introduce an `AuthProvider` port and make the HTTP handler depend on it.

### Step 2

Move runtime fallback implementations out of `tests/` into an internal adapter package.

### Step 3

Ensure handlers talk only to use cases, not directly to repositories.

### Step 4

Centralize shared mapping logic used by both service and response layers.

## Final Assessment

The BFF is already close to a real hexagonal architecture. The remaining work is mostly about tightening dependency direction and clarifying what belongs to application logic versus infrastructure composition.
