# Regras do Projeto

[Visão Geral](README.md) | [Guia](guide.md) | [Decisões](decisions.md) | [Melhores Práticas](best-practices.md) | [Regras do Projeto](regras-projeto.md)

**Status:** `[Canônico]` — regras específicas da Pokedex Platform. Todo código deve seguir.

---

## Estrutura do repositório

```
pokedex-platform/
├── core/
│   ├── app/                          ← serviços internos
│   │   ├── auth-service/
│   │   └── pokemon-catalog-service/
│   ├── bff/
│   │   └── mobile-bff/               ← BFF orientado ao cliente
│   ├── gateway/                      ← configuração do Kong API Gateway
│   ├── infra/                        ← infraestrutura compartilhada
│   └── docker-compose.yml
├── doc/
│   ├── architecture/
│   ├── guia-estilo/                  ← este guia
│   ├── BFF.md
│   ├── DECISIONS.md
│   └── SOLID-AND-PATTERNS.md
├── bruno/                            ← coleções de teste de API
├── .github/                          ← CI/CD, agents, skills
└── .specs/                           ← especificações de features
```

## Serviços

| Serviço | Módulo | Responsabilidade |
|---------|--------|-----------------|
| `mobile-bff` | `core/bff/mobile-bff` | Orquestração voltada ao cliente, composição de respostas, autenticação, favoritos |
| `pokemon-catalog-service` | `core/app/pokemon-catalog-service` | Fonte canônica do catálogo de Pokémon |
| `auth-service` | `core/app/auth-service` | Autenticação, ciclo de vida de tokens JWT, refresh, revogação |

---

## Arquitetura hexagonal — mobile-bff

O `mobile-bff` segue o padrão Ports & Adapters (Hexagonal):

```
internal/
├── domain/                       ← entidades, value objects, erros de domínio
│   ├── pokemon.go                ← Pokemon, Type, Region, Evolution, PokemonDetail, PokemonPage
│   ├── auth_session.go           ← AuthSession
│   ├── errors.go                 ← erros sentinela (var ErrXxx = errors.New(...))
│   └── type_colors.go            ← TypeColor(), AllTypeColors()
├── ports/
│   ├── inbound/                  ← interfaces de use cases
│   │   ├── pokemon_usecase.go    ← PokemonUseCase
│   │   ├── favorite_usecase.go   ← FavoriteUseCase
│   │   ├── auth_usecase.go       ← AuthUseCase
│   │   └── token_validator.go    ← TokenValidator
│   └── outbound/                 ← interfaces de repositórios e clientes externos
│       ├── pokemon_repository.go ← PokemonRepository
│       ├── favorite_repository.go← FavoriteRepository
│       └── auth.go               ← AuthProvider
├── service/                      ← implementação dos use cases
│   ├── pokemon_service.go        ← implementa inbound.PokemonUseCase
│   ├── favorite_service.go       ← implementa inbound.FavoriteUseCase
│   └── auth_service.go           ← implementa inbound.AuthUseCase
├── adapters/
│   ├── inbound/
│   │   └── http/                 ← handlers HTTP, DTOs, middleware
│   │       ├── handler.go        ← Handler struct, RegisterRoutes
│   │       ├── auth_handler.go   ← Signup, Login, Refresh, Logout
│   │       ├── favorite_handler.go
│   │       ├── pokemon_handler.go
│   │       ├── home_handler.go
│   │       ├── middleware.go      ← Auth, CORS, SecureHeaders, RateLimit, RequestLogger
│   │       ├── rate_limiter.go
│   │       ├── response_builder.go← ResponseBuilder, RespondJSON, RespondError, RespondDegraded
│   │       └── dto/
│   └── outbound/
│       ├── http/                 ← clients HTTP para serviços internos
│       │   ├── pokemon_catalog_client.go
│       │   ├── favorite_catalog_client.go
│       │   ├── auth_service_client.go
│       │   └── circuit_breaker.go
│       ├── postgres/             ← repositórios PostgreSQL
│       │   ├── database.go
│       │   ├── pokemon_repository.go
│       │   └── favorite_repository.go
│       └── memory/               ← repositórios in-memory (fallback/testes)
│           └── mock_repositories.go
├── config/
│   └── config.go                 ← LoadConfig() via variáveis de ambiente
├── infrastructure/
│   ├── logger/
│   │   └── logger.go             ← Setup() do slog
│   └── observability/
│       ├── tracing.go            ← OpenTelemetry
│       └── metrics.go            ← Prometheus
cmd/server/main.go                ← entry point, wiring de dependências
```

### Regras arquiteturais

1. **Handlers HTTP dependem de `ports/inbound`** (interfaces de use case), nunca de
   clients ou repositórios concretos.

```go
// Bom:
type Handler struct {
    pokemonUseCase  inbound.PokemonUseCase
    favoriteUseCase inbound.FavoriteUseCase
    authUseCase     inbound.AuthUseCase
}
```

```go
// Ruim:
type Handler struct {
    catalogClient *httpclient.PokemonCatalogServiceRepository
}
```

2. **Services (use cases) dependem de `ports/outbound`** (interfaces), implementadas
   pelos adapters. A injeção é feita no `cmd/server/main.go` (wiring manual).

```go
// Bom:
type PokemonService struct {
    repo outbound.PokemonRepository
}

func NewPokemonService(repo outbound.PokemonRepository) *PokemonService {
    return &PokemonService{repo: repo}
}
```

3. **Erros externos são normalizados no adapter ou na camada de serviço**, não no
   handler. O handler apenas mapeia erros de domínio para HTTP status codes.

```go
// Bom — no handler:
if errors.Is(err, domain.ErrPokemonNotFound) {
    RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
    return
}
```

4. **Novas entidades de domínio vivem em `domain/`**, não em `ports/` ou `dto/`.
   Entidades de domínio são os tipos core do negócio. DTOs são estruturas de
   transferência específicas para HTTP e vivem em `adapters/inbound/http/dto/`.

5. **O package `tests/` nunca é importado por código de produção.** É estritamente
   para testes.

6. **`context.Context` é passado explicitamente** em todas as camadas. Nunca é
   armazenado em structs.

```go
// Bom:
func (s *PokemonService) List(ctx context.Context, params SearchParams) (*PokemonPage, error)
```

```go
// Ruim:
type PokemonService struct {
    ctx context.Context // nunca faça isso
}
```

7. **Interface compliance em tempo de compilação**: todo adapter e service deve
   incluir a verificação `var _ Interface = (*Impl)(nil)`.

```go
// Em service/pokemon_service.go:
var _ inbound.PokemonUseCase = (*PokemonService)(nil)

// Em adapters/outbound/http/pokemon_catalog_client.go:
var _ outbound.PokemonRepository = (*PokemonCatalogServiceRepository)(nil)

// Em adapters/outbound/postgres/pokemon_repository.go:
var _ outbound.PokemonRepository = (*PostgresPokemonRepository)(nil)
```

### Nomenclatura de pacotes no projeto

| Pacote | Nome do diretório | Alias de import comum |
|--------|------------------|-----------------------|
| Handlers HTTP | `adapters/inbound/http` | `httpadapter` |
| Clients HTTP | `adapters/outbound/http` | `httpclient` |
| Repositórios PostgreSQL | `adapters/outbound/postgres` | — |
| Repositórios in-memory | `adapters/outbound/memory` | — |
| Ports inbound | `ports/inbound` | `inbound` |
| Ports outbound | `ports/outbound` | `outbound` |
| Services | `service` | — |
| Domain | `domain` | — |
| Logger | `infrastructure/logger` | `applogger` |

### Estratégia de fallback

- **Pokemons**: sempre via `pokemon-catalog-service` HTTP (obrigatório, sem fallback)
- **Favoritos**: PostgreSQL como preferencial, fallback para mock in-memory
- **Rate limiting**: Redis como preferencial, fallback para in-memory
- **Auth**: `IsTokenActive` retorna `true` se `AUTH_SERVICE_URL` vazio (permite tokens locais)

### DTOs e serialização JSON

DTOs (Data Transfer Objects) são definidos em `adapters/inbound/http/dto/` e
são a única representação exposta via API HTTP. As entidades de domínio
(`domain/`) nunca são serializadas diretamente nas respostas.

```go
// Bom — handler retorna DTO, não entidade de domínio:
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    pokemons, err := h.pokemonUseCase.List(r.Context(), params)
    // converte []domain.Pokemon → []dto.PokemonDTO
    items := make([]dto.PokemonDTO, 0, len(pokemons))
    for _, p := range pokemons {
        items = append(items, dto.ToPokemonDTO(p))
    }
    RespondJSON(w, http.StatusOK, items)
}
```

---

## Commits

Formato: **Conventional Commits** em português do Brasil.

```
tipo(escopo-opcional): descrição curta em português
```

### Tipos

| Tipo | Uso |
|------|-----|
| `feat` | Nova funcionalidade |
| `fix` | Correção de bug |
| `docs` | Documentação |
| `style` | Formatação, ponto e vírgula, etc. (não altera lógica) |
| `refactor` | Refatoração de código (sem feat nem fix) |
| `test` | Adição ou alteração de testes |
| `chore` | Manutenção, tarefas repetitivas |
| `ci` | CI/CD |
| `build` | Build, dependências |
| `perf` | Melhoria de desempenho |
| `revert` | Reversão de commit |

### Exemplos

```
feat(bff): adicionar circuit breaker nos clients HTTP
fix(auth): corrigir validação de token expirado
docs(guia-estilo): adicionar regras de nomenclatura
test(service): adicionar testes para listagem com filtros
refactor(domain): extrair erros para arquivo dedicado
```

### Regras

- Um commit atômico por task concluída
- Mensagens em português do Brasil
- Sempre começar com letra minúscula após o `:`

---

## CI/CD

### Workflows

| Workflow | Descrição |
|----------|-----------|
| `go-ci.yml` | Build, test, vet, lint |
| `conventional-commits.yml` | Validação de mensagens de commit |
| `pr-title-conventional.yml` | Validação de títulos de PR |

### Ferramentas

- **Build**: `go build ./...`
- **Test**: `go test ./... -cover`
- **Vet**: `go vet ./...`
- **Lint**: `golangci-lint run`

---

## Configuração

### Variáveis de ambiente

Toda configuração é feita via variáveis de ambiente, carregada por
`config.LoadConfig()`. Nenhum valor hardcoded.

```go
// Bom:
type Config struct {
    Port                 string        `env:"MOBILE_BFF_PORT" envDefault:"8080"`
    DatabaseURL          string        `env:"DATABASE_URL"`
    PokemonCatalogSvcURL string        `env:"POKEMON_CATALOG_SERVICE_URL"`
    AuthServiceURL       string        `env:"AUTH_SERVICE_URL"`
    RedisURL             string        `env:"REDIS_URL"`
    JWTSecret            string        `env:"JWT_SECRET"`
    ReadTimeout          time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
    WriteTimeout         time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
}
```

### Constantes de tempo

Use constantes nomeadas para timeouts e durations, nunca valores mágicos:

```go
// Bom:
const timeout = 30 * time.Second
const maxRetries = 3
```

```go
// Ruim:
time.Sleep(30) // 30 o quê? segundos? milissegundos?
```

---

## Segurança

### Cookies de autenticação

- `HttpOnly: true` — impede acesso via JavaScript
- `Secure: true` — apenas via HTTPS (baseado em `X-Forwarded-Proto` ou TLS)
- `SameSite: Lax` — proteção contra CSRF

### Rate limiting

Rotas de autenticação (`/auth/login`, `/auth/signup`, `/auth/refresh`, `/auth/logout`)
devem ter rate limiting por IP (padrão: 20 req/min).

### JWT

- Apenas algoritmos HMAC (`HS256`, `HS384`, `HS512`)
- Secret injetado via variável de ambiente `JWT_SECRET`
- Tokens de acesso e refresh com tempos de expiração distintos

### Payload HTTP

- `MaxBytesReader` limitando payload de auth a 8KB
- `json.Decoder` com `DisallowUnknownFields()` em handlers de auth

### Headers de segurança

| Header | Valor |
|--------|-------|
| `Strict-Transport-Security` | `max-age=31536000; includeSubDomains` |
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |

---

## Observabilidade

### Logging

- Logger estruturado via `slog` (Go 1.21+)
- Níveis configuráveis via `LOG_LEVEL` (debug, info, warn, error)
- Formato configurável via `LOG_FORMAT` (json, text)
- Use `slog.InfoContext(ctx, ...)` para logs com trace context

### Tracing

- OpenTelemetry com exportador OTLP
- Configurado via `OTEL_EXPORTER_OTLP_ENDPOINT`
- Middleware extrai e propaga trace context via `otel.GetTextMapPropagator()`

### Métricas

- Prometheus metrics expostas no endpoint `/metrics`
- Métricas padrão: `http_requests_total`, `http_request_duration_seconds`
- Métricas de circuit breaker: estado e contagem de transições

---

## Desenvolvimento spec-driven

Features seguem o fluxo **tlc-spec-driven**:

```
Specify → Design → Tasks → Execute
```

Artefatos em `.specs/features/[feature]/`:

| Arquivo | Descrição |
|---------|-----------|
| `spec.md` | Requisitos com IDs rastreáveis |
| `design.md` | Arquitetura e componentes |
| `tasks.md` | Tarefas atômicas com verificação |
| `validation.md` | Relatório do Verifier |

Estado do projeto registrado em `.specs/STATE.md` (Decisions + Handoff).
