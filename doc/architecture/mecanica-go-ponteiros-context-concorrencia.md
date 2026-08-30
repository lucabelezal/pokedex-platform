# Mecanismos de Go na Prática — Ponteiros, Context, Concorrência, Generics

> Complemento de [hexagonal-go-como-funciona.md](hexagonal-go-como-funciona.md).
> Cada mecanismo é explicado com **código real deste projeto** (`core/bff/mobile-bff`).
> Objetivo: entender **como a linguagem funciona** por baixo, no contexto de um código que você já conhece.

---

## 1. Ponteiros (`*T`) — Valor vs Referência

### O que é um ponteiro

Um ponteiro guarda o **endereço de memória** de um valor, não o valor em si. `*T` é "ponteiro para T". `&x` cria um ponteiro para `x`. `*p` lê/escreve o valor apontado.

### Onde o projeto usa ponteiros

**1. Structs retornadas como ponteiro — "não existe" é `nil`:**
```go
// service/pokemon_service.go
func (s *PokemonService) GetPokemonDetails(...) (*domain.PokemonDetail, error) {
    pokemon, err := s.pokemonRepo.GetByID(ctx, pokemonID)
    if err != nil {
        return nil, err        // ← nil = "nada aqui"
    }
    return detail, nil
}
```
O ponteiro `*domain.PokemonDetail` permite retornar `nil` quando o pokemon não existe. Se fosse valor (`domain.PokemonDetail`), você teria que criar um struct vazio e o chamador não saberia diferenciar "vazio" de "não encontrado".

**2. Campo opcional — "não sei o valor" também é `nil`:**
```go
// domain/pokemon.go
type PokemonDetail struct {
    GenderMale   *float64   // ← ponteiro = "pode não ter valor"
    GenderFemale *float64
}
```
Um `float64` sempre tem valor (0.0). Um `*float64` pode ser `nil` — perfeito para "a proporção de macho não foi informada". O catalog-service faz:
```go
if genderMale >= 0 {
    detail.GenderMale = &genderMale   // só preenche se existir
}
```

**3. Receiver de método — mutar ou evitar cópia:**
```go
// adapters/outbound/postgres/favorite_repository.go
func (r *PostgresFavoriteRepository) AddFavorite(...) error {
    // r é *PostgresFavoriteRepository → não copia o struct inteiro
}
```

### Regra de bolso

| Situação | Usa |
|----------|-----|
| "Não existe" / "opcional" / "não tem valor" | `*T` + `nil` |
| Struct grande que só vai ser lida | ponteiro `*T` (evita cópia) |
| Struct que o receiver vai mutar | receiver de ponteiro `(r *T)` |
| Value object pequeno e imutável (cor, tipo) | valor `T` |
| Slice/map/interface já são referências | não precisa de ponteiro |

### Armadilhas

- `nil` de ponteiro ≠ `nil` de interface. Um `(*Pokemon)(nil)` dentro de uma interface **não é** `== nil`:
  ```go
  var p *Pokemon = nil
  var any interface{} = p   // any NÃO é nil! tem tipo *Pokemon e valor nil
  if any == nil { /* false! */ }
  ```
- Chamar método em ponteiro `nil` panica se o método acessar campos.

---

## 2. `context.Context` — o "fio" que atravessa as camadas

### O que é

`context.Context` é um objeto que viaja **como primeiro parâmetro** de toda função que toca I/O. Ele carrega três coisas:

| Carga | O que faz | Exemplo no projeto |
|-------|-----------|--------------------|
| **Deadline/timer** | "você tem até X para responder" | `context.WithTimeout` |
| **Cancelamento** | "pode parar, o cliente foi embora" | o cancelamento do timeout |
| **Valores** | dados de request-scope (userID, email) | `context.WithValue` |

### 1. Timeout — em todo handler

```go
// adapters/inbound/http/handler.go:62
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()   // ← IMPORTANTE: libera o timer quando o handler terminar
    userID := getUserIDFromContext(ctx)
    ...
}
```

`r.Context()` é o contexto da requisição HTTP. `WithTimeout` cria um filho que **cancela sozinho após 5s** (ou quando `cancel()` for chamado). Esse `ctx` desce pela pilha até o Postgres.

### 2. Propagação — o context desce, nunca sobe

```go
handler (ctx, cancel)
  → service.GetPokemonDetails(ctx, ...)
    → repo.GetByID(ctx, ...)
      → pgx Query(ctx, ...)   // o driver postgres respeita o cancelamento
```

O `context.Context` **flui para dentro** (handler → service → repo → driver). Ele é sempre o **primeiro** parâmetro. Nunca é armazenado num struct (exceção: raro).

O benefício: se o cliente HTTP desconectar, o `r.Context()` cancela, o timeout dispara, e o `pgx` **aborta a query em voo** em vez de deixá-la rodando até o fim.

### 3. Valores — userID via contexto

```go
// adapters/inbound/http/middleware.go:187
ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
ctx = context.WithValue(ctx, UserEmailContextKey, userEmail)
next.ServeHTTP(w, r.WithContext(ctx))
```

O middleware extrai o userID do JWT, coloca no contexto, e o handler lê depois:
```go
// middleware.go:313
func getUserIDFromContext(ctx context.Context) string {
    userID, ok := ctx.Value(UserIDContextKey).(string)  // type assertion
    if !ok { return "" }
    return userID
}
```

### Boas práticas observadas no projeto

- **Chave tipada, não string crua**: `type contextKey string; const UserIDContextKey contextKey = "userID"` — evita colisão entre packages.
- **`defer cancel()` sempre** — vaza timer se esquecer.
- **Timeout pequeno por camada** — handler 5s, client HTTP 5-10s, DB via contexto.

---

## 3. Goroutines — o que é e o que NÃO é

### O modelo

Go lida com concorrência através de **goroutines**: funções que rodam "em paralelo" gerenciadas pelo runtime, não pelo sistema operacional. São baratas (começam com ~2KB de stack).

```go
go minhaFuncao()   // roda em uma goroutine separada
```

### Concorrência ≠ Paralelismo (a pergunta clássica)

| Conceito | Definição | Analogia |
|----------|-----------|----------|
| **Concorrência** | Várias tarefas **em progresso** no mesmo período, intercalando | Um chef alterna entre 3 panelas na mesma boca de fogão |
| **Paralelismo** | Várias tarefas **executando de verdade ao mesmo tempo** | 3 chefs, 3 bocas, cada um na sua |

Go permite **concorrência** por design (goroutines + channels). **Paralelismo** depende de hardware (vários cores) e é ativado pelo runtime automaticamente quando `GOMAXPROCS > 1` (padrão = nº de núcleos).

**Concorrência é sobre ESTRUTURA do código** (compor tarefas independentes).
**Paralelismo é sobre EXECUÇÃO** (rodar ao mesmo tempo).
Você escreve código **concorrente**; o runtime decide se executa em paralelo.

### Onde o projeto usa concorrência de verdade

**`sync.Mutex` no rate limiter em memória** — protege o mapa contra acesso simultâneo:

```go
// adapters/inbound/http/rate_limiter.go:66
type inMemoryRateLimiter struct {
    mu      sync.Mutex              // ← trava
    entries map[string]rateLimitEntry
}

func (l *inMemoryRateLimiter) Allow(_ context.Context, clientID string) bool {
    l.mu.Lock()                     // ← adquire
    defer l.mu.Unlock()             // ← garante liberação
    ...
}
```

**Por que precisa?** Cada requisição HTTP roda em uma goroutine do `net/http`. Duas requisições simultâneas do mesmo IP podem ler/escrever `l.entries` ao mesmo tempo — **race condition**. O `Mutex` serializa o acesso: só uma goroutine mexe no mapa por vez.

### Race condition — o inimigo

Duas goroutines lendo/escrevendo o mesmo dado sem sincronização = **data race**. Detectar:

```bash
go test -race ./...    # detector embutido — rode SEMPRE em CI
```

O projeto roda `-race` no CI (`.github/workflows/go-ci.yml`) e no Makefile.

### Mecânica de `sync` no projeto

| Ferramenta | Uso | Onde |
|------------|-----|------|
| `sync.Mutex` | Proteger dado compartilhado | `rate_limiter.go` (map de entradas) |
| `sync.RWMutex` | Muitas leituras, poucas escritas | não usado (Mutex basta aqui) |
| `sync.WaitGroup` | Esperar N goroutines terminarem | não usado no bff atual |
| `context` | Cancelamento cooperativo | em todos os handlers |

### Channels — quando usar

O projeto usa Redis e `sync.Mutex` em vez de channels porque o rate limiting é **estado compartilhado** (um mapa de contadores). Channels são para **comunicação entre goroutines** (passar valores). Regra de bolso do próprio Go:

> "Não comunique compartilhando memória; compartilhe memória comunicando." — mas: se só há 1 goroutine escrevendo e outras lendo, um Mutex é mais simples e legível.

Para um worker pool, fan-out/fan-in, você usaria channels — não é o caso do BFF atual (as operações são request-response síncronas).

---

## 4. Generics — tipos parametrizados

### O que é

Desde Go 1.18, funções e tipos podem ter **parâmetros de tipo** entre `[...]`. Permite escrever UMA função que serve para vários tipos, com type-safety em tempo de compilação.

### Onde o projeto usa

**`any`** = atalho para `interface{}` (aceita qualquer tipo). Usado quando a função precisa decodificar JSON sem saber o tipo:

```go
// adapters/outbound/http/pokemon_catalog_client.go:118
func (r *PokemonCatalogServiceRepository) getJSON(ctx context.Context, endpoint string, out any) (int, error) {
    ...
    if err := json.NewDecoder(resp.Body).Decode(out); err != nil { ... }
}
```
`out any` aceita `*domain.Pokemon`, `*domain.PokemonPage`, `[]domain.Type`... o chamador decide:
```go
var out domain.PokemonPage
_, err := r.getJSON(ctx, endpoint, &out)   // decode direto no tipo certo
```

### `any` vs generics de verdade

`any` é o **escape hatch** (perde type-safety). **Generics** mantêm type-safety:

```go
// Exemplo didático — se o projeto precisasse
func Chunk[T any](items []T, size int) [][]T {
    var out [][]T
    for i := 0; i < len(items); i += size {
        end := min(i+size, len(items))
        out = append(out, items[i:end])
    }
    return out
}

chunks := Chunk(pokemons, 20)   // [][]domain.Pokemon — o compilador sabe o tipo
```

### Quando usar generics (e quando não)

| Usa | Não usa |
|-----|---------|
| Coleção/container genérico (stack, filter, map) | Quando um `interface{}`/`any` basta |
| Algoritmo idêntico para vários tipos | Quando tipos têm métodos diferentes → interface é melhor |
| Helpers de slice/map reutilizáveis | Só para "parecer elegante" |

**Interface vs generics — como decidir:**
- **Interface** = "qualquer coisa que tenha estes métodos" (comportamento polimórfico). É o coração do hexagonal.
- **Generics** = "qualquer tipo T, mantendo o tipo exato" (estrutura uniforme). Ex.: container, coleção.

O projeto usa interfaces (hexagonal) em toda parte e `any` nos pontos de decodificação. Não abusa de generics — coerente com "Go simples".

---

## 5. Slices e Maps — as estruturas que você vai usar 90% do tempo

### Slices são "janelas" sobre arrays

```go
pokemons := make([]domain.Pokemon, 0)       // vazio, pronto pra append
pokemons = append(pokemons, p)              // cresce
```

**Padrão do projeto — retornar slice vazio, não nil:**
```go
// adapters/outbound/postgres/pokemon_repository.go
favorites := make([]string, 0)   // ← NUNCA var favorites []string (que é nil)
```
**Por quê?** `nil` serializa como `null` no JSON; `make([]T,0)` serializa como `[]`. A AGENTS.md exige: "use `make([]T, 0)` ao retornar slices via API/JSON" — para o cliente receber `[]` e não `null`.

### capacity — pré-alocar para performance

```go
// response_builder.go
pokemons := make([]dto.RichPokemonResponse, len(page.Content))  // tamanho exato
types := make([]dto.TypeDTO, len(p.Types))                      // sem append desnecessário
```

Quando você sabe o tamanho, `make([]T, n)` pré-aloca e evita re-alocações no `append`.

### Maps — cuidado com nil e concorrência

```go
entries: make(map[string]rateLimitEntry)   // inicializado no construtor
```
- Ler mapa `nil` funciona; **escrever** em mapa `nil` PANICA.
- Maps **não são thread-safe** → por isso o rate limiter usa `sync.Mutex`.
- `map[string]struct{}` = set (conjunto). O projeto usa para `favoriteSet` no handler:
  ```go
  favoriteSet := make(map[string]struct{}, len(favorites))
  for _, id := range favorites { favoriteSet[id] = struct{}{} }
  // depois: _, isFav := favoriteSet[id]  → bool
  ```

---

## 6. Ponteiros, Interface e o "nil trap" — os 3 erros clássicos

### Erro 1 — comparar `nil` com interface
```go
var err error           // nil
var p *Pokemon = nil
var i interface{} = p   // i tem TIPO *Pokemon, então NÃO é nil
if i == nil { /* false! erro clássico */ }
```

### Erro 2 — receiver nil
```go
func (s *PokemonService) List(...) {
    s.pokemonRepo.GetAll(...)   // se s for nil → panic
}
```
Por isso o construtor valida deps quando necessário (`service/auth_service.go` checa `s.authProvider == nil`).

### Erro 3 — escrever em mapa nil
```go
var m map[string]int
m["x"] = 1   // PANIC: assignment to entry in nil map
```

---

## 7. Como tudo se encaixa numa requisição (contexto + ponteiro + concorrência)

```
Cliente HTTP
  │
  ▼
goroutine do net/http (1 por requisição — já é concorrência!)
  │
  ▼
middleware: cria ctx com userID (WithValue) ─────────────────────────┐
  │                                                                 │
  ▼                                                                 │
handler: ctx, cancel := WithTimeout(r.Context(), 5s)                │
  │        defer cancel()                                           │
  ▼                                                                 │
service.GetPokemonDetails(ctx, id, userID)                          │
  │   retorna *domain.PokemonDetail (ponteiro, nil = não existe)    │
  ▼                                                                 │
repo.GetByID(ctx, id)                                               │
  │   retorna *domain.Pokemon                                       │
  ▼                                                                 │
pgx Query(ctx, ...)  ← o ctx carrega deadline; se estourar, aborta  │
  │                                                                 │
  └── o ctx viajou por TODAS as camadas ────────────────────────────┘
```

- **Concorrência**: o `net/http` já sobe uma goroutine por request; o Mutex do rate limiter protege o estado compartilhado entre elas.
- **Ponteiro**: `*domain.Pokemon` permite `nil` = não encontrado; `*float64` = campo opcional.
- **Context**: carrega timeout + userID e atravessa todas as camadas como 1º parâmetro.
- **Generics/`any`**: `getJSON(..., out any)` decodifica JSON em qualquer struct.

---

## 8. Cheatsheet rápido

| Você quer... | Use |
|--------------|-----|
| "não existe" ou opcional | `*T` + `nil` |
| Evitar cópia de struct grande | `*T` |
| Timeout numa operação | `context.WithTimeout(ctx, 5*time.Second)` + `defer cancel()` |
| Passar userID entre camadas | `context.WithValue(ctx, key, val)` com key tipada |
| Proteger mapa/mutável compartilhado | `sync.Mutex` + `Lock/Unlock` |
| Detectar race | `go test -race ./...` |
| Aceitar "qualquer tipo" sem type-safe | `any` (ex.: decode JSON) |
| Manter tipo exato com vários tipos | generics `func F[T any](...)` |
| Retornar lista vazia em API | `make([]T, 0)`, nunca `nil` |
| "set" (conjunto) | `map[T]struct{}` |
| Rodar função em paralelo | `go func()` (mas sincronize!) |

---

## Referências

- Código real citado: `core/bff/mobile-bff/internal/**`
- [Effective Go](https://go.dev/doc/effective_go) — ponteiros, slices, maps, concorrência
- [Go Concurrency Patterns](https://go.dev/blog/concurrency-patterns) — goroutines/channels (Rob Pike)
- [context package](https://pkg.go.dev/context) — doc oficial
- [Generics](https://go.dev/doc/tutorial/generics) — tutorial oficial
- Projeto: `doc/aprendizado/` — guias de estudo em pt-BR (`fase-04-concorrencia.md`, `fase-03-structs-interfaces.md`)
