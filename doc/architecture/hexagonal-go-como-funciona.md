# Hexagonal em Go — Cada Camada, o Que Faz, e Como a Linguagem Realiza Isso

> Guia prático que liga a **teoria** (doc/architecture/hexagonal.md) à **mecânica de Go**.
> Cada seção mostra: (1) o que a camada faz, (2) onde mora no `mobile-bff`, (3) como Go materializa isso.
> Foco: entender **por que** o código está organizado assim e **como** o compilador garante as regras.

---

## 1. O Desenho de Uma Vez

```
                    ┌─────────────────────────────────────────────────────┐
                    │                  MUNDO EXTERNO                       │
                    │   HTTP / JSON / rede / banco / serviços Go          │
                    └─────────────────────────────────────────────────────┘
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    ▼                     ▼                     ▼
         ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
         │  ADAPTADOR IN    │  │  ADAPTADOR IN    │  │  ADAPTADOR OUT   │
         │  handlers HTTP   │  │  middleware      │  │  client HTTP     │
         │  (net/http)      │  │  (auth, cors)    │  │  (catalog/auth)  │
         └────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘
                  │                     │                      │
                  │  usa                │  usa                │  implementa
                  ▼                     ▼                      │
         ┌─────────────────────────────────────────┐          │
         │           PORTAS (interfaces)           │          │
         │  ports/inbound/*  ports/outbound/*      │◄─────────┘
         └──────────────┬──────────────────────────┘
                        │
                        ▼
         ┌─────────────────────────────────────────┐
         │      HEXÁGONO (aplicação central)       │
         │  domain/   → entidades, regras, erros   │
         │  service/  → casos de uso               │
         └─────────────────────────────────────────┘
```

**Regra de ouro (Regra da Dependência):** setas apontam **para dentro** do hexágono.
O hexágono não importa nada de fora. Só `domain` + `ports` + `service` — nada de HTTP, pgx ou Redis.

---

## 2. As Camadas Uma a Uma

### 2.1 `domain/` — O Coração (entidades + regras + erros)

**O que faz:** define o vocabulário do negócio. Nada de infraestrutura aqui — nem HTTP, nem banco.

**Arquivos reais** (`core/bff/mobile-bff/internal/domain/`):
| Arquivo | Conteúdo |
|---------|----------|
| `pokemon.go` | `Pokemon`, `Type`, `PokemonPage`, `Favorite`, `User`... |
| `auth_session.go` | `AuthSession` (payload retornado pelo auth-service) |
| `errors.go` | erros sentinela: `ErrPokemonNotFound`, `ErrFavoriteAlreadyExists`... |
| `type_colors.go` | cor de cada tipo (mapa puro, sem DB) |

**Mecânica de Go:** são structs + `errors.New`. Nenhum import além de stdlib (`time`, `errors`). Se um dia `domain/pokemon.go` importar `github.com/jackc/pgx` — **violação**, o domínio deixou de ser puro.

```go
// internal/domain/pokemon.go
type Pokemon struct {
    ID     string
    Name   string
    Number string
    Types  []string
}

// internal/domain/errors.go — erros SENTINELA
var ErrPokemonNotFound = errors.New("pokemon nao encontrado")
```

---

### 2.2 `ports/` — As Portas (interfaces = contratos)

**O que faz:** declara **o que** a aplicação precisa, não **como**. Go não tem herança — **interfaces são a ferramenta central**.

**Mecânica de Go que muda tudo: satisfação implícita.**
Em Go, um tipo **implementa** uma interface apenas por ter os métodos certos. Nenhum `implements`, nenhum `extends`. Isso é o que torna Ports & Adapters natural em Go.

```
┌── ports/outbound/pokemon_repository.go ──┐
│  type PokemonRepository interface {      │  ← O QUE o hexágono exige
│      GetByID(ctx, id) (*domain.Pokemon, error)
│      GetAll(ctx, page, size) (*domain.PokemonPage, error)
│      Search(...)
│      ...
│  }                                       │
└──────────────────────────────────────────┘
        ▲ implementa implicitamente
        │
┌───────┴──────────────────────────────────┐
│  adapters/outbound/http/                 │
│  pokemon_catalog_client.go               │  ← COMO o hexágono obtém
│  type PokemonCatalogServiceRepository struct{...}
│  func (r *PokemonCatalogServiceRepository) GetByID(...)  // método idêntico
│  ...
└──────────────────────────────────────────┘
```

**Quem define a interface?** Quem **consome** (o hexágono), não quem implementa. Por isso as interfaces vivem em `ports/`, nunca no adapter. Isso inverte a dependência: o adapter (baixo nível) passa a depender do contrato (alto nível).

**Três pontos de confusão comuns:**

1. **"Interface no Go é tipo abstrato?"** — Não é classe. É um conjunto de métodos. Uma struct pode satisfazer N interfaces ao mesmo tempo sem declarar nada.
2. **`var _ Interface = (*Impl)(nil)`** — truque de compilação. Garante em tempo de build que `Impl` satisfaz `Interface`. Se alguém quebrar um método, o build falha na hora.
   ```go
   // pokemon_catalog_client.go:143
   var _ outbound.PokemonRepository = (*PokemonCatalogServiceRepository)(nil)
   // auth_service_client.go:296-297 — o MESMO client satisfaz DUAS portas
   var _ outbound.AuthProvider  = (*AuthServiceClient)(nil)
   var _ inbound.TokenValidator = (*AuthServiceClient)(nil)
   ```
3. **Adaptador que faz dois papéis:** `AuthServiceClient` implementa `outbound.AuthProvider` (login/signup para o `AuthService`) **e** `inbound.TokenValidator` (validação de token para o middleware). Uma struct, duas portas, zero declaração.

---

### 2.3 `service/` — Os Casos de Uso (aplicação)

**O que faz:** orquestra a lógica de negócio. Não sabe HTTP, não sabe Postgres — **só conhece interfaces de `ports/outbound`**.

**Mecânica de Go: injeção de dependência manual via construtor.**
O `service` recebe as portas como parâmetros do construtor e guarda nos campos. Em runtime, `main.go` decide **qual implementação concreta** entra.

```go
// internal/service/pokemon_service.go
type PokemonService struct {
    pokemonRepo  outbound.PokemonRepository   // ← interface, não *pgxpool.Pool
    favoriteRepo outbound.FavoriteRepository  // ← interface
}

func NewPokemonService(pokemonRepo outbound.PokemonRepository,
                       favoriteRepo outbound.FavoriteRepository) *PokemonService {
    return &PokemonService{pokemonRepo: pokemonRepo, favoriteRepo: favoriteRepo}
}
```

**O ponto-chave:** `pokemon_service.go` **não importa** `adapters/outbound/http` nem `postgres`. Ele só importa `domain`, `ports/inbound`, `ports/outbound`. Isso é o que permite testar com stub:

```go
// internal/service/service_test.go — substitui o adapter real por stub
svc := service.NewPokemonService(mockPokemonRepo, mockFavoriteRepo)
// mockPokemonRepo implementa outbound.PokemonRepository sem banco nenhum
```

---

### 2.4 `adapters/inbound/http/` — O que o mundo chama

**O que faz:** traduz HTTP/JSON → chamada de caso de uso. Recebe `*http.Request`, extrai dados, chama `inbound.*`, serializa resposta.

**Mecânica de Go: handlers são funções com assinatura `func(http.ResponseWriter, *http.Request)`.**
O `Handler` guarda as portas inbound e cada handler chama o caso de uso:

```go
// internal/adapters/inbound/http/handler.go
type Handler struct {
    pokemonUseCase  inbound.PokemonUseCase   // ← interface
    favoriteUseCase inbound.FavoriteUseCase  // ← interface
    authUseCase     inbound.AuthUseCase      // ← interface
}

// handler.go:56
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
    RespondJSON(w, http.StatusOK, h.responseBuilder.BuildHealthResponse())
}
```

**Regra arquitetural do projeto:** handler importa apenas `ports/inbound`. **Nunca** um repositório concreto. Se um handler precisar de dados, ele pede ao caso de uso — o caso de uso decide de onde vem.

**Rotas** (`RegisterRoutes`): Go 1.22+ permite patterns com método e path params:
```go
mux.HandleFunc("GET /api/v1/pokemons/{id}/details", h.GetPokemonDetails)
mux.HandleFunc("POST /api/v1/pokemons/{id}/favorite", h.RequireAuth(h.AddFavorite))
```

**Middleware:** `AuthMiddleware`, `CORSMiddleware`, `AuthRateLimitMiddleware` envolvem o mux. O `AuthMiddleware` recebe `inbound.TokenValidator` — novamente uma interface — e em `main.go` recebe o `AuthServiceClient` concreto.

---

### 2.5 `adapters/outbound/` — O que o hexágono chama

**O que faz:** implementa as portas outbound. Três tipos neste projeto:

| Adapter | Pasta | Implementa | Faz |
|---------|-------|------------|-----|
| HTTP client | `outbound/http/` | `outbound.PokemonRepository`, `outbound.AuthProvider`, `outbound.FavoriteCatalogProvider`, `inbound.TokenValidator` | Chama outros serviços Go via HTTP |
| PostgreSQL | `outbound/postgres/` | `outbound.PokemonRepository`, `outbound.FavoriteRepository` | Fala com o banco via `pgx/v5` |
| Memória | `outbound/memory/` | `outbound.*` | Mocks/fixtures para testes e fallback |

**Mecânica de Go: struct com um cliente concreto + métodos.**
O adapter guarda a dependência de infra (ex.: `*pgxpool.Pool` ou `*http.Client` com circuit breaker) e expõe métodos que batem com a interface da porta:

```go
// adapters/outbound/postgres/pokemon_repository.go
type PostgresPokemonRepository struct {
    db *pgxpool.Pool
}

func (r *PostgresPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    // SQL + pgx — INFRA aqui, escondida atrás da interface
}
```

**Normalização de erros no adapter:** o adapter traduz erros de infra para erros de **domínio** (`pgx.ErrNoRows` → `domain.ErrPokemonNotFound`). O caso de uso e o handler só conhecem erros de domínio.

```go
// postgres/pokemon_repository.go — err de banco vira err de domínio
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domain.ErrPokemonNotFound
}
```

**Equivalência estrutural:** o `PokemonCatalogServiceRepository` (HTTP) e o `PostgresPokemonRepository` (SQL) implementam **a mesma porta** `outbound.PokemonRepository`. Trocar o armazenamento = trocar uma linha no `main.go`, sem tocar em `service/` nem `handler/`.

---

### 2.6 `cmd/server/main.go` — Composition Root (o fio condutor)

**O que faz:** é o ÚNICO lugar que conhece as implementações concretas. Monta tudo na ordem certa: cria adaptadores → injeta nos serviços → injeta nos handlers → empilha middleware → sobe o servidor.

**Mecânica de Go: o `main` é quem decide as dependências reais.**

```go
// main.go (resumido)
// 1. Adapter concreto de catálogo (HTTP client)
pokemonRepo = httpclient.NewPokemonCatalogServiceRepository(cfg.PokemonCatalogServiceURL)

// 2. Adapter concreto de favoritos: por padrão via catalog-service (REST);
//    fallback para Postgres ou memória conforme config (FAVORITES_VIA_CATALOG).
favoriteRepo = favoriteCatalogClient      // default: httpclient.NewFavoriteCatalogClient(...)
// favoriteRepo = postgres.NewPostgresFavoriteRepository(db.Pool)   // se FAVORITES_VIA_CATALOG=false
// favoriteRepo = memory.NewFavoriteRepository()                    // fallback

// 3. Caso de uso recebe as INTERFACES (que na verdade são os adapters)
pokemonService := service.NewPokemonService(pokemonRepo, favoriteRepo)

// 4. Handler recebe o caso de uso
h := httpadapter.NewHandler(pokemonService, favoriteService, authService)
```

**A magia:** `NewPokemonService(pokemonRepo, favoriteRepo)` aceita *qualquer* tipo que satisfaça `outbound.PokemonRepository`. `main.go` passa o concreto. O teste passa um stub. **O mesmo construtor serve produção e teste.**

---

## 3. Como Go Realiza a Regra da Dependência

Em Go, **dependência = import**. A Regra da Dependência vira: *"quem pode importar quem"*.

```
IMPORTS REAIS (o que cada camada pode importar):

domain/          →  nada (só stdlib)
ports/inbound/   →  domain
ports/outbound/  →  domain
service/         →  domain, ports/inbound, ports/outbound
adapters/inbound/http/  →  ports/inbound, domain        (NUNCA adapters/outbound!)
adapters/outbound/*     →  ports/outbound, domain       (+ libs de infra)
cmd/server/main.go      →  TUDO (é o Composition Root)
```

**Por que isso é garantido pelo compilador:** se um `handler` importar `adapters/outbound/postgres`, o código **compila** mas viola a arquitetura. Para o Go, import ilegal ≠ erro de compilação. Por isso o projeto usa:

1. **Checklist arquitetural** no AGENTS.md (revisão manual/agent).
2. **`var _ Porta = (*Impl)(nil)`** — garante que o adapter satisfaz a porta.
3. **`harness/eval`** — avaliações que verificam se um LLM respeita as regras.
4. **Estrutura de pastas** — `internal/` impede import de fora do módulo; `adapters/` separado de `service/` torna a violação visível.

---

## 4. Fluxo Completo de uma Requisição (com os arquivos reais)

```
Cliente HTTP
   │  GET /api/v1/pokemons/1/details
   ▼
Kong Gateway (:8000)
   │  roteia para mobile-bff
   ▼
adapters/inbound/http/middleware.go     AuthMiddleware (valida JWT via inbound.TokenValidator
   │                                     = AuthServiceClient)
   ▼
adapters/inbound/http/handler.go        RegisterRoutes → Handler.GetPokemonDetails
   │  Handler guarda inbound.PokemonUseCase
   ▼
service/pokemon_service.go              PokemonService.GetPokemonDetails(ctx, id, userID)
   │  usa outbound.PokemonRepository.GetByID + outbound.FavoriteRepository.IsFavorite
   ▼
adapters/outbound/http/                 PokemonCatalogServiceRepository.GetByID
pokemon_catalog_client.go               faz HTTP GET http://pokemon-catalog-service:8081/v1/pokemons/1
   │  (com circuit breaker + retry)
   ▼
pokemon-catalog-service                 (outro serviço Go — fonte canônica do catálogo)
```

Cada `▼` cruza uma fronteira **por uma interface**, exceto dentro do próprio hexágono.

---

## 5. Por Que Isso Dá Certo em Go (resumo da mecânica)

| Pergunta | Resposta em Go |
|----------|----------------|
| Como declaro um contrato? | `type X interface { ... }` |
| Como uma struct cumpre o contrato? | **Implicitamente** — basta ter os métodos com a mesma assinatura |
| Quem define a interface? | Quem **consome** (o hexágono). Vive em `ports/`, nunca no adapter |
| Como o service chama o banco sem conhecer o banco? | Guarda `outbound.Repository` (interface) no campo; `main.go` injeta o concreto |
| Como troco Postgres por memória? | Troco 1 linha em `main.go`; `service` e `handler` nem recompilam diferente |
| Como sei que o adapter cumpre a porta? | `var _ outbound.PokemonRepository = (*Impl)(nil)` — falha de build se não cumprir |
| Como testo sem rede/banco? | Crio um stub que implementa a mesma interface e injeto no construtor |
| O que impede um handler de importar postgres? | Nada no compilador — é disciplina + checklist + evals + revisão |

---

## 6. O que é de Cada Camada — Regra de Bolso

```
Você está escrevendo...                      Coloque em...
--------------------------------------------------------------
Uma struct/regra do negócio                 domain/
Um erro de domínio ("nao encontrado")       domain/errors.go
Um contrato de repositório/client           ports/outbound/
Um contrato de caso de uso                  ports/inbound/
A lógica "o que fazer" (orquestração)       service/
Uma rota HTTP + parsing de JSON             adapters/inbound/http/
SQL / pgx / transação                       adapters/outbound/postgres/
Um HTTP client p/ outro serviço             adapters/outbound/http/
Conexão, pool, env vars                     adapters/outbound/postgres/database.go, config/
Wiring de tudo (o "cola tudo")              cmd/server/main.go
```

---

## 7. Erros Comuns de Quem Chega de Outra Linguagem

1. **"Onde fica o `implements`?"** — Não existe. Satisfação implícita. Se os métodos baterem, implementa.
2. **"Interface no `service` é redundante"** — É o contrário: é o que **permite** trocar a infra e testar sem banco.
3. **"Vou criar interface para TUDO"** — Errado. Cria interface **onde há dependência externa ou variação**. Não cria para struct de dados puro (`Pokemon` não precisa de interface).
4. **"Adapter é igual a controller"** — Não. Controller (MVC) acopla rota à lógica. Adapter **traduz protocolo** e delega para a porta; a lógica mora no `service`.
5. **"Struct embutida = herança"** — Composição, não herança. O `AuthServiceClient` satisfaz 2 interfaces **composicionalmente**, sem embutir nada.

---

## Referências

- [doc/architecture/hexagonal.md](hexagonal.md) — teoria (Alistair Cockburn)
- [doc/SOLID-AND-PATTERNS.md](../SOLID-AND-PATTERNS.md) — SOLID + patterns com exemplos Go
- Código real: `core/bff/mobile-bff/internal/**` (este documento cita arquivos e linhas exatas)
