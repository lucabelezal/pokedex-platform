---
name: golang-style-google
description: >-
  Google Go Style Guide for Go projects. Covers naming (packages, functions, receivers, initialisms),
  godoc comments, indent error flow, variable declarations, error wrapping with %w vs %v,
  interface definition placement, goroutine lifecycle, and import grouping.
  Use when writing or reviewing Go code for style, naming decisions, comment quality, or idiomatic patterns.
  Trigger examples: "naming convention", "receiver name", "godoc comment", "error wrapping",
  "indent error flow", "import grouping", "Go style", "initialism", "interface placement".
  Do NOT use for testing patterns (use golang-testing) or security (use golang-security).
---

# Google Go Style Guide — Skill

> Based on: Google Go Style Guide, Google Go Style Decisions, and Google Go Best Practices

## When to use this skill

Load this skill when:
- Reviewing or writing Go code that needs to follow Google's canonical style
- Doing code review focused on clarity and idiomaticity
- Resolving questions about naming, package organization, comments, or errors
- Deciding between several idiomatic ways to express the same thing in Go

---

## 1. Naming

### Packages

- Package name: **singular, lowercase, no underscore**. Ex.: `user`, `pokemon`, `auth`.
- Avoid generic names: `util`, `common`, `helpers` — they hinder discoverability and cause conflicts.
- The package name is part of the API: `auth.Service`, `pokemon.Repository`. Don’t repeat the package in the symbol name.
  - ✅ `auth.NewService()` — instead of `auth.NewAuthService()`
  - ✅ `pokemon.Repository` — instead of `pokemon.PokemonRepository`

### Functions and methods

- Short, descriptive names: the scope of use determines the appropriate length.
- **No `Get` prefix** for simple getters: `Name()` not `GetName()`, `Count()` not `GetCount()`.
- Constructors: `New` + type. Ex.: `NewPokemonService`, `NewHandler`.
- Boolean predicates: `Is`, `Has`, `Can`, `Should`. Ex.: `IsFavorite`, `HasAccess`.

### Variables and fields

- Short loop variable names are acceptable: `i`, `j`, `k`, `v`.
- **Initialisms** always in full case: `ID`, `URL`, `HTTP`, `JSON`, `gRPC`, `DB`, `API`.
  - ✅ `userID`, `baseURL`, `httpClient`, `jsonData`
  - ❌ `userId`, `baseUrl`, `httpClient` (incorrect only when `http` is an isolated initialism)
- Receivers: 1-2 letters, abbreviation of the type. Never `this`, `self`, `me`.
  ```go
  // Correct
  func (s *PokemonService) ListPokemons(...) {}
  func (c *AuthServiceClient) Login(...) {}

  // Incorrect
  func (this *PokemonService) ListPokemons(...) {}
  func (self *AuthServiceClient) Login(...) {}
  ```

### Errors

- Error sentinel variables: `Err` prefix. Ex.: `ErrNotFound`, `ErrInvalidToken`.
- Custom error types: `Error` suffix. Ex.: `ValidationError`, `NotFoundError`.
- Error strings: **lowercase, no trailing period**.
  - ✅ `"pokemon not found"`
  - ❌ `"Pokemon not found."` — will be concatenated with additional context

---

## 2. Comments

### Godoc

Every exported symbol must have a godoc comment starting with its name:

```go
// PokemonService manages listing and search operations for Pokémons.
type PokemonService struct { ... }

// ListPokemons returns a page of Pokémons with favorites information.
// Returns ErrInvalidPage if page is negative.
func (s *PokemonService) ListPokemons(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error) { ... }
```

### Complete sentences, first letter uppercase

```go
// Correct: full sentence starting with the symbol name
// ParseConfig reads the configuration from the environment.

// Incorrect: fragment
// reads the configuration from the environment
```

### Don't document the obvious

```go
// Bad: repeats the code
// i is the loop index
for i := range items { ... }

// Good: explains the why
// reverse iterates backwards to avoid copying the slice
for i := len(items) - 1; i >= 0; i-- { ... }
```

---

## 3. Control Flow

### Indent error flow

Return early when encountering an error; the happy path should not be inside an `else`:

```go
// ✅ Correct: happy path without nesting
func (s *Service) GetPokemon(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id == "" {
        return nil, domain.ErrInvalidInput
    }

    pokemon, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("fetch pokemon %s: %w", id, err)
    }

    return pokemon, nil
}

// ❌ Incorrect: happy path nested in else
func (s *Service) GetPokemon(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id != "" {
        pokemon, err := s.repo.GetByID(ctx, id)
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

### Switch vs if-else

Prefer `switch` when there are 3+ branches on the same variable or expression:

```go
// ✅ Correct
switch resp.StatusCode {
case http.StatusOK:
    return parseOK(body)
case http.StatusNotFound:
    return nil, domain.ErrNotFound
case http.StatusUnauthorized:
    return nil, domain.ErrInvalidToken
default:
    return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
}
```

---

## 4. Variable Declarations

Use the style that best communicates intention:

```go
// var x T — explicit and intentional zero value
var count int
var repo PokemonRepository

// x := value — initializes with non-zero value
name := "Pikachu"
client := &http.Client{Timeout: 5 * time.Second}

// x := T{} — when the type needs to be visible on the left side
handler := Handler{
    pokemonUseCase: svc,
}
```

---

## 5. Error Wrapping

```go
// %w — when the caller can or should inspect the error with errors.Is / errors.As
if err := s.repo.Save(ctx, pokemon); err != nil {
    return fmt.Errorf("save pokemon: %w", err)
}

// %v — when it's just context annotation, without need for inspection
log.Printf("operation completed with warning: %v", err)
```

**Rule**: use `%w` by default in error returns; `%v` only in logs where the error chain doesn’t matter.

---

## 6. Interfaces

### Defina no lado do consumidor

Interfaces pertencem ao pacote que **usa** o comportamento, não ao pacote que o implementa:

```go
// ports/outbound/repository.go — definida no pacote que usa (serviço de aplicação)
type PokemonRepository interface {
    GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
}

// adapters/repository/postgres_repository.go — implementa a interface
type PostgresPokemonRepository struct { db *sql.DB }
func (r *PostgresPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) { ... }
```

### Interfaces pequenas são melhores

Prefira interfaces de 1-3 métodos. Interfaces grandes são difíceis de mockar e testar.

### Asserção de interface em tempo de compilação

Use para detectar erros de implementação cedo:

```go
// Faz o build falhar se AuthService não implementar inbound.AuthUseCase
var _ inbound.AuthUseCase = (*AuthService)(nil)
```

---

## 7. Goroutines

- **Document the lifecycle**: when the goroutine terminates and who is responsible.
- Use `context.Context` to signal cancellation.
- Never launch goroutines without a termination strategy (WaitGroup, channel, context).

```go
// ✅ Correct: lifecycle documented, context used for cancellation
func (s *Scheduler) Start(ctx context.Context) {
    // terminates when ctx is cancelled
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case <-s.ticker.C:
                s.doWork(ctx)
            }
        }
    }()
}
```

---

## 8. Code Organization

### Order of elements in a file

1. Package declaration + package godoc (if applicable)
2. Imports (stdlib / external / internal, separated by blank lines)
3. Constants
4. Package variables
5. Types + constructors
6. Methods

### Import groups

```go
import (
    // stdlib
    "context"
    "fmt"
    "net/http"

    // external
    "github.com/stretchr/testify/assert"

    // internal
    "pokedex-platform/core/bff/mobile-bff/internal/domain"
    inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
)
```

---

## Sources

- [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)
- [Google Go Best Practices](https://google.github.io/styleguide/go/best-practices)
