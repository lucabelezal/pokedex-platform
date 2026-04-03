# Mobile BFF - API da Pokedex

## Visao Geral

O Mobile-BFF e um Backend-for-Frontend (BFF) em Go para atender clientes mobile/web da Pokedex.
O projeto usa arquitetura hexagonal (Ports and Adapters), PostgreSQL, testes unitarios e testes de integracao.
Ele faz parte da pasta `core/`, onde ficam os artefatos executaveis da plataforma.

Importante: o `POKEMON_CATALOG_SERVICE_URL` e obrigatorio para iniciar o BFF. O catalogo nao possui mais fallback direto para repositorio de pokemons em Postgres dentro do BFF.

## Localizacao No Repositorio

```text
core/bff/mobile-bff
```

## Arquitetura

### Padrao Hexagonal (Ports and Adapters)

```text
internal/
├─ domain/                          # Hexagono: entidades e erros de dominio
│  ├─ auth_session.go
│  ├─ errors.go
│  └─ pokemon.go
│
├─ ports/
│  ├─ inbound/                      # Contratos que adaptadores HTTP consomem
│  │  ├─ auth_usecase.go            # AuthUseCase
│  │  ├─ favorite_usecase.go        # FavoriteUseCase
│  │  ├─ pokemon_usecase.go         # PokemonUseCase
│  │  └─ token_validator.go         # TokenValidator (usado pelo middleware)
│  └─ outbound/                     # Contratos que os servicos exigem de recursos externos
│     ├─ auth.go                    # AuthProvider
│     ├─ favorite_repository.go     # FavoriteRepository
│     └─ pokemon_repository.go      # PokemonRepository
│
├─ service/                         # Hexagono: implementacoes dos casos de uso
│  ├─ auth_service.go
│  ├─ favorite_service.go
│  └─ pokemon_service.go
│
├─ adapters/
│  ├─ inbound/
│  │  └─ http/                      # Adaptadores de entrada: handlers HTTP (pkg httphandler)
│  │     ├─ handler.go              # Registro de rotas
│  │     ├─ auth_handler.go
│  │     ├─ favorite_handler.go
│  │     ├─ home_handler.go
│  │     ├─ pokemon_handler.go
│  │     ├─ middleware.go           # Auth, CORS, rate-limit, request logger
│  │     ├─ response_builder.go
│  │     └─ dto/
│  └─ outbound/
│     ├─ http/                      # Clientes HTTP externos
│     │  ├─ auth_service_client.go  # implementa outbound.AuthProvider e inbound.TokenValidator
│     │  └─ pokemon_catalog_client.go # implementa outbound.PokemonRepository
│     └─ postgres/                  # Adaptadores PostgreSQL
│        ├─ database.go             # Pool de conexao (pgx/v5)
│        ├─ favorite_repository.go  # implementa outbound.FavoriteRepository
│        └─ pokemon_repository.go   # implementa outbound.PokemonRepository
│
└─ infrastructure/
   └─ logger/                       # Setup do slog (LOG_LEVEL, LOG_FORMAT)
      └─ logger.go

tests/
├─ unit/                            # Testes unitarios (sem infraestrutura)
├─ integration/                     # Testes com banco PostgreSQL real
└─ mocks/                           # Repositorios fake para testes
```

## Endpoints

### Publicos

- `GET /health`
- `GET /api/v1/pokemons`
- `GET /api/v1/pokemons/search`
- `GET /api/v1/pokemons/{id}/details`
- `GET /api/v1/home`

Para filtro por tipo, use `GET /api/v1/pokemons?type=Electric&page=0&size=20`.

### Autenticados

- `GET /api/v1/me`
- `POST /api/v1/pokemons/{id}/favorite`
- `DELETE /api/v1/pokemons/{id}/favorite`
- `GET /api/v1/me/favorites`

### Sessao/Auth

- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

Os endpoints autenticados aceitam token por:
- Header: `Authorization: Bearer <jwt>`
- Cookie HTTP-only: `auth_token`

Regra de precedencia: quando os dois forem enviados, o header `Authorization` tem prioridade.

## Contrato Entre BFF E Catalog Service

O BFF atende o contrato rico do front, enquanto o `pokemon-catalog-service` concentra o catalogo canonico de Pokemon.

### BFF (mobile-bff)

- Compoe respostas para UI (`home`, cards e detalhe rico)
- Aplica regras de experiencia para cliente mobile/web
- Mantem regra de favoritos por usuario (escrita)

### Catalog Service (pokemon-catalog-service)

- Exposicao de dados canonicos de catalogo (leitura)
- Endpoints consumidos pelo BFF:
   - `GET /v1/pokemons`
   - `GET /v1/pokemons/search`
   - `GET /v1/pokemons/type/{type}`
   - `GET /v1/pokemons/{id}`
   - `GET /v1/pokemon-details/{id}`

### Decisao Atual

- Catalogo: no `pokemon-catalog-service`
- Favoritos: no BFF por usuario autenticado
- Evolucao futura: extrair favoritos para um servico proprio quando login/cadastro estiverem integrados

## Stack

- Go 1.24
- net/http
- PostgreSQL + pgx/v5
- Testify
- Docker / Docker Compose

## Como Executar

### Variaveis de Ambiente

| Variavel | Obrigatoria | Descricao |
|---|---|---|
| `POKEMON_CATALOG_SERVICE_URL` | Sim | URL base do `pokemon-catalog-service` |
| `JWT_SECRET` | Sim | Chave para validacao de tokens JWT |
| `AUTH_SERVICE_URL` | Nao | URL do `auth-service` (funcionalidades de auth) |
| `DATABASE_URL` | Nao | PostgreSQL para favoritos (usa mock se ausente) |
| `MOBILE_BFF_PORT` | Nao | Porta HTTP (padrao: 8080) |
| `LOG_LEVEL` | Nao | Nivel de log: `debug`, `info`, `warn`, `error` (padrao: `info`) |
| `LOG_FORMAT` | Nao | Formato de log: `json` (padrao) ou `text` (legivel no terminal) |

### Modo local (com catalog service)

```bash
export MOBILE_BFF_PORT=8080
export POKEMON_CATALOG_SERVICE_URL="http://localhost:8081"
export LOG_FORMAT=text   # saida legivel no terminal
go run ./cmd/server/main.go
```

### Modo PostgreSQL

```bash
docker compose -f docker-compose.test.yml up -d
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pokedex"
export MOBILE_BFF_PORT=8080
export POKEMON_CATALOG_SERVICE_URL="http://localhost:8081"
go run ./cmd/server/main.go
```

No modo PostgreSQL do BFF, o banco e usado para favoritos e dados de suporte locais. O catalogo continua vindo do `pokemon-catalog-service`.

### Com Makefile

```bash
make unit-test
make integration-test
make coverage
make build
make run
```

### Com scripts

```bash
chmod +x scripts/*.sh
./scripts/dev.sh local
./scripts/dev.sh postgres
./scripts/test-integration.sh up
./scripts/test-integration.sh test
```

### Via compose da plataforma

A partir da raiz do repositorio:

```bash
docker compose -p pokedex -f core/docker-compose.yml up --build
```

## Testes

### Unitarios

```bash
go test ./tests/unit -v
```

### Integracao

```bash
make integration-test
```

## Cobertura

```bash
make coverage
```

Meta minima: 75% (ideal: 90%).
