---
name: go-style-combined
description: >-
  Complete Go style guide for the Pokedex Platform, combining the Uber Go Style Guide and the Google Go Style Guide
  with project-specific decisions already made. Covers naming, receivers, initialisms, 3-group imports,
  compile-time interface assertions, make() for slices, time.Duration for time constants,
  errors (wrapping, sentinels, strings), goroutines with context and lifecycle, mutex, slices/maps,
  struct init, grouped declarations, table-driven tests, functional options, panic, init(),
  mutable global variables, atomic, embed, and linting.
  Use for ANY Go code review or writing in the platform. This skill is the canonical reference.
  Trigger examples: "Go style", "naming", "receiver", "import", "error", "goroutine", "mutex",
  "nil slice", "interface compliance", "time.Duration", "table test", "functional options",
  "idiomatic Go", "code review", "review code".
applyTo:
  - "core/bff/mobile-bff/**/*.go"
  - "core/app/**/*.go"
  - "core/gateway/**/*.go"
---

# Go Style Guide — Pokedex Platform

> Canonical synthesis of the **Uber Go Style Guide** + **Google Go Style Guide**  
> with project-specific decisions documented.
>
> Language for code comments: **Portuguese (Brazil)**  
> When there is a conflict between guides, **local consistency** takes precedence (Google principle §Consistency).

---

## Contents

1. [Principles](#1-principles)
2. [Naming](#2-naming)
3. [Comments and Godoc](#3-comments-and-godoc)
4. [Imports](#4-imports)
5. [Declarations and groupings](#5-declarations-and-groupings)
6. [Struct, slice, and map initialization](#6-struct-slice-and-map-initialization)
7. [Interfaces and compile-time assertions](#7-interfaces-and-compile-time-assertions)
8. [Errors](#8-errors)
9. [Context](#9-context)
10. [Goroutines and concurrency](#10-goroutines-and-concurrency)
11. [Mutex and synchronization](#11-mutex-and-synchronization)
12. [Performance](#12-performance)
13. [Control flow](#13-control-flow)
14. [Table-driven tests](#14-table-driven-tests)
15. [Functional Options](#15-functional-options)
16. [Panic and init()](#16-panic-and-init)
17. [Mutable global variables](#17-mutable-global-variables)
18. [Linting and tools](#18-linting-and-tools)
19. [Project-specific decisions](#19-project-specific-decisions)

---

## 1. Principles

The Google Go Style Guide lists the attributes of readable code, in order of importance:

| Priority | Attribute | Description |
|:---:|---|---|
| 1 | **Clarity** | The purpose and reasoning are obvious to the reader |
| 2 | **Simplicity** | The goal is achieved in the simplest possible way |
| 3 | **Conciseness** | High signal-to-noise ratio — every line counts |
| 4 | **Maintainability** | Easy to maintain over time |
| 5 | **Consistency** | Consistent with the rest of the codebase |

> **Golden rule**: prefer the simplest mechanism that solves the problem.  
> Native channel, slice, map, or loop > stdlib > third-party > custom code.

---

## 2. Naming

### 2.1 Packages

- Lowercase, no underscore, no camelCase. Ex.: `pokemon`, `auth`, `http`.
- No info/helper/util/common/shared/lib suffix — these names communicate nothing.
- The package name is part of the symbol: `auth.NewService()`, not `auth.NewAuthService()`.
- Not plural: `net/url`, not `net/urls`.

```go
// ✅
package pokemon

// ❌
package pokemonUtils
package pokemonHelper
```

### 2.2 Receivers

- 1–2 letters, abbreviation of the type. **Never** `this`, `self`, `me`.
- Consistent across all methods of the same type.

```go
// ✅
func (s *PokemonService) List(ctx context.Context) ([]domain.Pokemon, error)
func (c *AuthServiceClient) Login(ctx context.Context, email, password string) (*domain.Session, error)
func (r *PostgresPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error)

// ❌
func (this *PokemonService) List(ctx context.Context) ([]domain.Pokemon, error)
func (service *PokemonService) List(ctx context.Context) ([]domain.Pokemon, error)
```

### 2.3 Initialisms

Always in full case (or all lowercase for unexported):

| Term | Exported | Unexported |
|------|----------|------------|
| ID | `ID` | `id` |
| URL | `URL` | `url` |
| HTTP | `HTTP` | `http` |
| JSON | `JSON` | `json` |
| DB | `DB` | `db` |
| API | `API` | `api` |
| gRPC (exported) | `GRPC` | `gRPC` |

```go
// ✅
type PokemonID string
func GetByID(id PokemonID) {}
var baseURL = "http://..."
type GRPCClient struct{}   // exported
var gRPCClient *GRPCClient // unexported

// ❌
type PokemonId string
func GetById(id PokemonId) {}
var baseUrl = "http://..."
```

### 2.4 Functions and methods

- **No `Get` prefix** for getters: `Count()` not `GetCount()`, `Name()` not `GetName()`.
  - Exception: expensive operation (RPC, IO) may use `Fetch`, `Compute`, `Load`.
- Constructors: `New` + type. Ex.: `NewPokemonService`, `NewHandler`, `NewAuthServiceClient`.
- Boolean predicates: `Is`, `Has`, `Can`, `Should`. Ex.: `IsFavorite`, `HasAccess`.
- Test function names with underscore for grouping: `TestListPokemons_WhenUserNotFound`.

### 2.5 Constants

- MixedCaps. Ex.: `MaxRetryAttempts`, `DefaultPageSize`.
- **Never** ALL_CAPS or `K` prefix. Ex.: ❌ `MAX_RETRY`, ❌ `kDefaultPage`.
- Name by role, not by value: ❌ `Twelve = 12`.

### 2.6 Variables

- Length proportional to scope. Short loop → `i`, `v`, `k`. Large scope → `userCount`.
- Omit the type from the name: `users` not `userSlice`; `count` not `numUsers`.
- Local names by content in the current context, not by origin.

---

## 3. Comments and Godoc

### 3.1 Project rules (PT-BR)

All comments must be in **Portuguese (Brazil)**.

#### When to comment (high value)

- Explanation of non-obvious logic ("why", not "what")
- Design/architecture reasons: why this pattern and not another
- Performance warnings, known limitations
- Integration context between services

#### When NOT to comment (zero value)

- Obvious from the name → no comment
- Describes the code literally: ❌ `// checks error`
- Self-explanatory code with descriptive naming

### 3.2 Godoc

Every **exported** symbol must have a godoc comment starting with its name:

```go
// PokemonService manages listing, search, and favorites operations for Pokémons.
type PokemonService struct { ... }

// GetByID returns a Pokémon by its canonical ID.
// Returns domain.ErrNotFound if the Pokémon does not exist in the catalog.
func (s *PokemonService) GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
```

#### Package comments:

```go
// Package pokemon implements use cases related to the Pokémon catalog.
// Depends on PokemonRepository (outbound) and implements PokemonUseCase (inbound).
package pokemon
```

### 3.3 Naked parameters

Add inline comments `/* isLocal */` when a boolean/int passed without a name is ambiguous:

```go
// ✅
printInfo("pikachu", true /* isElectric */, false /* isFavorite */)

// Better still: use a named type
type ElectricType bool
```

---

## 4. Imports

### 4.1 Three groups (project decision)

The project uses **3 groups** separated by blank lines. Uber uses 2, Google uses 4+.  
Decision: keep 3 for consistency with the current codebase.

```go
import (
    // 1. stdlib
    "context"
    "fmt"
    "net/http"
    "time"

    // 2. third-party
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    // 3. internal (pokedex-platform module)
    "pokedex-platform/core/bff/mobile-bff/internal/domain"
    inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
    outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)
```

### 4.2 Import alias

Use alias when:
- The package name doesn’t match the last element of the path (e.g., versioned paths `v2`)
- There is a collision between two packages with the same name
- The project has already established a canonical alias (e.g., `inbound`, `outbound`)

```go
// ✅
inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"

// ❌ — unnecessary rename
pokemonDomain "pokedex-platform/core/bff/mobile-bff/internal/domain"
```

### 4.3 Dot and blank imports

- `import .` → **forbidden** (obscures the origin of symbols)
- `import _` → only in `main` or test files that need side effects

---

## 5. Declarations and groupings

### 5.1 Group related declarations

```go
// ✅
const (
    DefaultPageSize = 20
    MaxPageSize     = 100
)

var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidToken = errors.New("invalid token")
)

// ❌ — mixes unrelated items in the same block
const (
    DefaultPageSize = 20
    EnvKey          = "APP_ENV"  // different concept → separate block
)
```

### 5.2 Variable declaration

```go
// var x T — explicit and intentional zero value
var count int
var repo outbound.PokemonRepository

// x := value — initializes with non-zero value
name := "Pikachu"
client := &http.Client{Timeout: 5 * time.Second}

// var x T when all fields are zero
var cfg Config  // not: cfg := Config{}
```

### 5.3 Minimum scope

Reduce variable scope without sacrificing readability:

```go
// ✅
if err := os.WriteFile(name, data, 0644); err != nil {
    return err
}

// use long declaration when the result is used outside the if
data, err := os.ReadFile(name)
if err != nil {
    return err
}
fmt.Println(data)
```

---

## 6. Struct, slice, and map initialization

### 6.1 Structs: always with field names

```go
// ✅
svc := &PokemonService{
    pokemonRepo:  repo,
    favoriteRepo: favoriteRepo,
    authProvider: authProvider,
}

// ❌ — breaks if field order changes
svc := &PokemonService{repo, favoriteRepo, authProvider}
```

Exception: test tables with ≤3 simple fields may omit names.

### 6.2 Omit zero-value fields

```go
// ✅
cfg := Config{
    Host: "localhost",
    Port: 8080,
}

// ❌
cfg := Config{
    Host:    "localhost",
    Port:    8080,
    Timeout: 0,     // zero value — unnecessary
    Debug:   false, // zero value — unnecessary
}
```

Exception: in tests, include zero-value fields when the name provides relevant context:

```go
tests := []struct {
    give string
    want int
}{
    {give: "0", want: 0},  // want: 0 documents the intent
}
```

### 6.3 Slices — project decision

**PROJECT RULE**: Use `make([]T, 0)` when the slice is returned via JSON or API.

Rationale: a `nil` slice serializes as `null` in JSON; `make([]T, 0)` serializes as `[]`.

```go
// ✅ — in repositories and use cases that return via API
pokemons := make([]domain.Pokemon, 0)
favorites := make([]string, 0)

// ✅ — with known capacity (preferable)
pokemons := make([]domain.Pokemon, 0, len(ids))

// ✅ — var is ok for internal slices not returned via JSON
var nums []int
if condition {
    nums = append(nums, 1)
}

// ✅ — check empty with len(), not == nil
if len(pokemons) == 0 { ... }
```

### 6.4 Maps

```go
// ✅ — empty map
m := make(map[string]*domain.Pokemon)

// ✅ — with known capacity
m := make(map[string]*domain.Pokemon, len(ids))

// ✅ — map with fixed set of elements
statusMessages := map[int]string{
    http.StatusOK:           "ok",
    http.StatusNotFound:     "not found",
    http.StatusUnauthorized: "unauthorized",
}

// ❌ — nil map (panic on write)
var m map[string]*domain.Pokemon
m["pikachu"] = p // panic!
```

### 6.5 Struct references: use `&T{}` not `new(T)`

```go
// ✅
svc := &PokemonService{repo: repo}

// ❌ — inconsistent
svc := new(PokemonService)
svc.repo = repo
```

---

## 7. Interfaces and compile-time assertions

### 7.1 Define on the consumer side

```go
// ports/outbound/repository.go — consumer defines the interface
type PokemonRepository interface {
    GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
    ListAll(ctx context.Context) ([]domain.Pokemon, error)
}

// adapters/outbound/postgres/pokemon_repository.go — implements
type PostgresPokemonRepository struct { db *sql.DB }
```

### 7.2 Small interfaces

Prefer 1–3 methods. Large interfaces are hard to mock and violate ISP:

```go
// ✅ — focused
type TokenValidator interface {
    ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error)
}

// ❌ — too broad
type AuthEverything interface {
    Login(...)
    Logout(...)
    Register(...)
    ValidateToken(...)
    RefreshToken(...)
    ChangePassword(...)
}
```

### 7.3 Compile-time interface assertion — REQUIRED in the project

For each struct that implements an interface, declare the assertion at the **end of the file**:

```go
// At the end of each implementation file:

// adapters/outbound/postgres/pokemon_repository.go
var _ outbound.PokemonRepository = (*PostgresPokemonRepository)(nil)

// adapters/outbound/postgres/favorite_repository.go
var _ outbound.FavoriteRepository = (*PostgresFavoriteRepository)(nil)

// adapters/outbound/http/auth_service_client.go (implements 2 interfaces)
var _ outbound.AuthProvider    = (*AuthServiceClient)(nil)
var _ inbound.TokenValidator   = (*AuthServiceClient)(nil)

// adapters/outbound/http/pokemon_catalog_client.go
var _ outbound.PokemonRepository = (*PokemonCatalogServiceRepository)(nil)

// service/pokemon_service.go
var _ inbound.PokemonUseCase = (*PokemonService)(nil)

// service/favorite_service.go
var _ inbound.FavoriteUseCase = (*FavoriteService)(nil)
```

---

## 8. Errors

### 8.1 Never silently discard errors

```go
// ❌ — forbidden in production code
_ = repo.Save(ctx, pokemon)

// ✅ — at least log if you cannot return
if err := repo.Save(ctx, pokemon); err != nil {
    log.Printf("failed to save pokemon %s: %v", id, err)
}
```

### 8.2 Error strings

- Lowercase, no trailing period (they will be concatenated)
- Log/UI messages may have normal capitalization

```go
// ✅
return fmt.Errorf("pokemon not found: %w", err)
return errors.New("invalid token")

// ❌
return fmt.Errorf("Pokemon not found.")
return errors.New("Invalid token.")
```

### 8.3 Wrapping: `%w` vs `%v`

| Use | When |
|-----|------|
| `%w` | The caller can or will use `errors.Is` / `errors.As` |
| `%v` | Log/annotation only, no inspection needed |

```go
// ✅ — return to caller
if err := s.repo.GetByID(ctx, id); err != nil {
    return nil, fmt.Errorf("fetch pokemon %s: %w", id, err)
}

// ✅ — log only
log.Printf("operation completed with warning: %v", err)
```

### 8.4 Sentinel errors

Use `var ErrX = errors.New(...)` for expected errors that callers check:

```go
var (
    ErrNotFound      = errors.New("not found")
    ErrInvalidToken  = errors.New("invalid token")
    ErrAlreadyExists = errors.New("already exists")
    ErrUnauthorized  = errors.New("unauthorized")
)
```

### 8.5 Custom error types

Use when you need to carry structured data:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("field %s: %s", e.Field, e.Message)
}
```

### 8.6 How to handle errors (4 options)

1. **Handle and return** — the most common
2. **Wrap and return** with `%w` to the caller
3. **Log and degrade** gracefully (non-critical operation)
4. **Match and degrade** using `errors.Is` / `errors.As`

```go
// 1. Handle and return
if err := validate(input); err != nil {
    return nil, fmt.Errorf("validate input: %w", err)
}

// 3. Log and degrade
if err := emitMetrics(); err != nil {
    // metrics should not bring down the application
    log.Printf("failed to emit metrics: %v", err)
}

// 4. Match and degrade
tz, err := getUserTimeZone(id)
if err != nil {
    if errors.Is(err, ErrNotFound) {
        tz = time.UTC
    } else {
        return fmt.Errorf("get timezone for user %s: %w", id, err)
    }
}
```

### 8.7 Type assertion: use "comma ok"

```go
// ✅
t, ok := i.(string)
if !ok {
    return fmt.Errorf("unexpected type: %T", i)
}

// ❌ — panic if the type is wrong
t := i.(string)
```

### 8.8 In-band errors: avoid

```go
// ❌ — forces the caller to check a sentinel value
func Lookup(key string) int  // returns -1 if not found

// ✅ — multiple return values
func Lookup(key string) (value string, ok bool)
func GetByID(ctx context.Context, id string) (*Pokemon, error)
```

---

## 9. Context

### 9.1 First parameter, always

```go
// ✅
func (s *PokemonService) ListPokemons(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)

// ❌
func (s *PokemonService) ListPokemons(page, pageSize int, ctx context.Context) (*domain.PokemonPage, error)
```

### 9.2 Never store context in a struct

```go
// ❌ — forbidden
type Service struct {
    ctx context.Context
}

// ✅ — pass as parameter in each call
func (s *Service) Do(ctx context.Context) error { ... }
```

### 9.3 Use context for cancellation in goroutines

```go
func (s *Scheduler) Start(ctx context.Context) {
    // terminates when ctx is cancelled by the caller
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case <-s.ticker.C:
                s.processWork(ctx)
            }
        }
    }()
}
```

---

## 10. Goroutines and concurrency

### 10.1 Never launch goroutines without a termination strategy

```go
// ❌ — orphan goroutine
func init() {
    go backgroundWorker()
}

// ✅ — lifecycle documented, context used for cancellation
func NewScheduler(ctx context.Context) *Scheduler {
    s := &Scheduler{}
    // goroutine terminates when ctx is cancelled
    go s.run(ctx)
    return s
}
```

### 10.2 Document when the goroutine terminates

```go
// Start starts the background worker.
// The goroutine terminates when ctx is cancelled. Wait with WaitGroup if needed.
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
    wg.Add(1)
    go func() {
        defer wg.Done()
        w.loop(ctx)
    }()
}
```

### 10.3 Explicit channel size

```go
// ✅ — buffer justified with comment
// buffered channel of 1 to decouple producer and consumer
results := make(chan *domain.Pokemon, 1)

// ❌ — no buffer when the consumer may be slow → deadlock
results := make(chan *domain.Pokemon)
```

### 10.4 Atomic operations: use `go.uber.org/atomic`

```go
// ❌ — easy to forget to use atomic operation
type Service struct {
    running int32 // atomic
}

// ✅ — type-safe
type Service struct {
    running atomic.Bool
}

func (s *Service) start() {
    if s.running.Swap(true) {
        return // already running
    }
}
```

---

## 11. Mutex and synchronization

### 11.1 Declare mutex next to what it protects

```go
// ✅
type MockPokemonRepository struct {
    mu       sync.RWMutex
    pokemons map[string]*domain.Pokemon // protected by mu
}
```

### 11.2 RWMutex when reads outweigh writes

```go
func (m *MockPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    p, ok := m.pokemons[id]
    if !ok {
        return nil, domain.ErrNotFound
    }
    return p, nil
}

func (m *MockPokemonRepository) Save(ctx context.Context, p *domain.Pokemon) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.pokemons[p.ID] = p
    return nil
}
```

### 11.3 Avoid embedding types in public structs

```go
// ❌ — leaks implementation details, prevents evolution
type ConcreteList struct {
    *AbstractList
}

// ✅ — explicit composition
type ConcreteList struct {
    list *AbstractList
}

func (l *ConcreteList) Add(e Entity) { l.list.Add(e) }
```

---

## 12. Performance

### 12.1 Prefer `strconv` over `fmt` for numeric conversions

```go
// ✅ — faster, 1 alloc
s := strconv.Itoa(42)
n, _ := strconv.Atoi("42")

// ❌ — slower, 2 allocs
s := fmt.Sprintf("%d", 42)
```

### 12.2 Avoid repeated string→byte conversions

```go
// ✅
separator := []byte(", ")  // convert once
for _, item := range items {
    w.Write(separator)
    w.Write([]byte(item))
}

// ❌ — recreates the slice on every iteration
for range items {
    w.Write([]byte(", "))
}
```

### 12.3 Pre-allocate slices and maps with known capacity

```go
// ✅
pokemons := make([]domain.Pokemon, 0, len(ids))
index := make(map[string]*domain.Pokemon, len(ids))

// ❌ — reallocates as it grows
var pokemons []domain.Pokemon
index := map[string]*domain.Pokemon{}
```

### 12.4 `time.Duration` for time constants — project decision

**PROJECT RULE**: Time duration constants and variables must use `time.Duration`.

```go
// ✅
const defaultAuthRateLimitWindow = 60 * time.Second

// direct use — no cast needed
time.NewTicker(defaultAuthRateLimitWindow)

// ❌ — int type, requires cast every use
const defaultAuthRateLimitWindowSeconds = 60

// bad usage — repeated cast
time.NewTicker(time.Duration(defaultAuthRateLimitWindowSeconds) * time.Second)
```

### 12.5 Format string as constant

```go
// ✅ — go vet can verify at compile time
const msgFmt = "pokemon %s not found (status=%d)"
fmt.Printf(msgFmt, id, code)

// ❌ — go vet cannot reach
msg := "pokemon %s not found (status=%d)"
fmt.Printf(msg, id, code)
```

---

## 13. Control flow

### 13.1 Indent error flow — return early

Return when an error is found; the happy path stays unindented.

```go
// ✅
func (s *PokemonService) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id == "" {
        return nil, domain.ErrInvalidInput
    }

    pokemon, err := s.pokemonRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("fetch pokemon %s: %w", id, err)
    }

    return pokemon, nil
}

// ❌ — happy path nested in else
func (s *PokemonService) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id != "" {
        pokemon, err := s.pokemonRepo.GetByID(ctx, id)
        if err == nil {
            return pokemon, nil
        } else {
            return nil, fmt.Errorf("fetch pokemon %s: %w", id, err)
        }
    } else {
        return nil, domain.ErrInvalidInput
    }
}
```

### 13.2 Switch vs if-else

Prefer `switch` when there are 3+ branches on the same variable:

```go
// ✅
switch resp.StatusCode {
case http.StatusOK:
    return parseResponse(body)
case http.StatusNotFound:
    return nil, domain.ErrNotFound
case http.StatusUnauthorized:
    return nil, domain.ErrInvalidToken
case http.StatusConflict:
    return nil, domain.ErrAlreadyExists
default:
    return nil, fmt.Errorf("unexpected status from auth-service: %d", resp.StatusCode)
}
```

### 13.3 Reduce nesting

```go
// ✅ — early return reduces nesting
for _, item := range items {
    if item == nil {
        continue
    }
    if !item.IsValid() {
        continue
    }
    process(item)
}

// ❌ — pyramid of doom
for _, item := range items {
    if item != nil {
        if item.IsValid() {
            process(item)
        }
    }
}
```

### 13.4 Raw string literals to avoid escapes

```go
// ✅ — readable
wantError := `unknown name:"test"`

// ❌ — hard to read
wantError := "unknown name:\"test\""
```

---

## 14. Table-driven tests

### 14.1 When to use table-driven

Use when there are multiple cases that vary only in inputs/outputs with the same verification logic:

```go
func TestAuthServiceClientErrorMapping(t *testing.T) {
    tests := []struct {
        name      string
        method    func(client *httpclient.AuthServiceClient) error
        transport roundTripFunc
        wantErr   error
    }{
        {
            name: "login returns ErrInvalidCredentials for 401",
            method: func(c *httpclient.AuthServiceClient) error {
                _, err := c.Login(context.Background(), "user@test.com", "wrongpass")
                return err
            },
            transport: mockStatusCode(http.StatusUnauthorized),
            wantErr:   domain.ErrInvalidCredentials,
        },
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := httpclient.NewAuthServiceClient(
                "http://test",
                httpclient.WithHTTPClient(&http.Client{Transport: tt.transport}),
            )
            err := tt.method(client)
            require.ErrorIs(t, err, tt.wantErr)
        })
    }
}
```

### 14.2 When NOT to use table-driven

- Cases with very different setup/teardown between them
- Cases with complex, unique assertions
- Happy path tests with many field checks (use explicit test)

```go
// ✅ — explicit test for happy path with multiple assertions
func TestAuthServiceClientSignupSuccess(t *testing.T) {
    // specific setup
    client := setupClientWithMockResponse(t, validSignupResponse)
    
    session, err := client.Signup(context.Background(), validSignupRequest)
    
    require.NoError(t, err)
    assert.Equal(t, "expected-token", session.AccessToken)
    assert.Equal(t, "user@test.com", session.UserEmail)
    assert.False(t, session.ExpiresAt.IsZero())
}
```

### 14.3 Subtest naming

- Use `t.Run(tt.name, ...)` with names describing the expected behavior
- Pattern: `"method returns ErrorX for condition Y"`

### 14.4 Unnecessary fields in the test struct

Never declare fields in the test struct that are not used in the `t.Run` body:

```go
// ❌ — wantPath declared but never used
tests := []struct {
    name      string
    wantErr   error
    wantPath  string  // declared but never used
}{ ... }

// ✅ — only what is verified
tests := []struct {
    name    string
    wantErr error
}{ ... }
```

---

## 15. Functional Options

Use when a constructor has many optional parameters:

```go
type ClientOption func(*AuthServiceClient)

func WithTimeout(d time.Duration) ClientOption {
    return func(c *AuthServiceClient) {
        c.httpClient.Timeout = d
    }
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
    return func(c *AuthServiceClient) {
        c.httpClient = httpClient
    }
}

func NewAuthServiceClient(baseURL string, opts ...ClientOption) *AuthServiceClient {
    c := &AuthServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

// Usage
client := NewAuthServiceClient(
    "http://auth-service:8080",
    WithTimeout(30 * time.Second),
)

// In tests
client := NewAuthServiceClient(
    "http://test",
    WithHTTPClient(&http.Client{Transport: mockTransport}),
)
```

---

## 16. Panic and init()

### 16.1 Do not use panic in production code

```go
// ❌
func run(args []string) {
    if len(args) == 0 {
        panic("required argument")
    }
}

// ✅ — return error and let the caller decide
func run(args []string) error {
    if len(args) == 0 {
        return errors.New("required argument")
    }
    return nil
}

func main() {
    if err := run(os.Args[1:]); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

Exception: program initialization with `template.Must`, `regexp.MustCompile`:

```go
// ✅ — panic only at initialization, not during requests
var statusTemplate = template.Must(template.New("status").Parse(statusHTML))
```

In tests, use `t.Fatal` / `t.FailNow`, never panic:

```go
// ✅
f, err := os.CreateTemp("", "test")
if err != nil {
    t.Fatal("failed to set up test:", err)
}

// ❌
if err != nil {
    panic("failed to set up test")
}
```

### 16.2 Avoid init()

`init()` runs too early, makes testing hard, and hides dependencies.

```go
// ❌
func init() {
    go backgroundWorker()    // uncontrolled goroutine
    db = openDatabase()      // hidden side effect
}

// ✅ — explicit initialization in the constructor
func New(ctx context.Context, dbURL string) (*Service, error) {
    db, err := openDatabase(dbURL)
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }
    svc := &Service{db: db}
    go svc.run(ctx) // goroutine with controlled lifecycle
    return svc, nil
}
```

The only acceptable use of `init()`: registering drivers/codecs that cannot be done any other way.

---

## 17. Mutable global variables

Avoid — use dependency injection:

```go
// ❌ — global mutation makes parallel tests difficult
var _timeNow = time.Now

func sign(msg string) string {
    return signWithTime(msg, _timeNow())
}

// ✅ — injectable
type Signer struct {
    now func() time.Time
}

func NewSigner() *Signer {
    return &Signer{now: time.Now}
}

func (s *Signer) Sign(msg string) string {
    return signWithTime(msg, s.now())
}

// In tests:
func TestSign(t *testing.T) {
    s := NewSigner()
    s.now = func() time.Time { return fixedTime }
    assert.Equal(t, want, s.Sign(give))
}
```

### Canonical names for unexported global variables

Use `_` prefix to distinguish globals from locals:

```go
// ✅
var (
    _defaultClient = &http.Client{}
    _logger        = log.New(os.Stderr, "", 0)
)
```

---

## 18. Linting and tools

### 18.1 Essential tools

| Tool | Purpose |
|------|---------|
| `golangci-lint` | Aggregated lint — single runner |
| `go vet` | Basic compiler analysis |
| `gofmt` / `goimports` | Formatting + automatic imports |
| `errcheck` | Detects discarded errors |
| `staticcheck` | Advanced static analysis |
| `revive` | Additional style rules |

### 18.2 Soft line length limit

Soft limit: **99 characters** (Uber). Break before, but it's not strict.  
Exception: long URLs in comments don't need to be wrapped.

### 18.3 Custom Printf functions

If creating a Printf wrapper, end with `f`:

```go
// ✅ — go vet detects automatically
func Wrapf(format string, args ...interface{}) error { ... }

// ❌
func Wrap(format string, args ...interface{}) error { ... }
```

---

## 19. Project-specific decisions

These are decisions already made in the Pokedex Platform that diverge from or complement the guides:

### 19.1 Imports: 3 groups (not 2 like Uber, not 4 like Google)

See [section 4.1](#41-three-groups-project-decision).

### 19.2 Slices returned via API: always `make([]T, 0)`

See [section 6.3](#63-slices--project-decision).

### 19.3 Interface compliance: `var _ I = (*T)(nil)` at end of file

See [section 7.3](#73-compile-time-interface-assertion--required-in-the-project).

### 19.4 Time constants: always `time.Duration`

See [section 12.4](#124-timeduration-for-time-constants--project-decision).

### 19.5 Comments: always Portuguese (Brazil)

See [section 3.1](#31-project-rules-pt-br).

### 19.6 Hexagonal structure: dependency direction

```
adapters/inbound/http  →  ports/inbound (use cases)  ←  service
                                                              ↓
                          ports/outbound  ←  adapters/outbound/{postgres,http}
```

- HTTP handlers do **not** depend on concrete clients
- Errors normalized in the outbound adapter or application layer
- New domain entities go in `domain/`, never in `ports/`

### 19.7 gRPC exported: `GRPC` not `Grpc`

Follow Google Style Guide:

```go
// ✅
type GRPCClient struct{}
func NewGRPCHandler() *GRPCHandler {}

// ❌
type GrpcClient struct{}
func NewGrpcHandler() *GrpcHandler {}
```

### 19.8 Official service names

- `mobile-bff` — BFF focused on client experience
- `pokemon-catalog-service` — canonical source of the catalog
- `auth-service` — authentication and token lifecycle

Never use `pokedex-service` (legacy).

---

## References

| Document | Link | Normative | Canonical |
|----------|------|:---------:|:---------:|
| Uber Go Style Guide (PT-BR) | https://github.com/alcir-junior-caju/uber-go-style-guide-pt-br | — | — |
| Google Go Style Guide | https://google.github.io/styleguide/go/guide | ✅ | ✅ |
| Google Go Style Decisions | https://google.github.io/styleguide/go/decisions | ✅ | ❌ |
| Google Go Best Practices | https://google.github.io/styleguide/go/best-practices | ❌ | ❌ |
| Google Go Style (overview) | https://google.github.io/styleguide/go | — | — |
| Effective Go | https://go.dev/doc/effective_go | — | — |
| Go Code Review Comments | https://github.com/golang/go/wiki/CodeReviewComments | — | — |
