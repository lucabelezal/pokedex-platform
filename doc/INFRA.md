# Infrastructure

## Purpose

The `core/infra/` directory contains the shared technical assets required to run the platform locally and in containerized environments.

## PostgreSQL

PostgreSQL is used as the main persistent store.

### Files

- `core/infra/postgres/schema/schema.sql`: database schema.
- `core/infra/postgres/seeds/init-data.sql`: generated seed data.
- `core/infra/postgres/source-json/`: source of truth for catalog content.
- `core/infra/postgres/json2sql/`: tool that converts source JSON files into SQL seed data.

### Current Data Pipeline

```text
source-json/*.json
  -> json2sql
    -> seeds/init-data.sql
      -> PostgreSQL container initialization
```

This is a good fit for deterministic local environments and study projects because it keeps the dataset versioned and reproducible.

## Redis

Redis is provisioned in Docker Compose, but its architectural role is not yet strongly visible in the application code.

This usually means one of two things:

- It is planned for future use.
- It is infrastructure-ready but not yet integrated into the core flows.

That should be documented explicitly as the project evolves.

## Docker Compose

The `core/docker-compose.yml` describes the full local topology:

- PostgreSQL
- Redis
- `pokemon-catalog-service`
- `auth-service`
- `mobile-bff`
- Kong

This is one of the clearest documents of the real runtime architecture today.

## Improvement Opportunities

- Document which service owns which tables.
- Clarify whether the BFF may persist its own data long term or only temporarily.
- Decide if Redis will be used for cache, sessions, rate limiting, or not at all.
- Add environment variable documentation per service.

