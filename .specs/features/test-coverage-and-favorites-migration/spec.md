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
| Frontend | Frontend está em outro repositório |
| Migração para gRPC | HTTP é suficiente para o escopo atual (per spec Out of Scope) |

---

## Assumptions & Open Questions

| Assumption / decision | Chosen default | Rationale | Confirmed? |
|-----------------------|---------------|-----------|------------|
| Mock pgxpool para catalog-service | `pashagolub/pgxmock` | Mais próximo do driver real que testify/mock genérico | n |
| Handlers de signup/login | `core/app/auth-service/internal/http/handlers.go` | Já existem, apenas adicionar testes | y |
| Endpoint de favoritos write | `POST/DELETE /v1/pokemons/{id}/favorite` | RESTful, consistente com GET batch existente | n |
| Feature flag favoritos | `FAVORITES_VIA_CATALOG` env (default true) | Rollback instantâneo sem deploy | n |
| Tabela de favoritos | `user_favorites` existente | Sem migração de dados | y |

**Open questions:** none — todas as ambiguidades estão resolvidas ou registradas acima.

---

## User Stories

### P1: Testes no Pokemon Catalog Service ⭐ MVP

**User Story**: Como desenvolvedor da plataforma, quero testes unitários no catalog-service com cobertura >75% para que mudanças possam ser feitas com confiança.

**Why P1**: Sem testes, qualquer mudança no catálogo pode causar regressões não detectadas.

**Acceptance Criteria**:

1. WHEN `go test -coverprofile=coverage.out ./internal/...` é executado no catalog-service THEN a cobertura SHALL ser >75%
2. WHEN handlers HTTP são testados THEN cada endpoint SHALL ter teste de happy path + erro
3. WHEN repository é testado THEN queries SHALL ser testadas com mock pgxpool

**Independent Test**: `go test -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out | grep total`

---

### P2: Testes de Login e Signup no Auth Service

**User Story**: Como desenvolvedor da plataforma, quero que os handlers de login e signup do auth-service tenham testes unitários para que o fluxo crítico de autenticação seja confiável.

**Why P2**: Login e signup são fluxos críticos sem cobertura atual.

**Acceptance Criteria**:

1. WHEN o handler Signup é testado THEN SHALL cobrir: sucesso (201), email duplicado (409), senha curta (400)
2. WHEN o handler Login é testado THEN SHALL cobrir: sucesso (200 + token), senha incorreta (401)

**Independent Test**: `go test -v -race ./internal/http/...` no auth-service com `handlers_test.go`

---

### P3: Migração de Escritas de Favoritos para Catalog-Service

**User Story**: Como desenvolvedor da plataforma, quero que as escritas de favoritos (add/remove) sejam feitas via catalog-service em vez de acesso direto ao PostgreSQL para que a arquitetura hexagonal seja respeitada.

**Why P3**: BFF acessando DB diretamente quebra a separação de responsabilidades.

**Acceptance Criteria**:

1. WHEN o BFF adiciona um favorito THEN SHALL chamar `POST /v1/pokemons/{id}/favorite` no catalog-service
2. WHEN o BFF remove um favorito THEN SHALL chamar `DELETE /v1/pokemons/{id}/favorite` no catalog-service
3. WHEN o catalog-service processa um favorito THEN SHALL usar a tabela `user_favorites` existente
4. WHEN o BFF inicializa com `FAVORITES_VIA_CATALOG=true` THEN SHALL não conectar ao PostgreSQL para operações de favoritos

**Independent Test**: Adicionar favorito via BFF e verificar via `GET /v1/pokemons/favorites?ids=...` no catalog-service.

---

## Requirement Traceability

| ID | Story | Status |
|----|-------|--------|
| TC-01 | P1: Testes no catalog-service | Pending |
| TC-02 | P2: Testes de login/signup no auth-service | Pending |
| TC-03 | P3: Migração de escritas de favoritos | Pending |

**Coverage:** 3 total, 0 mapped to tasks, 3 unmapped

---

## Success Criteria

- [ ] `go test -coverprofile=coverage.out ./internal/...` no catalog-service com cobertura >75%
- [ ] `go test -v -race ./internal/http/...` no auth-service com testes de Signup e Login
- [ ] Endpoints `POST/DELETE /v1/pokemons/{id}/favorite` no catalog-service
- [ ] BFF sem acesso direto ao DB para operações de escrita de favoritos
