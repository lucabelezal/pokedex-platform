# Fase 01 — Fundamentos

**Aulas do curso:** 1 a 26

## Sumário

| Passo | Recurso | Tempo est. |
|-------|---------|-----------|
| 1 | Assistir aulas do curso | ~2h30 |
| 2 | Executar exemplos Go by Example | ~20min |
| 3 | Ler Effective Go | ~15min |
| 4 | Conferir roadmap.sh | ~5min |
| 5 | Ler styleguide | ~15min |
| 6 | Fazer exercício prático | ~15min |

---

## 1. Aulas do curso

| # | Aula | Conceito |
|---|------|----------|
| 1 | Visão Geral do Curso | Introdução |
| 2 | Baixando a Apostila | Material de apoio |
| 3 | Links: Apostila & Repositório | Referências |
| 4 | Go: História e Características | Origem, filosofia |
| 5 | Usando o terminal | Terminal básico |
| 6 | Terminais | Variantes |
| 7 | Aviso importante aos usuários de Windows | Compatibilidade |
| 8 | Windows e Mac | Instalação |
| 9 | Linux | Instalação |
| 10 | Ambiente GO: GOROOT & GOPATH | Configuração do ambiente |
| 11 | Atualização da execução do GO | Versões |
| 12 | Primeiro Programa em Go | `package main`, `func main()` |
| 13 | Explorando os Comandos do Terminal | `go run`, `go build`, `go fmt` |
| 14 | Constantes e Variáveis | `const`, `var`, `:=` |
| 15 | Imprimindo Valores no Console | `fmt.Println`, `fmt.Printf` |
| 16 | Tipos Básicos | `int`, `float64`, `string`, `bool` |
| 17 | Tipos Básicos: Os Zeros | Zero values de cada tipo |
| 18 | Conversão entre Tipos Básicos | `int(x)`, `string(b)` |
| 19 | Funções Básicas | `func nome(params) retorno` |
| 20 | Operadores Aritméticos | `+`, `-`, `*`, `/`, `%` |
| 21 | Operadores de Atribuição | `=`, `+=`, `-=`, `*=`, `/=` |
| 22 | Operadores Relacionais | `==`, `!=`, `<`, `>`, `<=`, `>=` |
| 23 | Operadores Lógicos | `&&`, `||`, `!` |
| 24 | Operadores Unários | `++`, `--` (são statements, não expressões) |
| 25 | Operador Ternário??? | **Não existe em Go.** Use `if/else` |
| 26 | Ponteiros em Go | `*T`, `&v`, desreferenciação |

---

## 2. Go by Example — exemplos desta fase

| Exemplo | Link | O que observar |
|---------|------|---------------|
| Hello World | https://gobyexample.com/hello-world | `package main`, `func main()`, `import "fmt"` |
| Values | https://gobyexample.com/values | Strings, inteiros, floats, booleanos, concatenação |
| Variables | https://gobyexample.com/variables | `var`, `:=`, tipos explícitos e inferidos |
| Constants | https://gobyexample.com/constants | `const`, expressões constantes, `iota` inicial |
| Pointers | https://gobyexample.com/pointers | `&` para endereço, `*` para desreferenciar |
| Strings and Runes | https://gobyexample.com/strings-and-runes | `len()`, indexação, runas, UTF-8 |

---

## 3. Effective Go — seções para ler

| Seção | Link | O que aprender |
|-------|------|---------------|
| **Introduction** | https://go.dev/doc/effective_go#introduction | "To write Go well, understand its idioms" |
| **Formatting** | https://go.dev/doc/effective_go#formatting | `gofmt` resolve tudo. Use. Não discuta com a máquina. |
| **Names** (Package names, Getters, Interface names, MixedCaps) | https://go.dev/doc/effective_go#names | A visibilidade é determinada pela capitalização |
| **Semicolons** | https://go.dev/doc/effective_go#semicolons | Por que a chave `{` não pode ir na linha seguinte |

---

## 4. Roadmap.sh — tópicos desta fase

| Categoria | Tópicos |
|-----------|---------|
| Introduction to Go | Why use Go, History of Go |
| Setting up the Environment | Hello World, `go` command |
| Language Basics | Variables & Constants, `var` vs `:=`, Zero Values, `const` and `iota`, Scope and Shadowing |
| Data Types | Boolean, Numeric Types (Integers, Floats, Complex), Runes, Strings (Raw vs Interpreted), Type Conversion |
| Pointers | Pointers Basics |

---

## 5. Guia de Estilo — regras desta fase

> Leia o [`guide.md`](../guia-estilo/guide.md) completo. Ele é curto (~280 linhas) e estabelece as
> regras fundamentais que se aplicam a **todo código Go** que você escrever.

| Regra | Seção | Por que importa agora |
|-------|-------|----------------------|
| `gofmt` sempre | [Formatação](../guia-estilo/guide.md#gofmt) | `goimports` ao salvar o arquivo. Configure seu editor. |
| MixedCaps, nunca snake_case | [MixedCaps](../guia-estilo/guide.md#mixedcaps) | `maxLength`, não `max_length`. Desde a aula 14 você já aplica isso. |
| Sem limite fixo de linha | [Tamanho de linha](../guia-estilo/guide.md#tamanho-de-linha) | Se a linha está longa, refatore — não a quebre. |
| Nomes de pacotes: minúsculo, singular, curto | [Nomes de pacotes](../guia-estilo/guide.md#nomes-de-pacotes) | `package user`, nunca `package utils` |
| Sem prefixo `Get` | [Nomes de funções](../guia-estilo/guide.md#nomes-de-funções) | `Count()`, não `GetCount()` |
| Initialisms: `ID`, `URL`, `HTTP`, `JSON` | [Initialisms](../guia-estilo/guide.md#initialisms) | `userID`, não `userId`. `URLPony`, não `UrlPony`. |
| `const MaxSize = 100`, não `MAX_SIZE` | [Nomes de constantes](../guia-estilo/guide.md#nomes-de-constantes) | MixedCaps também para constantes |
| Variáveis: nome proporcional ao escopo | [Nomes de variáveis](../guia-estilo/guide.md#nomes-de-variáveis) | `i` no loop, `user` na função, `DefaultTimeout` no pacote |

---

## 6. No código do projeto

Veja como essas regras aparecem no código real da Pokedex Platform:

### `const` com MixedCaps

```go
// Em internal/config/config.go — constantes com MixedCaps, não UPPER_CASE:
const defaultPort = "8080"
const maxRetries = 3
```

### `var` com valor zero explícito vs `:=`

```go
// Em internal/infrastructure/logger/logger.go:
var handler slog.Handler  // valor zero explícito — o tipo é importante para clareza

// Em internal/service/pokemon_service.go:
page := make([]domain.Pokemon, 0, len(pokemons))  // declaração curta
```

### Initialisms aplicados

```go
// Em todo o projeto: userID, JSON, HTTP, URL, DTO, UUID, JWT
type PokemonDTO struct {
    ID   string `json:"id"`
}
```

---

## 7. Exercício prático

**Objetivo:** Aplicar as regras de nomenclatura e formatação.

1. Navegue até `core/bff/mobile-bff/internal/domain/errors.go`
2. Observe como os erros são nomeados: `var ErrPokemonNotFound = errors.New(...)`
3. Crie uma nova constante no pacote `config/config.go` seguindo MixedCaps:

```go
// Antes (se existisse em snake_case):
// const API_TIMEOUT = 30

// Depois — MixedCaps:
const APITimeout = 30 * time.Second
```

4. Encontre uma variável no código e verifique:
   - Está com `:=` quando apropriado? (escopo local, valor inicial fornecido)
   - Ou com `var` quando o tipo precisa ser explícito?
   - O nome é proporcional ao escopo?

---

**Próxima fase:** [Fase 02 — Controle & Coleções](fase-02-controle-colecoes.md)
