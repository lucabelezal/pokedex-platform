# Fase 03 — Structs, Interfaces & Erros

**Aulas do curso:** 44 a 71

Esta é a fase mais longa e mais importante. Structs, interfaces e tratamento de erros
são o coração do design em Go — e da arquitetura hexagonal da Pokedex Platform.

## Sumário

| Passo | Recurso | Tempo est. |
|-------|---------|-----------|
| 1 | Assistir aulas do curso | ~3h30 |
| 2 | Executar exemplos Go by Example | ~45min |
| 3 | Ler Effective Go | ~40min |
| 4 | Conferir roadmap.sh | ~5min |
| 5 | Ler styleguide | ~20min |
| 6 | Fazer exercício prático | ~25min |

---

## 1. Aulas do curso

| # | Aula | Conceito |
|---|------|----------|
| 44 | Funções Básicas | Revisão aprofundada |
| 45 | Pilha de Funções | Call stack |
| 46 | Retorno Nomeado | `func() (result int, err error)` |
| 47 | Armazenar Funções em Variáveis | Funções como first-class citizens |
| 48 | Passar Função como Parâmetro | Higher-order functions |
| 49 | Funções Variáticas #01 | `func(nums ...int)` |
| 50 | Funções Variáticas #02 | Spread `slice...` |
| 51 | Closure | Funções que capturam escopo externo |
| 52 | Recursividade | Função chama a si mesma |
| 53 | Recursividade (Simples) | Caso base + chamada recursiva |
| 54 | Defer | `defer fn()` — executa ao final da função |
| 55 | Passando Ponteiro para Função | Modificar valor original |
| 56 | Função Init | `func init()` — executada automaticamente |
| 57 | Usando Struct | `type X struct { ... }` |
| 58 | Struct Aninhada | Struct dentro de struct |
| 59 | Métodos em Structs | `func (r Receiver) Method()` |
| 60 | Pseudo-Herança em Structs | Embedding de structs |
| 61 | Tipo Personalizado | `type MyType int` |
| 62 | Usando Interfaces #01 | `type Interface interface { Method() }` |
| 63 | Usando Interfaces #02 | Polimorfismo |
| 64 | Composição de Interfaces | Interface composta de outras interfaces |
| 65 | Tipo Interface | Interface como valor, type assertion |
| 66 | Convertendo uma Struct em JSON | `encoding/json`, tags `json:"name"` |
| 67 | Executando Múltiplos Arquivos | `go run *.go` |
| 68 | Atualização do Go | Nova versão |
| 69 | Pacotes & Visibilidade | Exportado = Maiúscula, Não exportado = minúscula |
| 70 | Criando um Pacote Reutilizável | Estrutura de diretórios e `go.mod` |
| 71 | Criando & Instalando um Pacote do Github | `go get`, `go install` |

---

## 2. Go by Example — exemplos desta fase

| Exemplo | Link | O que observar |
|---------|------|---------------|
| Structs | https://gobyexample.com/structs | Definição, inicialização com nomes de campo, `&T{}` |
| Methods | https://gobyexample.com/methods | `func (r *Rect) Area() float64` — receiver pointer vs value |
| Interfaces | https://gobyexample.com/interfaces | Implementação implícita, `type geometry interface { area() }` |
| Enums | https://gobyexample.com/enums | `iota` para enumerados, `String()` method |
| Struct Embedding | https://gobyexample.com/struct-embedding | Composição sobre herança |
| Generics | https://gobyexample.com/generics | `func Sum[T int|float64](slice []T) T` — Go 1.18+ |
| Errors | https://gobyexample.com/errors | `errors.New()`, retorno `(result, error)`, `if err != nil` |
| Custom Errors | https://gobyexample.com/custom-errors | `type argError struct { ... }` com método `Error()` |
| Defer | https://gobyexample.com/defer | LIFO, argumentos avaliados no momento do `defer` |
| Panic | https://gobyexample.com/panic | `panic()` interrompe execução normal |
| Recover | https://gobyexample.com/recover | `recover()` dentro de `defer` captura panic |
| JSON | https://gobyexample.com/json | `Marshal`, `Unmarshal`, tags `json:"name"` |
| String Formatting | https://gobyexample.com/string-formatting | `%v`, `%+v`, `%#v`, `%T`, `%q`, `%x` |
| Text Templates | https://gobyexample.com/text-templates | `template.Must`, `Execute` |
| Regular Expressions | https://gobyexample.com/regular-expressions | `regexp.MustCompile`, `FindString` |
| Time | https://gobyexample.com/time | `time.Now`, `time.Date`, `Add`, `Sub` |
| Epoch | https://gobyexample.com/epoch | `Unix()`, `UnixMilli()`, `UnixNano()` |
| Time Formatting / Parsing | https://gobyexample.com/time-formatting-parsing | `time.RFC3339`, formato de referência Go |
| Random Numbers | https://gobyexample.com/random-numbers | `rand.Intn`, `rand.Float64` |
| Number Parsing | https://gobyexample.com/number-parsing | `strconv.ParseFloat`, `ParseInt`, `Atoi` |
| URL Parsing | https://gobyexample.com/url-parsing | `url.Parse`, esquema, host, path, query |
| SHA256 Hashes | https://gobyexample.com/sha256-hashes | `sha256.Sum256`, `hex.EncodeToString` |
| Base64 Encoding | https://gobyexample.com/base64-encoding | `base64.StdEncoding`, `URLEncoding` |
| Embed Directive | https://gobyexample.com/embed-directive | `//go:embed` — Go 1.16+ |
| XML | https://gobyexample.com/xml | `xml.MarshalIndent`, tags `xml:"name"` |

---

## 3. Effective Go — seções para ler

| Seção | Link | O que aprender |
|-------|------|---------------|
| **Functions** | https://go.dev/doc/effective_go#functions | Múltiplos retornos, retorno nomeado — o padrão Go para erros |
| **Defer** | https://go.dev/doc/effective_go#defer | LIFO, limpeza de recursos, tracing |
| **Initialization** (Constants, Variables, init) | https://go.dev/doc/effective_go#initialization | `iota` para enums, `init()` e ordem de inicialização |
| **Methods: Pointers vs Values** | https://go.dev/doc/effective_go#pointers_vs_values | Pointer receiver modifica, value receiver não |
| **Interfaces and other types** | https://go.dev/doc/effective_go#interfaces_and_types | Implementação implícita, conversões, type assertions, generality |
| **The blank identifier** | https://go.dev/doc/effective_go#blank | `_` em atribuição múltipla, imports não usados, side-effect imports, interface checks |
| **Embedding** | https://go.dev/doc/effective_go#embedding | Composição de structs e interfaces, promoção de métodos |
| **Errors** | https://go.dev/doc/effective_go#errors | `error` interface, `fmt.Errorf`, type assertion em erros |

---

## 4. Roadmap.sh — tópicos desta fase

| Categoria | Tópicos |
|-----------|---------|
| Structs | Definition, Struct Tags & JSON, Embedding |
| Methods and Interfaces | Methods vs Functions, Pointer Receivers, Value Receivers, Interfaces, Empty Interfaces, Embedding Interfaces, Type Assertions, Type Switch |
| Generics | Why Generics?, Generic Functions, Generic Types/Interfaces, Type Constraints, Type Inference |
| Error Handling | `error` interface, `errors.New`, `fmt.Errorf`, Wrapping/Unwrapping, Sentinel Errors, `panic` and `recover`, Stack Traces |
| Code Organization | Modules & Dependencies (`go mod init`, `go mod tidy`), Packages, Import Rules, Using 3rd Party Packages, Publishing Modules |
| Standard Library | `encoding/json`, `time`, `flag`, `regexp`, `bufio` |

---

## 5. Guia de Estilo — regras desta fase

> Leia estas seções:

| Regra | Seção | Por que importa agora |
|-------|-------|----------------------|
| Receivers: 1-2 letras, nunca `this`/`self` | [Nomes de receivers](../guia-estilo/guide.md#nomes-de-receivers) | `s` para Service, `c` para Client, `h` para Handler |
| Inicialização de structs com nomes de campo | [Inicialização de structs](../guia-estilo/decisions.md#inicialização-de-structs) | `User{Name: "Maria"}` nunca `User{"Maria"}` |
| Omitir campos de valor zero | [Omitir campos](../guia-estilo/decisions.md#omitir-campos-de-valor-zero) | Só declare o que não for zero value |
| `var t T` para zero value struct | [Structs de valor zero](../guia-estilo/decisions.md#structs-de-valor-zero) | `var opts options` em vez de `opts := options{}` |
| Interface compliance: `var _ Interface = (*Impl)(nil)` | [Verificar conformidade](../guia-estilo/decisions.md#verificar-conformidade-de-interface) | **Regra crítica** — vai em TODO adapter e service |
| Ponteiros para interfaces: nunca | [Ponteiros para interfaces](../guia-estilo/decisions.md#ponteiros-para-interfaces) | `func(p PokemonRepository)`, não `func(p *PokemonRepository)` |
| Erros: `%w` para inspecionável, `%v` para obscurecer | [Envoltório de erros](../guia-estilo/decisions.md#envoltório-de-erros-wrapping) | `fmt.Errorf("buscar pokemon: %w", err)` |
| Error strings minúsculas, sem ponto final | [Error strings](../guia-estilo/decisions.md#error-strings) | `errors.New("pokemon nao encontrado")` |
| Tratar erro uma vez (logar OU retornar) | [Tratar erros uma vez](../guia-estilo/decisions.md#tratar-erros-uma-vez) | Nunca `log.Printf(...); return err` |
| Indent error flow | [Indent error flow](../guia-estilo/decisions.md#indent-error-flow) | Trate o erro primeiro, happy path sem indentação |
| Type assertion com "comma ok" | [Falhas de asserção de tipo](../guia-estilo/decisions.md#lidar-com-falhas-de-asserção-de-tipo) | `v, ok := x.(T)` — nunca `v := x.(T)` sem verificação |
| Evitar `init()` | [Evitar init()](../guia-estilo/best-practices.md#evitar-init) | Prefira inicialização explícita |
| `defer` para limpeza | [defer para limpeza](../guia-estilo/best-practices.md#defer-para-limpeza) | `defer f.Close()` logo após `os.Open` |
| Enums com `iota + 1` | [Enums com iota](../guia-estilo/best-practices.md#enums-com-iota) | Comece em 1 para distinguir do zero value |
| Tags de struct consistentes | [Tags de campos](../guia-estilo/best-practices.md#tags-de-campos-em-structs-serializadas) | `json:"id"`, espaços entre tags |
| Evitar nomes embutidos | [Evitar nomes embutidos](../guia-estilo/best-practices.md#evitar-usar-nomes-embutidos) | Não use `error`, `string`, `len` como nome de variável |

---

## 6. No código do projeto

### Interface compliance em todo adapter e service

```go
// Em internal/service/pokemon_service.go (linha ~120):
var _ inbound.PokemonUseCase = (*PokemonService)(nil)

// Em internal/adapters/outbound/http/pokemon_catalog_client.go (linha ~143):
var _ outbound.PokemonRepository = (*PokemonCatalogServiceRepository)(nil)

// Em internal/adapters/outbound/postgres/pokemon_repository.go (linha ~371):
var _ outbound.PokemonRepository = (*PostgresPokemonRepository)(nil)
```

### Erros de domínio como sentinel errors

```go
// Em internal/domain/errors.go:
var (
    ErrPokemonNotFound       = errors.New("pokemon nao encontrado")
    ErrUserAlreadyExists     = errors.New("usuario ja existe")
    ErrFavoriteAlreadyExists = errors.New("favorito ja existe")
    ErrInvalidCredentials    = errors.New("credenciais invalidas")
    ErrInvalidToken          = errors.New("token invalido")
    ErrServiceUnavailable    = errors.New("servico temporariamente indisponivel")
)

// Uso nos handlers — errors.Is para verificar:
if errors.Is(err, domain.ErrPokemonNotFound) {
    RespondError(w, http.StatusNotFound, "pokemon nao encontrado", "NOT_FOUND")
    return
}
```

### Receivers padronizados

```go
// Em internal/service/pokemon_service.go:
func (s *PokemonService) List(ctx context.Context, ...)  // s = Service

// Em internal/service/favorite_service.go:
func (s *FavoriteService) Add(ctx context.Context, ...)   // s = Service

// Em internal/adapters/inbound/http/handler.go:
func (h *Handler) GetPokemonDetails(w http.ResponseWriter, ...)  // h = Handler

// Em internal/adapters/inbound/http/response_builder.go:
func (rb *ResponseBuilder) BuildPokemonResponse(...)  // rb = ResponseBuilder
```

### Arquitetura hexagonal: ports e adapters

```go
// Ports INBOUND (interfaces que os handlers consomem):
// internal/ports/inbound/pokemon_usecase.go
type PokemonUseCase interface {
    List(ctx context.Context, params SearchParams) (*domain.PokemonPage, error)
    GetByID(ctx context.Context, id string) (*domain.PokemonDetail, error)
    Search(ctx context.Context, query string) ([]domain.Pokemon, error)
}

// Ports OUTBOUND (interfaces que os services exigem):
// internal/ports/outbound/pokemon_repository.go
type PokemonRepository interface {
    GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
    List(ctx context.Context, params SearchParams) (*PokemonPage, error)
    Search(ctx context.Context, query string) ([]Pokemon, error)
}
```

---

## 7. Exercício prático

**Objetivo:** Rastrear a implementação de uma interface da porta até o adapter.

1. Comece em `internal/ports/inbound/pokemon_usecase.go` — veja a interface `PokemonUseCase`
2. Vá para `internal/service/pokemon_service.go` — veja como `PokemonService` implementa:
   ```go
   var _ inbound.PokemonUseCase = (*PokemonService)(nil)  // compile-time check
   ```
3. Veja como o service depende de `outbound.PokemonRepository` (interface)
4. Vá para `internal/adapters/outbound/http/pokemon_catalog_client.go` — veja como o client HTTP implementa:
   ```go
   var _ outbound.PokemonRepository = (*PokemonCatalogServiceRepository)(nil)
   ```
5. Desenhe o diagrama no papel:

```
Handler → inbound.PokemonUseCase (interface)
              ↑ implementada por
         PokemonService (struct)
              ↓ depende de
         outbound.PokemonRepository (interface)
              ↑ implementada por
         PokemonCatalogServiceRepository (struct) — chamadas HTTP reais
```

6. **Desafio extra:** Crie um novo adapter outbound que implementa `PokemonRepository` usando um arquivo JSON local como fonte de dados. Inclua o `var _ outbound.PokemonRepository = (*SeuAdapter)(nil)`.

---

**Próxima fase:** [Fase 04 — Concorrência](fase-04-concorrencia.md)
