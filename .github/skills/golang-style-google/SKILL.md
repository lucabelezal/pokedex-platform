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

> Baseado em: Google Go Style Guide, Google Go Style Decisions e Google Go Best Practices

## Quando usar esta skill

Carregue esta skill ao:
- Revisar ou escrever código Go que precisa seguir o estilo canônico do Google
- Fazer code review com foco em clareza e idiomaticidade
- Resolver dúvidas sobre naming, organização de packages, comentários ou erros
- Decidir entre várias formas idiomáticas de expressar a mesma coisa em Go

---

## 1. Naming

### Pacotes

- Nome do pacote: **singular, minúsculo, sem underscore**. Ex.: `user`, `pokemon`, `auth`.
- Evite nomes genéricos: `util`, `common`, `helpers` — dificultam descoberta e causam conflitos.
- O nome do pacote é parte da API: `auth.Service`, `pokemon.Repository`. Não repita o pacote no nome do símbolo.
  - ✅ `auth.NewService()` — em vez de `auth.NewAuthService()`
  - ✅ `pokemon.Repository` — em vez de `pokemon.PokemonRepository`

### Funções e métodos

- Nomes curtos e descritivos: o escopo de uso determina o comprimento apropriado.
- **Sem prefixo `Get`** para getters simples: `Name()` não `GetName()`, `Count()` não `GetCount()`.
- Construtores: `New` + tipo. Ex.: `NewPokemonService`, `NewHandler`.
- Predicados booleanos: `Is`, `Has`, `Can`, `Should`. Ex.: `IsFavorite`, `HasAccess`.

### Variáveis e campos

- Nomes de variáveis de loop curtos são aceitáveis: `i`, `j`, `k`, `v`.
- **Initialisms** sempre em caixa completa: `ID`, `URL`, `HTTP`, `JSON`, `gRPC`, `DB`, `API`.
  - ✅ `userID`, `baseURL`, `httpClient`, `jsonData`
  - ❌ `userId`, `baseUrl`, `httpClient` (incorreto apenas se `http` for initialism isolado)
- Receivers: 1-2 letras, abreviação do tipo. Nunca `this`, `self`, `me`.
  ```go
  // Correto
  func (s *PokemonService) ListPokemons(...) {}
  func (c *AuthServiceClient) Login(...) {}

  // Incorreto
  func (this *PokemonService) ListPokemons(...) {}
  func (self *AuthServiceClient) Login(...) {}
  ```

### Erros

- Variáveis de erro como sentinel: prefixo `Err`. Ex.: `ErrNotFound`, `ErrInvalidToken`.
- Tipos de erro customizados: sufixo `Error`. Ex.: `ValidationError`, `NotFoundError`.
- Strings de erro: **minúsculas, sem ponto final**.
  - ✅ `"pokemon não encontrado"`
  - ❌ `"Pokemon não encontrado."` — será concatenada com contexto adicional

---

## 2. Comentários

### Godoc

Todo símbolo exportado deve ter comentário godoc começando pelo nome:

```go
// PokemonService gerencia operações de listagem e busca de Pokémons.
type PokemonService struct { ... }

// ListPokemons retorna uma página de Pokémons com informações de favoritos.
// Retorna ErrInvalidPage se page for negativo.
func (s *PokemonService) ListPokemons(ctx context.Context, page, pageSize int, userID string) (*domain.PokemonPage, error) { ... }
```

### Frases completas, primeira letra maiúscula

```go
// Correto: frase completa começando com o nome do símbolo
// ParseConfig lê a configuração do ambiente.

// Incorreto: fragmento
// lê a configuração do ambiente
```

### Não documente o óbvio

```go
// Ruim: repete o código
// i é o índice do loop
for i := range items { ... }

// Bom: explica o porquê
// reverse percorre de trás para frente para evitar cópia do slice
for i := len(items) - 1; i >= 0; i-- { ... }
```

---

## 3. Estrutura de Controle

### Indent error flow

Retorne cedo ao encontrar um erro; o caminho normal não deve estar dentro de um `else`:

```go
// ✅ Correto: caminho feliz sem aninhamento
func (s *Service) GetPokemon(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id == "" {
        return nil, domain.ErrInvalidInput
    }

    pokemon, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("buscar pokemon %s: %w", id, err)
    }

    return pokemon, nil
}

// ❌ Incorreto: caminho feliz aninhado no else
func (s *Service) GetPokemon(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id != "" {
        pokemon, err := s.repo.GetByID(ctx, id)
        if err == nil {
            return pokemon, nil
        } else {
            return nil, fmt.Errorf("buscar pokemon %s: %w", id, err)
        }
    } else {
        return nil, domain.ErrInvalidInput
    }
}
```

### Switch vs if-else

Prefira `switch` quando houver 3+ ramos sobre a mesma variável ou expressão:

```go
// ✅ Correto
switch resp.StatusCode {
case http.StatusOK:
    return parseOK(body)
case http.StatusNotFound:
    return nil, domain.ErrNotFound
case http.StatusUnauthorized:
    return nil, domain.ErrInvalidToken
default:
    return nil, fmt.Errorf("status inesperado: %d", resp.StatusCode)
}
```

---

## 4. Declaração de Variáveis

Use o estilo que melhor comunica a intenção:

```go
// var x T — zero value explícito e intencional
var count int
var repo PokemonRepository

// x := value — inicializa com valor não-zero
name := "Pikachu"
client := &http.Client{Timeout: 5 * time.Second}

// x := T{} — quando precisa do tipo visível no lado esquerdo
handler := Handler{
    pokemonUseCase: svc,
}
```

---

## 5. Wrapping de Erros

```go
// %w — quando o chamador pode ou deve inspecionar o erro com errors.Is / errors.As
if err := s.repo.Save(ctx, pokemon); err != nil {
    return fmt.Errorf("salvar pokemon: %w", err)
}

// %v — quando é só anotação de contexto, sem necessidade de inspeção
log.Printf("operação concluída com aviso: %v", err)
```

**Regra**: use `%w` por padrão em retornos de erro; `%v` apenas em logs onde a cadeia de erro não importa.

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

- **Documente o ciclo de vida**: quando a goroutine termina e quem é responsável.
- Use `context.Context` para sinalizar cancelamento.
- Nunca lance goroutines sem estratégia de finalização (WaitGroup, canal, context).

```go
// ✅ Correto: ciclo de vida documentado, context usado para cancelamento
func (s *Scheduler) Start(ctx context.Context) {
    // encerra quando ctx for cancelado
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

## 8. Organização do Código

### Ordem dos elementos em um arquivo

1. Declaração do pacote + godoc do pacote (se aplicável)
2. Imports (stdlib / external / internal, separados por linha em branco)
3. Constantes
4. Variáveis de pacote
5. Tipos + construtores
6. Métodos

### Grupos de import

```go
import (
    // stdlib
    "context"
    "fmt"
    "net/http"

    // externos
    "github.com/stretchr/testify/assert"

    // internos
    "pokedex-platform/core/bff/mobile-bff/internal/domain"
    inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
)
```

---

## Fontes

- [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)
- [Google Go Best Practices](https://google.github.io/styleguide/go/best-practices)
