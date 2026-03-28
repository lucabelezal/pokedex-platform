# Core

This directory contains the runtime code and infrastructure of the Pokedex Platform.

## Structure

```text
core/
├── app/           # Internal backend services
├── bff/           # Backend for Frontend applications
├── gateway/       # API gateway configuration
├── infra/         # Shared infrastructure assets
├── bruno/         # API request collection
├── observability/ # Monitoring and operational assets
└── docker-compose.yml
```

## Purpose

The `core/` folder keeps implementation artifacts together so the repository root can stay focused on cross-cutting materials such as documentation and contribution files.

## Main Areas

- `app/`: service-specific business capabilities such as `auth-service` and `pokemon-catalog-service`.
- `bff/`: client-facing composition layer, currently `mobile-bff`.
- `gateway/`: Kong declarative configuration.
- `infra/`: PostgreSQL schema, seeds, source JSON files, Redis config, and data tooling.

## Local Runtime

From the repository root:

```bash
docker compose -f core/docker-compose.yml up --build
```

## Related Documentation

- `../README.md`
- `../doc/SYSTEM-OVERVIEW.md`
- `../doc/BFF.md`
- `../doc/INFRA.md`

