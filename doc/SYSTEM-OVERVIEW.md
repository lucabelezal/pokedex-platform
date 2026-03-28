# System Overview

## Purpose

The Pokedex Platform is organized as a small backend ecosystem built around a mobile-oriented BFF, internal services, an API gateway, and shared infrastructure.

## Main Flow

```text
Client
  -> Kong Gateway
    -> mobile-bff
      -> pokemon-catalog-service
      -> auth-service
      -> PostgreSQL
      -> Redis
```

## Repository Areas

### `core/app/`

Contains internal backend services that expose narrower business capabilities.

- `auth-service`: authentication and token lifecycle.
- `pokemon-catalog-service`: canonical Pokemon catalog access.

### `core/bff/`

Contains `mobile-bff`, the Backend for Frontend that shapes responses for client experience and orchestrates multiple dependencies.

### `core/gateway/`

Contains Kong declarative configuration used as the public entry point.

### `core/infra/`

Contains shared infrastructure assets such as PostgreSQL schema, seed generation inputs, generated seed data, and Redis configuration.

## Architectural Style

At repository level, the platform follows a service-oriented composition:

- Gateway as entry point.
- BFF as client-facing orchestrator.
- Internal services for focused capabilities.
- Infrastructure kept outside application code.

Inside the `mobile-bff`, the intended style is hexagonal architecture. That intent is visible in the `domain`, `ports`, `service`, and `adapters` packages, although a few implementation details still create coupling to concrete infrastructure.

## Current Strengths

- Clear top-level separation between BFF, services, gateway, and infrastructure.
- Good use of Docker Compose to represent the runtime topology.
- BFF already uses ports and adapters terminology consistently.
- Tests exist for unit and integration scenarios in the BFF.

## Current Improvement Areas

- The BFF bootstrap still imports test mocks directly in production composition.
- Some HTTP handlers depend on concrete infrastructure clients instead of ports.
- Some business formatting rules are duplicated across layers.
- Service boundaries are documented in practice, but not yet fully formalized as contracts.

