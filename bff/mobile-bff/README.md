# Mobile BFF - API da Pokedex

## Visao Geral

O Mobile-BFF e um Backend-for-Frontend (BFF) em Go para atender clientes mobile/web da Pokedex.
O projeto usa arquitetura hexagonal (Ports and Adapters), PostgreSQL, testes unitarios e testes de integracao.

## Arquitetura

### Padrao Hexagonal (Ports and Adapters)

```text
domain/                    # Regras de negocio (independente de framework)
├─ pokemon.go              # Modelos (Pokemon, PokemonDetail etc.)
└─ errors.go               # Erros de dominio

ports/                     # Contratos (interfaces)
├─ repository.go           # Portas de acesso a dados
└─ usecase.go              # Portas de casos de uso

service/                   # Servicos de aplicacao
└─ pokemon_service.go      # Operacoes de Pokemon e favoritos

adapters/
├─ http/                   # Adapter de entrada HTTP
│  ├─ handlers.go          # Endpoints REST
│  ├─ middleware.go        # Auth e CORS
│  ├─ response_builder.go  # Conversao dominio -> DTO
│  └─ dto/                 # Objetos de transferencia
└─ repository/             # Adapter PostgreSQL
   ├─ pokemon_repository.go
   ├─ favorite_repository.go
   └─ database.go          # Pool de conexao

tests/
├─ unit/                   # Testes unitarios
├─ integration/            # Testes de integracao com banco
└─ mocks/                  # Repositorios fake
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

- `POST /api/v1/pokemons/{id}/favorite`
- `DELETE /api/v1/pokemons/{id}/favorite`

## Contrato Entre BFF E Service

O BFF atende o contrato rico do front, enquanto o `pokedex-service` concentra o catalogo canonico de Pokemon.

### BFF (mobile-bff)

- Compoe respostas para UI (`home`, cards e detalhe rico)
- Aplica regras de experiencia para cliente mobile/web
- Mantem regra de favoritos por usuario (escrita)

### Service (pokedex-service)

- Exposicao de dados canonicos de catalogo (leitura)
- Endpoints consumidos pelo BFF:
   - `GET /v1/pokemons`
   - `GET /v1/pokemons/search`
   - `GET /v1/pokemons/type/{type}`
   - `GET /v1/pokemons/{id}`

### Decisao Atual

- Catalogo: no `pokedex-service`
- Favoritos: no BFF por usuario autenticado
- Evolucao futura: extrair favoritos para um servico proprio quando login/cadastro estiverem integrados

## Stack

- Go 1.22
- net/http
- PostgreSQL + pgx/v5
- Testify
- Docker / Docker Compose

## Como Executar

### Modo local (mock)

```bash
export MOBILE_BFF_PORT=8080
go run ./cmd/server/main.go
```

### Modo PostgreSQL

```bash
docker compose -f docker-compose.test.yml up -d
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pokedex"
export MOBILE_BFF_PORT=8080
go run ./cmd/server/main.go
```

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
