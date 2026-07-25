# Fase 02 — Controle & Coleções

**Aulas do curso:** 27 a 43

## Sumário

| Passo | Recurso | Tempo est. |
|-------|---------|-----------|
| 1 | Assistir aulas do curso | ~2h |
| 2 | Executar exemplos Go by Example | ~30min |
| 3 | Ler Effective Go | ~25min |
| 4 | Conferir roadmap.sh | ~5min |
| 5 | Ler styleguide | ~15min |
| 6 | Fazer exercício prático | ~20min |

---

## 1. Aulas do curso

| # | Aula | Conceito |
|---|------|----------|
| 27 | If/Else | `if`, `else`, blocos |
| 28 | If/Else If | Múltiplas condições |
| 29 | If com Init | `if err := fn(); err != nil { }` |
| 30 | Laço For | Único loop de Go: `for init; cond; post { }` |
| 31 | Switch #01 | `switch`, `case`, `default` |
| 32 | Switch #02 | Múltiplos valores por `case` |
| 33 | Resposta do Desafio Switch | — |
| 34 | Switch #03 | `switch` sem expressão (avalia `true`) |
| 35 | Trabalhando com Arrays | Tamanho fixo, parte do tipo |
| 36 | Percorrendo Arrays com For (Range) | `for i, v := range arr` |
| 37 | Conhecendo o Slice | Tamanho dinâmico, referência a array interno |
| 38 | Construindo Slices com Make | `make([]T, len, cap)` |
| 39 | Usando Mesmo Array Interno | Slices compartilham backing array |
| 40 | Slice: Usando Append e Copy | `append()`, `copy()` |
| 41 | Trabalhando com Maps #01 | `map[K]V`, `make()`, atribuição |
| 42 | Trabalhando com Maps #02 | `delete()`, "comma ok" idiom |
| 43 | Maps Aninhados | `map[string]map[string]int` |

---

## 2. Go by Example — exemplos desta fase

| Exemplo | Link | O que observar |
|---------|------|---------------|
| For | https://gobyexample.com/for | `for` como `while`, loop infinito, `continue`, `break` |
| If/Else | https://gobyexample.com/if-else | `if` com init, sem parênteses, chaves obrigatórias |
| Switch | https://gobyexample.com/switch | Sem `break` automático, múltiplos valores, sem expressão |
| Arrays | https://gobyexample.com/arrays | Tamanho fixo, `[5]int` é diferente de `[6]int` |
| Slices | https://gobyexample.com/slices | `make`, `append`, `copy`, slice operator `[lo:hi]` |
| Maps | https://gobyexample.com/maps | `make`, delete, blank identifier `_` |
| Range | https://gobyexample.com/range-over-built-in-types | `range` com índice, valor, chave; `_` para descartar |
| Functions | https://gobyexample.com/functions | Parâmetros tipados, retorno, chamada |
| Multiple Return Values | https://gobyexample.com/multiple-return-values | `(int, error)` — o padrão Go |
| Variadic Functions | https://gobyexample.com/variadic-functions | `func(nums ...int)` — número variável de argumentos |
| Closures | https://gobyexample.com/closures | Funções anônimas que capturam variáveis do escopo |
| Recursion | https://gobyexample.com/recursion | Funções que chamam a si mesmas |
| String Functions | https://gobyexample.com/string-functions | `strings.Contains`, `HasPrefix`, `Join`, `Replace` |

---

## 3. Effective Go — seções para ler

| Seção | Link | O que aprender |
|-------|------|---------------|
| **Control structures: If** | https://go.dev/doc/effective_go#if | `if` com init, omissão de `else` desnecessário |
| **Redeclaration and reassignment** | https://go.dev/doc/effective_go#redeclaration | Por que `err` pode ser redeclarado com `:=` |
| **Control structures: For** | https://go.dev/doc/effective_go#for | `for` é o único loop; o blank identifier `_` no range |
| **Control structures: Switch** | https://go.dev/doc/effective_go#switch | Switch sem expressão, múltiplos valores, `break` com label |
| **Data: Allocation with new** | https://go.dev/doc/effective_go#allocation_new | `new(T)` aloca e zera, retorna `*T` |
| **Data: Allocation with make** | https://go.dev/doc/effective_go#allocation_make | `make` para slices, maps e channels apenas |
| **Data: Arrays** | https://go.dev/doc/effective_go#arrays | Arrays são valores (copiados), use slices |
| **Data: Slices** | https://go.dev/doc/effective_go#slices | Slices são referências a arrays; `append` retorna novo slice |
| **Data: Maps** | https://go.dev/doc/effective_go#maps | "comma ok" idiom, `delete`, zero value em lookup |
| **Printing** | https://go.dev/doc/effective_go#printing | `%v`, `%+v`, `%#v`, `%T` — formatação poderosa |
| **Append** | https://go.dev/doc/effective_go#append | `append(slice, elem...)` — built-in essencial |

---

## 4. Roadmap.sh — tópicos desta fase

| Categoria | Tópicos |
|-----------|---------|
| Conditionals | `if`, `if-else`, `switch` |
| Loops | `for` loop, `for range`, iterating maps, iterating strings, `break`, `continue` |
| Composite Types | Arrays, Slices (Capacity and Growth, `make()`), Maps (Comma-Ok Idiom) |
| Functions | Basics, Variadic, Multiple Return Values, Anonymous, Closures |

---

## 5. Guia de Estilo — regras desta fase

> Leia estas seções do [`decisions.md`](../guia-estilo/decisions.md):

| Regra | Seção | Por que importa agora |
|-------|-------|----------------------|
| Reduzir aninhamento | [Reduzir aninhamento](../guia-estilo/decisions.md#reduzir-aninhamento) | Trate o erro primeiro, `continue` ou `return` cedo. Evite `if` dentro de `if`. |
| Else desnecessário | [Else desnecessário](../guia-estilo/decisions.md#else-desnecessário) | Se o `if` termina com `return`, não use `else`. |
| Reduzir escopo de variáveis | [Reduzir o escopo de variáveis](../guia-estilo/decisions.md#reduzir-o-escopo-de-variáveis) | Declare a variável o mais próximo possível do uso. |
| `nil` é um slice válido | [Slices](../guia-estilo/decisions.md#slices) | Use `nil` para slices vazios. Exceção: `make([]T, 0)` ao retornar via API/JSON. |
| Capacidade de containers | [Especificar capacidade](../guia-estilo/best-practices.md#especificar-capacidade-de-containers) | `make([]T, 0, len(data))` quando souber o tamanho. |
| Evitar conversões repetidas `string`↔`[]byte` | [Desempenho](../guia-estilo/best-practices.md#evitar-conversões-repetidas-de-string-para-bytes) | Converta uma vez, reutilize. Use `io.WriteString`. |
| `strconv` > `fmt` para conversões | [Preferir strconv](../guia-estilo/best-practices.md#preferir-strconv-a-fmt) | `strconv.Itoa(n)` em vez de `fmt.Sprintf("%d", n)` |

---

## 6. No código do projeto

### Indent error flow — o "happy path" sem aninhamento

```go
// Em internal/service/pokemon_service.go — padrão real do projeto:
func (s *PokemonService) GetByID(ctx context.Context, id string) (*domain.PokemonDetail, error) {
    pokemon, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("buscar pokemon: %w", err)
    }
    return pokemon, nil
}
```

### Capacidade de slice com `make`

```go
// Em internal/adapters/inbound/http/pokemon_handler.go:
items := make([]dto.PokemonDTO, 0, len(pokemons))
for _, p := range pokemons {
    items = append(items, dto.ToPokemonDTO(p))
}
// capacidade pré-alocada evita realocações no append
```

### "Comma ok" idiom com type assertion

```go
// Em internal/adapters/inbound/http/middleware.go:
userID, ok := r.Context().Value(userIDKey).(string)
if !ok {
    http.Error(w, "nao autorizado", http.StatusUnauthorized)
    return
}
```

---

## 7. Exercício prático

**Objetivo:** Refatorar código com `if` aninhado para indent error flow.

1. Abra `internal/service/favorite_service.go`
2. Encontre uma função com estrutura semelhante a:

```go
// Ruim — aninhamento desnecessário:
if user != nil {
    if fav, err := repo.Find(ctx, user.ID); err == nil {
        return fav, nil
    } else {
        return nil, err
    }
} else {
    return nil, errors.New("usuário não autenticado")
}
```

3. Refatore para:

```go
// Bom — indent error flow:
if user == nil {
    return nil, errors.New("usuário não autenticado")
}
fav, err := repo.Find(ctx, user.ID)
if err != nil {
    return nil, err
}
return fav, nil
```

4. Verifique se o código existente já segue esse padrão. Se encontrar violações, anote para uma futura correção.

---

**Próxima fase:** [Fase 03 — Structs, Interfaces & Erros](fase-03-structs-interfaces.md)
