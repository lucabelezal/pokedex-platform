# Decisions

## Purpose

This file records the main architectural decisions that are already visible in the codebase.

## Decision 1: Use A BFF For Client-Facing API

### Decision

Expose client-oriented endpoints through `mobile-bff` instead of exposing internal services directly.

### Why

- Client responses can be shaped for UX needs.
- Internal services stay narrower and more focused.
- Cross-service orchestration stays out of the client.

### Consequence

The BFF becomes an important composition layer and must avoid growing into a generic monolith.

## Decision 2: Keep The Canonical Pokemon Catalog In `pokemon-catalog-service`

### Decision

Use `pokemon-catalog-service` as the canonical read source for Pokemon catalog information.

### Why

- Catalog rules stay centralized.
- The BFF can stay focused on presentation and orchestration.

### Consequence

The BFF should avoid re-implementing catalog rules beyond presentation-specific formatting.

## Decision 3: Keep Favorites In The BFF Context For Now

### Decision

Favorites are currently handled by the BFF context instead of a dedicated service.

### Why

- Simpler implementation for the current project stage.
- Favorites are tightly connected to authenticated user experience.

### Consequence

This is acceptable now, but it may become a future extraction candidate if favorite logic grows, needs independent scaling, or must be shared by more clients.

## Decision 4: Use Hexagonal Architecture In The BFF

### Decision

Organize the BFF around domain, ports, services, and adapters.

### Why

- Improves separation of concerns.
- Makes testing easier.
- Reduces coupling to transport and persistence details.

### Consequence

The codebase should keep enforcing dependency direction. Concrete adapters must depend on ports, and entry adapters should not bypass the application layer.

## Decision 5: Prefer Deterministic Seed Generation Over Runtime Data Setup

### Decision

Maintain JSON files as source data and generate SQL seeds deterministically.

### Why

- Reproducible local environments.
- Easy review of source data changes.
- Clear pipeline from source content to database initialization.

### Consequence

Any change in catalog content must respect the JSON-to-SQL generation flow.

