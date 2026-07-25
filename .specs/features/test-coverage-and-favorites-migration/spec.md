# Test Coverage & Favorites Migration — Specification

## Problem Statement

A feature production-readiness deixou 3 gaps identificados: o pokemon-catalog-service não tem testes unitários (0% coverage), o auth-service não tem testes para os handlers de login e signup (fluxo crítico de autenticação), e as escritas de favoritos no BFF ainda acessam o PostgreSQL diretamente (`PostgresFavoriteRepository`), violando a separação de responsabilidades.

## Goals

- [ ] Cobertura de testes >75% no pokemon-catalog-service (handlers + repository com mocks)
- [ ] Testes table-driven para Signup (201, 409, 400) e Login (200, 401) no auth-service
- [ ] Endpoint de escrita de favoritos no catalog-service (`POST/DELETE /v1/pokemons/{id}/favorite`)
- [ ] BFF sem acesso direto ao DB para operações de escrita de favoritos

## Out of Scope

| Feature | Reason |
|---------|--------|
| Migração de banco de favoritos | Tabela user_favorites já existe e é usada; endpoint apenas expõe via REST |

## User Stories

### P1: Testes no Pokemon Catalog Service ⭐ MVP

**Acceptance Criteria**:
1. WHEN executamos `go test ./...` no catalog-service THEN cobertura >75%
2. WHEN testamos handlers HTTP THEN cada endpoint tem teste de happy path + erro
3. WHEN testamos repository THEN queries são testadas com mock pgxpool

**Independent Test**: `go test -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out | grep total`

### P2: Testes de Login e Signup no Auth Service

**Acceptance Criteria**:
1. WHEN testamos Signup THEN cobre: sucesso (201), email duplicado (409), senha curta (400)
2. WHEN testamos Login THEN cobre: sucesso (200 + token), senha incorreta (401)

### P3: Migração de Escritas de Favoritos para Catalog-Service

**Acceptance Criteria**:
1. WHEN BFF adiciona um favorito THEN chama `POST /v1/pokemons/{id}/favorite` no catalog-service
2. WHEN BFF remove um favorito THEN chama `DELETE /v1/pokemons/{id}/favorite` no catalog-service
3. WHEN catalog-service processa favorito THEN usa tabela `user_favorites` existente
4. WHEN BFF inicializa THEN não conecta mais ao PostgreSQL para favoritos

## Requirement Traceability

| ID | Story | Status |
|----|-------|--------|
| TC-01 | P1: Testes no catalog-service | Pending |
| TC-02 | P2: Testes de login/signup no auth-service | Pending |
| TC-03 | P3: Migração de escritas de favoritos | Pending |
