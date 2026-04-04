---
name: golang-style-uber
description: >-
  Uber Go Style Guide for production Go code. Covers performance patterns (slice/map pre-allocation,
  strconv over fmt), struct initialization, error handling (no silent discards, wrapping, sentinel errors),
  goroutines with context and lifecycle, mutex placement, channel sizing, context rules, and
  functional options pattern.
  Use when writing performance-sensitive Go code, reviewing production patterns, or implementing
  constructors with optional parameters.
  Trigger examples: "pre-allocate slice", "functional options", "channel buffer", "mutex",
  "struct initialization", "production Go", "Uber style", "context in struct".
  Do NOT use for test patterns (use golang-testing) or API design (use go-api-design).
---

# Uber Go Style Guide — Skill

> Baseado em: [Uber Go Style Guide — PT-BR](https://github.com/alcir-junior-caju/uber-go-style-guide-pt-br)

## Quando usar esta skill

Carregue esta skill ao:
- Precisar de orientações práticas e opinionadas sobre código Go de produção
- Revisar código com foco em performance, segurança e clareza operacional
- Tomar decisões sobre mutexes, slices, maps, panics e logging
- Escrever código Go que precisa escalar e ser mantido por equipes grandes

---

## 1. Diretrizes de Performance

### Inicialize maps e slices com capacidade conhecida

```go
// ✅ Correto: evita realocações
pokemons := make([]domain.Pokemon, 0, len(ids))
index := make(map[string]*domain.Pokemon, len(ids))

// ❌ Incorreto: realoca conforme cresce
var pokemons []domain.Pokemon
index := map[string]*domain.Pokemon{}
```

### Prefira `strconv` a `fmt` para conversões numéricas

```go
// ✅ Mais rápido
s := strconv.Itoa(42)
n, err := strconv.Atoi("42")

// ❌ Mais lento
s := fmt.Sprintf("%d", 42)
```

### Evite strings de formato em erros constantes

```go
// ✅ Sem alocação de string
const errMsg = "token expirado"

// ❌ Aloca nova string a cada chamada se for Sprintf
err := fmt.Errorf("token expirado")  // aceitável; sem Sprintf é ok
```

---

## 2. Inicialização de Structs

### Sempre use nomes de campos ao inicializar structs

```go
// ✅ Correto: legível e resistente a mudanças de ordem de campos
svc := &PokemonService{
    pokemonRepo:  repo,
    favoriteRepo: favoriteRepo,
}

// ❌ Incorreto: quebra silenciosamente se a ordem dos campos mudar
svc := &PokemonService{repo, favoriteRepo}
```

### Omita campos zero-value ao inicializar

```go
// ✅ Limpo
cfg := Config{
    Host:    "localhost",
    Port:    8080,
}
// Timeout fica como zero value (0)

// ❌ Verboso sem necessidade
cfg := Config{
    Host:    "localhost",
    Port:    8080,
    Timeout: 0,
}
```

---

## 3. Erros

### Nunca descarte erros silenciosamente

```go
// ❌ Silencia o erro — proibido em código de produção
_ = repo.Save(ctx, pokemon)

// ✅ Pelo menos logue o erro se não puder retorná-lo
if err := repo.Save(ctx, pokemon); err != nil {
    log.Printf("falha ao salvar pokemon: %v", err)
}
```

### Wrapping consistente

```go
// ✅ Adicione contexto com %w para preservar a cadeia
if err := s.authProvider.Login(ctx, email, password); err != nil {
    return nil, fmt.Errorf("autenticar usuário %s: %w", email, err)
}
```

### Tipos de erro vs sentinel errors

Use **sentinel errors** (`var ErrX = errors.New(...)`) para erros esperados que os chamadores verificam:

```go
var (
    ErrNotFound       = errors.New("não encontrado")
    ErrInvalidToken   = errors.New("token inválido")
    ErrAlreadyExists  = errors.New("já existe")
)
```

Use **tipos de erro** quando precisar transportar dados adicionais:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("campo %s: %s", e.Field, e.Message)
}
```

---

## 4. Channels e Goroutines

### Goroutine com tamanho explícito de canal

```go
// ✅ Buffer explícito — razão documentada
// canal com buffer de 1 para evitar bloqueio do produtor
results := make(chan *domain.Pokemon, 1)

// ❌ Canal sem buffer quando o receptor pode demorar — causa bloqueio
results := make(chan *domain.Pokemon)
```

### Evite goroutines em funções Init

`init()` é executado cedo demais; goroutines lançadas ali são difíceis de controlar:

```go
// ❌ Proibido
func init() {
    go backgroundWorker()
}

// ✅ Inicie explicitamente, com context
func NewScheduler(ctx context.Context) *Scheduler {
    s := &Scheduler{}
    go s.run(ctx)
    return s
}
```

---

## 5. Mutex e Sincronização

### Incorpore mutex junto ao que ele protege

```go
// ✅ Claro: o mutex protege os campos abaixo dele
type MockPokemonRepository struct {
    mu       sync.RWMutex
    pokemons map[string]*domain.Pokemon
}
```

### Prefira `sync.RWMutex` quando leituras são mais frequentes que escritas

```go
func (m *MockPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    // ...
}

func (m *MockPokemonRepository) Save(ctx context.Context, p *domain.Pokemon) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    // ...
}
```

---

## 6. Slices e Maps

### Nunca retorne nil slice e non-nil slice como equivalentes

```go
// ✅ Sempre retorne slice vazio, não nil, quando não há itens
func (s *Service) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
    if userID == "" {
        return []string{}, nil  // não nil
    }
    // ...
}
```

### Copie maps e slices ao receber ou retornar (boundary da API)

```go
// ✅ Faz cópia para evitar que o chamador mute o estado interno
func (m *MockPokemonRepository) GetAll() []domain.Pokemon {
    m.mu.RLock()
    defer m.mu.RUnlock()

    result := make([]domain.Pokemon, len(m.pokemons))
    copy(result, m.pokemons)
    return result
}
```

---

## 7. Context

### `context.Context` é sempre o primeiro parâmetro

```go
// ✅ Correto
func (s *PokemonService) GetPokemonDetails(ctx context.Context, id string) (*domain.PokemonDetail, error)

// ❌ Incorreto
func (s *PokemonService) GetPokemonDetails(id string, ctx context.Context) (*domain.PokemonDetail, error)
```

### Nunca armazene context em struct

Context pertence ao fluxo de uma requisição, não ao estado de um objeto:

```go
// ❌ Proibido
type Service struct {
    ctx context.Context
}

// ✅ Correto: passe como parâmetro
func (s *Service) Do(ctx context.Context) error { ... }
```

---

## 8. Linting e Ferramentas

O Uber recomenda o uso das seguintes ferramentas para reforçar o estilo:

| Ferramenta | Propósito |
|------------|-----------|
| `golangci-lint` | Lint agregado com múltiplos linters |
| `errcheck` | Detecta erros descartados |
| `staticcheck` | Análise estática avançada |
| `go vet` | Análise básica do compilador |
| `gofmt` / `goimports` | Formatação e imports |

---

## 9. Padrões de Opções para Construtores

Quando um construtor tiver muitos parâmetros opcionais, use functional options:

```go
type ClientOption func(*AuthServiceClient)

func WithTimeout(d time.Duration) ClientOption {
    return func(c *AuthServiceClient) {
        c.httpClient.Timeout = d
    }
}

func NewAuthServiceClient(baseURL string, opts ...ClientOption) *AuthServiceClient {
    c := &AuthServiceClient{
        baseURL:    baseURL,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}
```

---

## Fontes

- [Uber Go Style Guide — PT-BR](https://github.com/alcir-junior-caju/uber-go-style-guide-pt-br)
- [Uber Go Style Guide (original)](https://github.com/uber-go/guide/blob/master/style.md)
