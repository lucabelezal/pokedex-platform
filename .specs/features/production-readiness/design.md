# Production Readiness - Design

**Spec**: `.specs/features/production-readiness/spec.md`
**Status**: Draft

---

## Architecture Overview

A feature abrange 3 serviços e adiciona camadas transversais de resiliência, observabilidade e segurança. As mudanças são orquestradas em 3 níveis:

**Nível 1 — BFF (mobile-bff):** Circuit breaker + retry nos outbound HTTP clients, security headers middleware, OTel + Prometheus instrumentation, novo adapter de favoritos via REST (não direto ao DB), centralização de type colors.

**Nível 2 — Serviços internos (auth, catalog):** OTel + Prometheus, unificação de logging para slog, novo endpoint de batch favorites no catalog-service, health probes separados.

**Nível 3 — Infraestrutura:** Dockerfiles seguros, pgxpool config, rate limiter Redis.

```mermaid
graph TD
    subgraph "Mobile BFF"
        A[HTTP Handler] --> B[Circuit Breaker]
        B --> C[Retry Wrapper]
        C --> D[OTel Instrumented HTTP Client]
        A --> E[Security Headers MW]
        A --> F[Prometheus MW]
        A --> G[Rate Limiter MW]
    end

    subgraph "Pokemon Catalog Service"
        H[Batch Favorites Handler] --> I[Pokemon Repository]
        H --> J[OTel MW]
        I --> K[(PostgreSQL)]
        L[/metrics endpoint]
        M[/health + /ready endpoints]
    end

    subgraph "Auth Service"
        N[slog Logger] --> O[Structured Logging]
        P[OTel MW]
        Q[/metrics endpoint]
        R[/health + /ready endpoints]
    end

    D -->|W3C TraceContext| H
    D -->|W3C TraceContext| N
    G -->|Sliding Window| S[(Redis)]
```

---

## Code Reuse Analysis

### Existing Components to Leverage

| Component | Location | How to Use |
|-----------|----------|------------|
| `http.Client` com timeout | `mobile-bff/internal/adapters/outbound/http/*.go` | Wrappar com circuit breaker + retry, mantendo a mesma interface `PokemonCatalogRepository` e `AuthProvider` |
| Middleware pattern | `mobile-bff/internal/adapters/inbound/http/middleware.go` | Estender com security headers (HSTS, X-Content-Type-Options, X-Frame-Options) |
| `config.LoadConfig()` | `mobile-bff/internal/config/config.go` | Adicionar configs para circuit breaker (threshold, timeout), OTel endpoint, Jaeger URL |
| `log/slog` setup | `mobile-bff/internal/infrastructure/logger/logger.go` | Copiar padrão para auth-service e catalog-service |
| Kong declarative config | `core/gateway/kong/kong.yml` | Adicionar header de tracing (W3C) e rota para /metrics (interna) |
| PostgreSQL schema | `core/infra/postgres/schema/schema.sql` | Adicionar índice em `user_favorites(user_id)` se necessário |

### Integration Points

| System | Integration Method |
|--------|-------------------|
| Pokemon Catalog Service | HTTP REST — novo endpoint `GET /v1/pokemons/favorites?ids=...` |
| Auth Service | HTTP REST — existente `/v1/auth/introspect`, adicionar `/metrics`, `/ready` |
| PostgreSQL | pgxpool — adicionar config explícita de pool (MaxConns, MinConns, etc.) |
| Redis | go-redis/v9 — sliding window rate limiter (substitui in-memory map) |
| Jaeger | OTLP exporter via OpenTelemetry SDK |

---

## Components

### CircuitBreakerClient (mobile-bff)

- **Purpose**: Wrapper que adiciona circuit breaker + retry com backoff ao `http.Client` padrão
- **Location**: `core/bff/mobile-bff/internal/adapters/outbound/http/circuit_breaker.go`
- **Interfaces**:
  - `NewCircuitBreakerClient(name string, cfg CircuitBreakerConfig) *CircuitBreakerClient`
  - `Do(req *http.Request) (*http.Response, error)` — mesma assinatura do `http.Client.Do()`
- **Dependencies**: `sony/gobreaker`, `net/http`
- **Reuses**: padrão de timeout existente nos clients (`pokemon_catalog_client.go:5s`, `auth_service_client.go:10s`)

**Configuração:**

```go
type CircuitBreakerConfig struct {
    MaxRequests    uint32        // half-open max requests (default: 1)
    Interval       time.Duration // cyclic period for closed state (default: 60s)
    Timeout        time.Duration // open → half-open (default: 30s)
    FailureCount   uint32        // consecutive failures to trip (default: 5)
    RetryMax       int           // max retries (default: 3)
    RetryBackoff   []time.Duration // backoff sequence (default: [1s, 3s, 10s])
}
```

### SecureHeadersMiddleware (mobile-bff)

- **Purpose**: Middleware que adiciona security headers em todas as respostas
- **Location**: `core/bff/mobile-bff/internal/adapters/inbound/http/middleware.go` (extensão)
- **Interfaces**: `func SecureHeadersMiddleware(next http.Handler) http.Handler`
- **Headers adicionados**:
  - `Strict-Transport-Security: max-age=31536000; includeSubDomains`
  - `X-Content-Type-Options: nosniff`
  - `X-Frame-Options: DENY`

### TypeColors (domain compartilhado)

- **Purpose**: Source of truth única para mapeamento tipo → cor
- **Location**: `core/bff/mobile-bff/internal/domain/type_colors.go`
- **Interfaces**:
  - `func TypeColor(typeName string) string` — retorna cor hex (sem `#`) para um tipo
  - `func AllTypeColors() map[string]string` — mapa completo (útil para testes)
- **Dependencies**: nenhuma (domínio puro)

### BatchFavoritesHandler (catalog-service)

- **Purpose**: Novo handler HTTP para buscar detalhes de Pokémon em batch por IDs
- **Location**: `core/app/pokemon-catalog-service/internal/http/handlers.go` (extensão)
- **Endpoint**: `GET /v1/pokemons/favorites?ids=1,4,7`
- **Interfaces**:
  - `func (h *Handler) GetFavoritesBatch(w http.ResponseWriter, r *http.Request)`
- **Comportamento**:
  - Aceita query param `ids` com lista separada por vírgula (máx 100)
  - Retorna array de objetos Pokémon (mesmo schema do `GET /v1/pokemons/:id`)
  - IDs não encontrados são omitidos do resultado (não causam erro)
  - Lista vazia retorna `[]` (não null)

### OtelInstrumentation (todos os serviços)

- **Purpose**: Configurar OpenTelemetry SDK com tracing + métricas
- **Location**:
  - `core/bff/mobile-bff/internal/infrastructure/observability/tracing.go`
  - `core/app/auth-service/internal/observability/tracing.go`
  - `core/app/pokemon-catalog-service/internal/observability/tracing.go`
- **Interfaces**:
  - `func InitTracerProvider(endpoint string) (*sdktrace.TracerProvider, error)`
  - `func ShutdownTracerProvider(ctx context.Context) error`
- **Dependencies**: `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`

### PrometheusMetrics (todos os serviços)

- **Purpose**: Handler HTTP para expor métricas Prometheus
- **Location**: cada serviço adiciona `middleware.go` com middleware de métricas
- **Métricas**:
  - `http_requests_total{method, path, status}` — counter
  - `http_request_duration_seconds{method, path}` — histogram
  - `http_errors_total{method, path, error_type}` — counter
- **Endpoint**: `GET /metrics` (interno, não exposto via Kong)

### HealthHandlers (auth + catalog)

- **Purpose**: Separar liveness de readiness
- **Location**: handlers existentes (extensão)
- **Endpoints**:
  - `GET /health` — sempre 200 se processo vivo
  - `GET /ready` — 200 se dependências OK (DB ping), 503 se não

---

## Data Models

### TypeColorMap

```go
// domain/type_colors.go
var typeColors = map[string]string{
    "normal":   "A8A878",
    "fire":     "F08030",
    "water":    "6890F0",
    "electric": "F8D030",
    "grass":    "78C850",
    "ice":      "98D8D8",
    "fighting": "C03028",
    "poison":   "A040A0",
    "ground":   "E0C068",
    "flying":   "A890F0",
    "psychic":  "F85888",
    "bug":      "A8B820",
    "rock":     "B8A038",
    "ghost":    "705898",
    "dragon":   "7038F8",
    "dark":     "705848",
    "steel":    "B8B8D0",
    "fairy":    "EE99AC",
}
```

### CircuitBreakerState (interno ao wrapper)

Estado gerenciado pelo `sony/gobreaker` — não exposto como modelo de domínio.

---

## Error Handling Strategy

| Error Scenario | Handling | User Impact |
|---------------|----------|-------------|
| Circuit breaker aberto (catalog-service down) | Retorna 503 + JSON `{"error": "serviço de catálogo temporariamente indisponível", "degraded": true}` | Usuário vê mensagem de degradação, app pode mostrar estado offline |
| Circuit breaker aberto (auth-service down) | Retorna 503 + JSON `{"error": "serviço de autenticação temporariamente indisponível"}` | Login/refresh bloqueados temporariamente |
| Timeout no catalog-service (retry esgotado) | Circuit breaker registra falha, retorna 503 | Similar ao caso acima |
| Erro de banco no catalog-service | Wrappado como `ErrInternal` com mensagem genérica (não expõe detalhes do DB) | Usuário vê "erro interno", operadores veem detalhes nos logs |
| Jaeger indisponível | Log de warning, tracing descartado (não bloqueia operação) | Sem impacto no usuário |
| Redis indisponível (rate limiter) | Fallback para in-memory limiter com log de warning | Rate limiting por instância (não distribuído) |

---

## Risks & Concerns

| Concern | Location (file:line) | Impact | Mitigation |
|---------|---------------------|--------|------------|
| Produção importa package de testes | `cmd/server/main.go:18` | Binário de produção depende de mocks; quebra compilação se tests/ não existe | Extrair mocks para `adapters/outbound/memory/` (PR-14) |
| Type colors com valores inconsistentes | 5 arquivos (service, handler, postgres repo, mock, catalog repo) | Cor de tipo diferente entre BFF e catalog-service (ex: Normal = `#A8A878` vs `#A8A77A`) | Centralizar em `domain/type_colors.go` (PR-13) |
| TOCTOU race no AddFavorite | `postgres/favorite_repository.go:26-41` | Duas requisições concorrentes podem tentar inserir o mesmo favorito | `INSERT ... ON CONFLICT DO NOTHING` + remover check prévio (PR-21) |
| Ordenação lexicográfica na home | `home_handler.go:176-178` | Pokémon #100 aparece antes do #25 | Converter Number para int antes de comparar (PR-20) |
| N+1 queries no GetUserFavorites | `favorite_handler.go:84-98` | 50 favoritos = 50 chamadas HTTP sequenciais | Batch endpoint no catalog-service (PR-22) |
| Erros de DB mascarados como "não encontrado" | `postgres/pokemon_repository.go:48-53` | Timeout de DB tratado como Pokémon inexistente | Propagar erro real, distinguir `pgx.ErrNoRows` de outros erros (PR-23) |
| Rate limiter in-memory sem cleanup | `middleware.go:350-387` | Mapa cresce indefinidamente sob ataque sustentado | Substituir por Redis sliding window (PR-28) |
| catalog-service fallback com 2 Pokémon | `pokemon-catalog/main.go:24-25` | Serviço em produção serve dados incorretos sem indicar falha | Remover fallback, retornar erro explícito (PR-15) |
| Testes de integração pulam silenciosamente | `postgres_repository_test.go:34` | CI passa sem cobertura real se DB não disponível | `t.Fatal` em vez de `t.Skipf` quando variáveis de ambiente indicam que DB deveria estar disponível (PR-27) |
| Auht-service sem testes de login/signup | `auth-service/internal/http/handlers_test.go` | Regressões em fluxo crítico de autenticação não detectadas | Adicionar table-driven tests para Signup e Login (PR-26) |

---

## Tech Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Biblioteca de circuit breaker | `sony/gobreaker` | Leve (~500 LOC), idiomática Go, zero dependências, callback hooks para logging/métricas |
| Retry strategy | Exponencial fixo (1s, 3s, 10s) com jitter | 3 tentativas + jitter evita thundering herd; backoff exponencial é padrão da indústria |
| Tracing backend | Jaeger via OTLP HTTP | OTLP é padrão aberto da OTel; Jaeger é leve para dev; pode trocar por outro backend sem mudar código |
| Rate limiter backend | Redis sorted sets (sliding window) | Algoritmo preciso, testado no `go-ratelimiter-sliding-window` do msfidelis; Redis já existe na infra |
| BFF → Catalog communication | HTTP REST (não gRPC) | Manter consistência com arquitetura existente; gRPC seria adoção prematura para este escopo |
| Security headers | Middleware dedicado no BFF | Kong poderia adicionar, mas ter no BFF garante headers mesmo em dev local sem Kong |
| Health probes | Dois endpoints: `/health` (liveness) + `/ready` (readiness) | Padrão Kubernetes; `/ready` verifica DB, `/health` apenas processo vivo |

> **Project-level decisions:** Nenhuma decisão nesta feature é project-level (todas são feature-local). Decisões arquiteturais (hexagonal, net/http, PostgreSQL) já são convenções do projeto documentadas em AGENTS.md.
