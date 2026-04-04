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

> Based on: [Uber Go Style Guide — PT-BR](https://github.com/alcir-junior-caju/uber-go-style-guide-pt-br)

## When to use this skill

Load this skill when:
- You need practical, opinionated guidance on production Go code
- Reviewing code with a focus on performance, safety, and operational clarity
- Making decisions about mutexes, slices, maps, panics, and logging
- Writing Go code that needs to scale and be maintained by large teams

---

## 1. Performance Guidelines

### Initialize maps and slices with known capacity

```go
// ✅ Correct: avoids reallocations
pokemons := make([]domain.Pokemon, 0, len(ids))
index := make(map[string]*domain.Pokemon, len(ids))

// ❌ Incorrect: reallocates as it grows
var pokemons []domain.Pokemon
index := map[string]*domain.Pokemon{}
```

### Prefer `strconv` over `fmt` for numeric conversions

```go
// ✅ Faster
s := strconv.Itoa(42)
n, err := strconv.Atoi("42")

// ❌ Slower
s := fmt.Sprintf("%d", 42)
```

### Avoid format strings in constant errors

```go
// ✅ No string allocation
const errMsg = "token expired"

// ❌ Allocates a new string per call if using Sprintf
err := fmt.Errorf("token expired")  // acceptable; without Sprintf is ok
```

---

## 2. Struct Initialization

### Always use field names when initializing structs

```go
// ✅ Correct: readable and resilient to field order changes
svc := &PokemonService{
    pokemonRepo:  repo,
    favoriteRepo: favoriteRepo,
}

// ❌ Incorrect: breaks silently if field order changes
svc := &PokemonService{repo, favoriteRepo}
```

### Omit zero-value fields when initializing

```go
// ✅ Clean
cfg := Config{
    Host:    "localhost",
    Port:    8080,
}
// Timeout stays as zero value (0)

// ❌ Unnecessarily verbose
cfg := Config{
    Host:    "localhost",
    Port:    8080,
    Timeout: 0,
}
```

---

## 3. Errors

### Never silently discard errors

```go
// ❌ Silences the error — forbidden in production code
_ = repo.Save(ctx, pokemon)

// ✅ At least log the error if you cannot return it
if err := repo.Save(ctx, pokemon); err != nil {
    log.Printf("failed to save pokemon: %v", err)
}
```

### Consistent wrapping

```go
// ✅ Add context with %w to preserve the error chain
if err := s.authProvider.Login(ctx, email, password); err != nil {
    return nil, fmt.Errorf("authenticate user %s: %w", email, err)
}
```

### Error types vs sentinel errors

Use **sentinel errors** (`var ErrX = errors.New(...)`) for expected errors that callers check:

```go
var (
    ErrNotFound       = errors.New("not found")
    ErrInvalidToken   = errors.New("invalid token")
    ErrAlreadyExists  = errors.New("already exists")
)
```

Use **error types** when you need to carry additional data:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("field %s: %s", e.Field, e.Message)
}
```

---

## 4. Channels and Goroutines

### Goroutine with explicit channel size

```go
// ✅ Explicit buffer — reason documented
// buffered channel of 1 to avoid blocking the producer
results := make(chan *domain.Pokemon, 1)

// ❌ Unbuffered channel when the receiver may be slow — causes blocking
results := make(chan *domain.Pokemon)
```

### Avoid goroutines in Init functions

`init()` runs too early; goroutines launched there are hard to control:

```go
// ❌ Forbidden
func init() {
    go backgroundWorker()
}

// ✅ Start explicitly, with context
func NewScheduler(ctx context.Context) *Scheduler {
    s := &Scheduler{}
    go s.run(ctx)
    return s
}
```

---

## 5. Mutex and Synchronization

### Embed mutex next to what it protects

```go
// ✅ Clear: the mutex protects the fields below it
type MockPokemonRepository struct {
    mu       sync.RWMutex
    pokemons map[string]*domain.Pokemon
}
```

### Prefer `sync.RWMutex` when reads are more frequent than writes

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

## 6. Slices and Maps

### Never return nil slice and non-nil slice as equivalent

```go
// ✅ Always return an empty slice, not nil, when there are no items
func (s *Service) GetUserFavorites(ctx context.Context, userID string) ([]string, error) {
    if userID == "" {
        return []string{}, nil  // not nil
    }
    // ...
}
```

### Copy maps and slices when receiving or returning (API boundary)

```go
// ✅ Makes a copy to prevent the caller from mutating internal state
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

### `context.Context` is always the first parameter

```go
// ✅ Correct
func (s *PokemonService) GetPokemonDetails(ctx context.Context, id string) (*domain.PokemonDetail, error)

// ❌ Incorrect
func (s *PokemonService) GetPokemonDetails(id string, ctx context.Context) (*domain.PokemonDetail, error)
```

### Never store context in a struct

Context belongs to a request's flow, not to an object's state:

```go
// ❌ Forbidden
type Service struct {
    ctx context.Context
}

// ✅ Correct: pass as parameter
func (s *Service) Do(ctx context.Context) error { ... }
```

---

## 8. Linting and Tools

Uber recommends the following tools to enforce style:

| Tool | Purpose |
|------|---------|
| `golangci-lint` | Aggregated lint with multiple linters |
| `errcheck` | Detects discarded errors |
| `staticcheck` | Advanced static analysis |
| `go vet` | Basic compiler analysis |
| `gofmt` / `goimports` | Formatting and imports |

---

## 9. Constructor Options Pattern

When a constructor has many optional parameters, use functional options:

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

## Sources

- [Uber Go Style Guide — PT-BR](https://github.com/alcir-junior-caju/uber-go-style-guide-pt-br)
- [Uber Go Style Guide (original)](https://github.com/uber-go/guide/blob/master/style.md)
