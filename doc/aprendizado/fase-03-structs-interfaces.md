# Fase 03 — Structs, Interfaces & Erros

Esta é a fase mais importante do roteiro. Structs, interfaces e tratamento de erros
são o coração do design em Go. É aqui que você entende por que Go é diferente de
outras linguagens — e como essa diferença produz código mais simples e previsível.

---

## Structs

Struct é o tipo composto básico de Go. Agrupa campos sob um nome. É um **value type**:
quando você atribui uma struct a outra variável, todo o conteúdo é copiado.

```go
package main

import "fmt"

type Pokemon struct {
    Nome  string
    Level int
    Tipos []string
    HP    int
}

func main() {
    // inicialização com nomes de campo (idiomático)
    p1 := Pokemon{
        Nome:  "Pikachu",
        Level: 25,
        Tipos: []string{"Elétrico"},
        HP:    80,
    }

    // inicialização posicional (evite — frágil)
    p2 := Pokemon{"Charizard", 36, []string{"Fogo", "Voador"}, 120}

    // zero value — todos os campos com valor padrão
    var p3 Pokemon

    fmt.Println(p1.Nome, p1.Level)
    fmt.Println(p2)
    fmt.Printf("Zero value: Nome=%q Level=%d HP=%d\n", p3.Nome, p3.Level, p3.HP)
}
```

### Struct é value type

```go
p1 := Pokemon{Nome: "Pikachu", Level: 25}
p2 := p1        // cópia completa!
p2.Level = 99
fmt.Println(p1.Level)  // 25 (não foi alterado)
```

### Campo não exportado = privado

```go
type Pokemon struct {
    Nome  string   // exportado (maiúscula) — visível fora do pacote
    level int      // não exportado (minúscula) — privado ao pacote
}
```

### Swift vs Go — structs

```swift
// Swift — struct também é value type
struct Pokemon {
    var nome: String
    var level: Int
}
var p1 = Pokemon(nome: "Pikachu", level: 25)
var p2 = p1  // cópia
```

```go
// Go — praticamente igual, mas sem var/let distinction
type Pokemon struct {
    Nome  string
    Level int
}
p1 := Pokemon{Nome: "Pikachu", Level: 25}
p2 := p1  // cópia
```

**Atenção:** Go não tem classes. Todo comportamento é adicionado via funções com
receiver (métodos) — veja a próxima seção.

---

## Métodos

Métodos são funções associadas a um tipo via **receiver**. O receiver aparece
antes do nome do método:

```go
package main

import "fmt"

type Pokemon struct {
    Nome  string
    Level int
    HP    int
}

// value receiver — não modifica o Pokémon
func (p Pokemon) Descreve() string {
    return fmt.Sprintf("%s (lvl %d, HP %d)", p.Nome, p.Level, p.HP)
}

// pointer receiver — modifica o Pokémon
func (p *Pokemon) LevelUp() {
    p.Level++
    p.HP += 10
}

// pointer receiver — recebe dano
func (p *Pokemon) TakeDamage(dano int) {
    p.HP -= dano
    if p.HP < 0 {
        p.HP = 0
    }
}

func main() {
    pikachu := Pokemon{Nome: "Pikachu", Level: 25, HP: 80}

    fmt.Println(pikachu.Descreve())  // value receiver
    pikachu.LevelUp()                // pointer receiver (Go converte automaticamente)
    pikachu.TakeDamage(30)           // pointer receiver

    fmt.Println(pikachu.Descreve())
}
```

### Pointer receiver vs Value receiver

| Pointer receiver `(p *Pokemon)` | Value receiver `(p Pokemon)` |
|----------------------------------|------------------------------|
| Modifica o valor original | Recebe uma cópia, modifica a cópia |
| Não copia a struct (eficiente pra structs grandes) | Copia a struct inteira |
| Pode ser chamado em `nil` | Não pode (panic) |
| **Use como padrão** | Use quando a struct é pequena e imutável |

**Regra prática:** Se **qualquer** método do tipo usa pointer receiver, **todos**
devem usar pointer receiver. Seja consistente.

### Receiver idiomático

Em Go, receivers são nomeados com 1 ou 2 letras — a primeira letra do tipo:

```go
func (p *Pokemon) LevelUp()      // p = Pokemon
func (s *Service) List()         // s = Service
func (c *Client) Do()            // c = Client
func (h *Handler) ServeHTTP()    // h = Handler
func (rb *ResponseBuilder) Build() // rb = ResponseBuilder
```

**Nunca** use `this` ou `self`. Isso não é Go.

### Swift vs Go — métodos

```swift
// Swift — métodos dentro da struct (mutating quando modifica)
struct Pokemon {
    var level: Int
    mutating func levelUp() { level += 1 }
    func descreve() -> String { "lvl \(level)" }
}
```

```go
// Go — métodos fora da struct, receiver explícito
type Pokemon struct { Level int }
func (p *Pokemon) LevelUp() { p.Level++ }
func (p Pokemon) Descreve() string { return fmt.Sprintf("lvl %d", p.Level) }
```

---

## Como ler código Go

Antes de continuar para interfaces e erros, você precisa dominar uma skill
fundamental: **ler e navegar código Go**. As perguntas mais comuns de quem
está começando são exatamente sobre isso.

Vamos dissecar uma função real, passo a passo:

```go
func (s *Service) LevelUp(ctx context.Context, id string) error {
    p, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return fmt.Errorf("level up: %w", err)
    }
    p.Level++
    return s.repo.Save(ctx, p)
}
```

### Passo 1 — Entenda o receiver: `func (s *Service)`

`(s *Service)` é o **receiver** — o equivalente ao `self` do Swift. Toda função
com receiver é um **método** que pertence a um tipo.

```go
// Função (sem receiver) — não pertence a tipo nenhum
func LevelUp(p *Pokemon) { p.Level++ }

// Método (com receiver) — pertence ao tipo *Service
func (s *Service) LevelUp(ctx context.Context, id string) error { ... }
//   ^^^^^^^^^^^ receiver
//   s é uma variável local que referencia a instância de Service
//   *Service significa que o método recebe um ponteiro (pode modificar)
```

Dentro do método, `s` é como o `self` de Swift — você acessa campos e outros
métodos através dele: `s.repo`, `s.cache`, `s.OutroMetodo()`.

**Por que `s`?** Convenção Go: receiver com 1-2 letras — a primeira do tipo.
`Service` → `s`, `Client` → `c`, `Handler` → `h`. **Nunca** `this` ou `self`.

### Passo 2 — Descubra o que é `s.repo`

`s.repo` é um **campo** da struct `Service`. Para saber o tipo dele, vá até a
definição da struct:

```go
type Service struct {
    repo PokemonRepository   // ← aqui! repo é do tipo PokemonRepository
}

// PokemonRepository é uma interface
type PokemonRepository interface {
    FindByID(ctx context.Context, id string) (*Pokemon, error)
    Save(ctx context.Context, p *Pokemon) error
}
```

`PokemonRepository` é uma **interface** — um contrato. `Service` não sabe se
o repository real é PostgreSQL, HTTP ou memória. Sabe apenas que ele tem os
métodos `FindByID` e `Save`.

### Passo 3 — Leia a assinatura para saber o que retorna

```go
FindByID(ctx context.Context, id string) (*Pokemon, error)
//        ─────── parâmetros ──────    ──── retornos ────
//                                     *Pokemon = um Pokémon ou nil
//                                     error    = erro ou nil
```

O padrão Go: **se pode falhar, retorna `(resultado, error)`**.

| Situação | `p` | `err` | Significado |
|----------|-----|-------|-------------|
| Sucesso | Pokémon válido | `nil` | Deu certo, use `p` |
| Falha | `nil` | não-nil | Deu errado, trate `err` |

**Regra:** se `err != nil`, o primeiro valor **não é confiável**. Trate o erro
e não use o valor.

### Passo 4 — Entenda o fluxo da função

Agora lemos a função linha a linha, no fluxo real de execução:

```go
func (s *Service) LevelUp(ctx context.Context, id string) error {
    // LINHA 1: chama o repository. Pattern: resultado, erro := chamada()
    p, err := s.repo.FindByID(ctx, id)

    // LINHA 2: indent error flow — trate o erro PRIMEIRO
    if err != nil {
        // LINHA 3: wrap do erro com contexto adicional
        // %w preserva o erro original para errors.Is/errors.As
        return fmt.Errorf("level up: %w", err)
    }

    // LINHA 4: happy path — se chegou aqui, p é válido
    p.Level++

    // LINHA 5: Save também retorna error.
    // Como LevelUp retorna error, podemos retornar o resultado direto.
    return s.repo.Save(ctx, p)
}
```

Fluxo visual:

```
FindByID(ctx, id)
    │
    ├── err != nil? ──→ return err (wrap)
    │
    └── sucesso: p.Level++
                      │
                      └── return Save(ctx, p)
                               │
                               ├── err != nil? ──→ return err
                               └── nil ──→ return nil
```

### Passo 5 — Por que `ctx context.Context` sempre primeiro?

Context é o **primeiro parâmetro** de toda função que faz I/O (banco, HTTP,
arquivos). É **convenção da linguagem** — você encontra isso em TODO código Go.

```go
func (s *Service) LevelUp(ctx context.Context, id string) error {
//                          ^^^ sempre primeiro, sempre chamado ctx
```

Context carrega três coisas:
- **Deadline/timeout** — "esta operação tem 5 segundos"
- **Cancelamento** — "usuário fechou a aba, cancela tudo"
- **Valores** — tracing ID, user ID (use com moderação)

Sem context, você não tem timeout nem cancelamento. Toda operação de I/O
deve recebê-lo e propagá-lo.

### Passo 6 — Por que `return s.repo.Save(ctx, p)` funciona direto?

Porque os tipos batem:

```go
func (s *Service) LevelUp(...) error {  // retorna error
    // ...
    return s.repo.Save(ctx, p)          // Save retorna error também
}
```

`Save` retorna `error`. `LevelUp` retorna `error`. Você pode retornar o
resultado de `Save` diretamente — é idiomático e limpo.

Se `LevelUp` retornasse `(int, error)` e `Save` retornasse `error`, você
precisaria tratar:

```go
func (s *Service) LevelUp(...) (int, error) {
    // ...
    if err := s.repo.Save(ctx, p); err != nil {
        return 0, fmt.Errorf("save: %w", err)
    }
    return p.Level, nil
}
```

### Resumo: como ler qualquer função Go

```
1. Olhe o receiver   → que tipo é? (s *Service → struct Service)
2. Olhe os parâmetros → o que entra? (ctx, id)
3. Olhe o retorno     → o que sai? (error)
4. Siga o fluxo       → erro primeiro, happy path depois
5. Para cada chamada  → vá na definição da interface, leia a assinatura
```

Exemplo prático de navegação no projeto real:

```
Handler → chama svc.GetByID(ctx, id)
           │
           └─ vá para ports/inbound/pokemon_usecase.go
              │  type PokemonUseCase interface {
              │      GetByID(ctx context.Context, id string) (*PokemonDetail, error)
              │  }
              │
              └─ implementação em service/pokemon_service.go
                 │  func (s *PokemonService) GetByID(...)
                 │      s.repo.GetByID(ctx, id)
                 │      └─ vá para ports/outbound/pokemon_repository.go
                 │         type PokemonRepository interface {
                 │             GetByID(ctx context.Context, id string) (*Pokemon, error)
                 │         }
                 │
                 └─ implementações concretas em adapters/outbound/
                    ├── postgres/pokemon_repository.go
                    └── http/pokemon_catalog_client.go
```

---

## Embedding (Composição)

Go não tem herança. O mecanismo de reuso é **embedding** (composição):

```go
package main

import "fmt"

type Animal struct {
    Nome string
}

func (a Animal) Falar() string {
    return "..."
}

type Cachorro struct {
    Animal       // embedding — Cachorro "é um" Animal
    Raca string
}

// Cachorro sobrescreve Falar
func (c Cachorro) Falar() string {
    return "Au au!"
}

func main() {
    c := Cachorro{
        Animal: Animal{Nome: "Rex"},
        Raca:   "Labrador",
    }

    fmt.Println(c.Nome)     // Rex — campo promovido de Animal
    fmt.Println(c.Falar())  // Au au! — método sobrescrito
    fmt.Println(c.Raca)     // Labrador
}
```

### Embedding de interfaces

```go
type Leitor interface {
    Read(p []byte) (n int, err error)
}

type Escritor interface {
    Write(p []byte) (n int, err error)
}

type LeitorEscritor interface {
    Leitor
    Escritor
}
```

### Swift vs Go — composição

```swift
// Swift — protocol + extensão (não tem herança de classe como mecanismo principal)
protocol Animal { var nome: String { get } }
struct Cachorro: Animal { var nome: String; var raca: String }
```

```go
// Go — embedding de struct
type Animal struct { Nome string }
type Cachorro struct {
    Animal
    Raca string
}
```

**Atenção:** Embedding não é herança. `Cachorro` **não** é um subtipo de `Animal`.
Você não pode passar um `Cachorro` onde se espera um `Animal`. Para isso, use interfaces.

---

## Interfaces

A interface é o conceito mais poderoso — e mais mal compreendido — de Go.
A diferença fundamental para Swift: em Go, a implementação é **implícita**.

Você **não declara** que um tipo implementa uma interface. Se o tipo tem os métodos
certos, ele implementa automaticamente.

```go
package main

import (
    "fmt"
    "math"
)

// define a interface
type Descritor interface {
    Descreve() string
}

// Pokemon tem o método Descreve() string → implementa Descritor automaticamente
type Pokemon struct {
    Nome  string
    Level int
}

func (p Pokemon) Descreve() string {
    return fmt.Sprintf("%s (lvl %d)", p.Nome, p.Level)
}

// Treinador também implementa Descritor
type Treinador struct {
    Nome    string
    Medals int
}

func (t Treinador) Descreve() string {
    return fmt.Sprintf("Treinador %s (%d medals)", t.Nome, t.Medals)
}

// função que aceita qualquer Descritor
func imprime(d Descritor) {
    fmt.Println(d.Descreve())
}

func main() {
    pikachu := Pokemon{Nome: "Pikachu", Level: 25}
    ash := Treinador{Nome: "Ash", Medals: 8}

    imprime(pikachu)  // Pikachu (lvl 25)
    imprime(ash)      // Treinador Ash (8 medals)
}
```

### Interface compliance check

Para garantir em tempo de compilação que um tipo implementa uma interface:

```go
var _ Descritor = (*Pokemon)(nil)     // compila? Pokemon implementa Descritor
var _ Descritor = Pokemon{}           // mesmo check com value receiver
```

Se o tipo não implementar a interface, o código **não compila**.
Isso é colocado no arquivo onde o tipo é definido (não na interface).

### Swift vs Go — a diferença crucial

```swift
// Swift — conformidade EXPLÍCITA
protocol Descritor {
    func descreve() -> String
}
struct Pokemon: Descritor {     // ← ": Descritor" obrigatório
    func descreve() -> String { "Pikachu" }
}
```

```go
// Go — conformidade IMPLÍCITA
type Descritor interface {
    Descreve() string
}
type Pokemon struct { /* ... */ }
func (p Pokemon) Descreve() string { return "Pikachu" }
// Pokemon NÃO declara ": Descritor". Só de ter o método, já implementa.
```

**Por que implícita?** Permite que você defina interfaces pequenas e specificas
depois que os tipos já existem. Você pode criar uma interface para um tipo que
está em um pacote de terceiros — sem modificar o pacote original. Isso é impossível
em Swift (precisaria de uma extension declarando conformidade).

### Interface vazia `interface{}` / `any`

```go
var x interface{}  // aceita qualquer valor (Go 1.17-)
var x any           // alias moderno (Go 1.18+), preferido

x = 42
x = "hello"
x = Pokemon{Nome: "Pikachu"}
```

### Regras para interfaces

- **Mantenha interfaces pequenas.** 1-3 métodos é o ideal. `io.Reader` tem 1 método.
- **Aceite interfaces, retorne structs.** Funções devem receber interfaces e retornar tipos concretos.
- **Defina interfaces onde são consumidas**, não onde são implementadas.
- **Nunca use ponteiro para interface** (`*Descritor`) — interfaces já são referências.

---

## Type Assertion e Type Switch

Quando você tem um valor de interface, pode precisar recuperar o tipo concreto:

```go
package main

import "fmt"

func inspeciona(x any) {
    // type assertion — comma ok
    str, ok := x.(string)
    if ok {
        fmt.Println("É uma string:", str)
        return
    }

    num, ok := x.(int)
    if ok {
        fmt.Println("É um int:", num)
        return
    }

    fmt.Println("Tipo desconhecido")
}

func inspecionaComSwitch(x any) {
    // type switch — mais idiomático para múltiplos tipos
    switch v := x.(type) {
    case string:
        fmt.Println("String:", v)
    case int:
        fmt.Println("Int:", v)
    case bool:
        fmt.Println("Bool:", v)
    default:
        fmt.Printf("Tipo desconhecido: %T\n", v)
    }
}

func main() {
    inspeciona("hello")
    inspeciona(42)
    inspecionaComSwitch(true)
}
```

**Atenção:** Type assertion **sem** comma ok causa panic se o tipo não corresponder:

```go
str := x.(string)  // PANIC se x não for string — NUNCA faça isso sem verificar
```

---

## Generics

Go 1.18 introduziu generics. Permite escrever funções e tipos que trabalham com
qualquer tipo:

```go
package main

import "fmt"

// função genérica — aceita slice de qualquer tipo comparável
func Contem[T comparable](slice []T, elemento T) bool {
    for _, v := range slice {
        if v == elemento {
            return true
        }
    }
    return false
}

// função com restrição de tipo
func Soma[T int | float64](a, b T) T {
    return a + b
}

func main() {
    fmt.Println(Contem([]string{"a", "b", "c"}, "b")) // true
    fmt.Println(Contem([]int{1, 2, 3}, 5))            // false
    fmt.Println(Soma(1, 2))                           // 3
    fmt.Println(Soma(1.5, 2.5))                       // 4.0
}
```

### Swift vs Go — generics

```swift
// Swift
func contem<T: Equatable>(_ slice: [T], _ elemento: T) -> Bool {
    return slice.contains(elemento)
}
```

```go
// Go
func Contem[T comparable](slice []T, elemento T) bool { ... }
// comparable é uma constraint built-in: tipos que suportam == e !=
```

---

## Erros

Em Go, erro **não é uma exceção**. É um valor como outro qualquer, retornado
pela função. Esta é uma das maiores diferenças para Swift.

### A interface `error`

```go
type error interface {
    Error() string
}
```

Qualquer tipo com método `Error() string` é um erro.

### Criando erros

```go
package main

import (
    "errors"
    "fmt"
)

// sentinel error — erro simples e reutilizável
var ErrNotFound = errors.New("nao encontrado")

// erro customizado — com dados estruturados
type ValidationError struct {
    Campo  string
    Motivo string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("campo %s: %s", e.Campo, e.Motivo)
}

func buscaPokemon(nome string) (string, error) {
    if nome == "" {
        return "", &ValidationError{Campo: "nome", Motivo: "nao pode ser vazio"}
    }
    if nome != "Pikachu" {
        return "", ErrNotFound
    }
    return "Pikachu encontrado", nil
}

func main() {
    _, err := buscaPokemon("Charizard")
    if err != nil {
        fmt.Println("Erro:", err)
    }
}
```

### Wrapping de erros

Use `fmt.Errorf` com `%w` para **wrap** de erros (preserva o erro original para
inspeção com `errors.Is` e `errors.As`):

```go
func processaPokemon(nome string) error {
    pokemon, err := buscaPokemon(nome)
    if err != nil {
        return fmt.Errorf("processa pokemon %s: %w", nome, err)
        //                           ↑ %w = wrap, inspecionável
    }
    fmt.Println(pokemon)
    return nil
}
```

Use `%v` quando quiser **ocultar** o erro original (não inspecionável):

```go
return fmt.Errorf("erro interno ao processar: %v", err)
```

### `errors.Is` — verifica se um erro é (ou envolve) um erro sentinela

```go
err := processaPokemon("Charizard")
if errors.Is(err, ErrNotFound) {
    fmt.Println("Pokémon não existe")
}
// Funciona mesmo que ErrNotFound esteja wrapped várias vezes
```

### `errors.As` — extrai um tipo específico de erro da cadeia

```go
err := processaPokemon("")
var vErr *ValidationError
if errors.As(err, &vErr) {
    fmt.Printf("Campo inválido: %s — %s\n", vErr.Campo, vErr.Motivo)
}
```

### Swift vs Go — erros

```swift
// Swift — throws + try/catch
enum AppError: Error {
    case notFound
    case validation(campo: String, motivo: String)
}

func buscaPokemon(_ nome: String) throws -> String {
    guard !nome.isEmpty else { throw AppError.validation(campo: "nome", motivo: "vazio") }
    guard nome == "Pikachu" else { throw AppError.notFound }
    return "Pikachu encontrado"
}

do {
    let pokemon = try buscaPokemon("Charizard")
} catch AppError.notFound {
    print("Não encontrado")
} catch let AppError.validation(campo, motivo) {
    print("Inválido: \(campo) - \(motivo)")
} catch {
    print("Erro: \(error)")
}
```

```go
// Go — erro como valor de retorno
var ErrNotFound = errors.New("nao encontrado")

func buscaPokemon(nome string) (string, error) {
    if nome == "" {
        return "", &ValidationError{Campo: "nome", Motivo: "vazio"}
    }
    if nome != "Pikachu" {
        return "", ErrNotFound
    }
    return "Pikachu encontrado", nil
}

result, err := buscaPokemon("Charizard")
if errors.Is(err, ErrNotFound) {
    fmt.Println("Não encontrado")
}
```

**Atenção:** Go não tem `try`/`catch`. Não existe "lançar exceção". Erros são
valores de retorno e o chamador decide se quer tratá-los ou propagá-los.
Isso torna o fluxo de erro **explícito e visível** em cada chamada.

### Regras para erros

- **Error strings em minúsculas, sem ponto final:** `errors.New("nao encontrado")`
- **Trate o erro uma vez:** logue OU retorne — nunca os dois.
- **`%w` para wrap inspecionável, `%v` para obscurecer.**
- **Sentinel errors** para condições verificáveis (`ErrNotFound`).
- **Erros customizados** para dados estruturados (`ValidationError`).

---

## `encoding/json`

O pacote `encoding/json` converte entre structs Go e JSON:

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Pokemon struct {
    Nome  string   `json:"nome"`
    Level int      `json:"nivel"`
    Tipos []string `json:"tipos,omitempty"` // omitempty: omite se vazio
    Hp    int      `json:"-"`               // "-": exclui do JSON
    interno string  // minúscula = não exportado = não vai pro JSON
}

func main() {
    // Marshal — struct → JSON
    p := Pokemon{Nome: "Pikachu", Level: 25, Tipos: []string{"Elétrico"}, Hp: 80}
    data, _ := json.Marshal(p)
    fmt.Println(string(data))
    // {"nome":"Pikachu","nivel":25,"tipos":["Elétrico"]}

    // MarshalIndent — struct → JSON formatado
    data, _ = json.MarshalIndent(p, "", "  ")
    fmt.Println(string(data))

    // Unmarshal — JSON → struct
    jsonStr := `{"nome":"Charizard","nivel":36,"tipos":["Fogo","Voador"]}`
    var p2 Pokemon
    json.Unmarshal([]byte(jsonStr), &p2)
    fmt.Println(p2.Nome, p2.Level)  // Charizard 36
}
```

### Swift vs Go — JSON

```swift
// Swift — Codable
struct Pokemon: Codable {
    var nome: String
    var nivel: Int
    var tipos: [String]?
}
let data = try JSONEncoder().encode(p)
let p = try JSONDecoder().decode(Pokemon.self, from: data)
```

```go
// Go — struct tags + encoding/json
type Pokemon struct {
    Nome  string   `json:"nome"`
    Nivel int      `json:"nivel"`
    Tipos []string `json:"tipos,omitempty"`
}
data, _ := json.Marshal(p)
json.Unmarshal(data, &p)
```

---

## Pacotes e visibilidade

Em Go, a visibilidade é determinada pela **capitalização do nome**:

- **Maiúscula** (`Pokemon`, `Descreve`) = exportado (público)
- **Minúscula** (`pokemon`, `descreve`) = não exportado (privado ao pacote)

```go
package pokemon

type Pokemon struct {        // exportado
    Nome  string             // exportado
    level int                // NÃO exportado
}

func NovoPokemon(nome string) *Pokemon {  // exportado (construtor)
    return &Pokemon{Nome: nome, level: 1}
}

func (p *Pokemon) LevelUp() {  // exportado
    p.level++
}

func (p *Pokemon) level() int {  // NÃO exportado
    return p.level
}
```

**Nomes de pacotes:** minúsculos, singular, curtos. `package pokemon`, não
`package pokemons` ou `package pokemon_utils`.

---

## Exercícios da Fase 03

### 1. Implemente uma interface

Crie uma interface `Lutador` com método `Atacar() int` (retorna dano). Crie dois
tipos (`Pikachu`, `Charizard`) que implementam `Lutador`. Escreva uma função
`batalha(a, b Lutador) string` que simula uma batalha e retorna o vencedor.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "fmt"
    "math/rand"
)

type Lutador interface {
    Atacar() int
    Nome() string
}

type Pikachu struct{}

func (p Pikachu) Atacar() int { return rand.Intn(20) + 10 }
func (p Pikachu) Nome() string { return "Pikachu" }

type Charizard struct{}

func (c Charizard) Atacar() int { return rand.Intn(30) + 15 }
func (c Charizard) Nome() string { return "Charizard" }

func batalha(a, b Lutador) string {
    danoA := a.Atacar()
    danoB := b.Atacar()
    if danoA > danoB {
        return fmt.Sprintf("%s vence com %d de dano!", a.Nome(), danoA)
    }
    return fmt.Sprintf("%s vence com %d de dano!", b.Nome(), danoB)
}

func main() {
    fmt.Println(batalha(Pikachu{}, Charizard{}))
}
```
</details>

### 2. Erro customizado com wrapping

Crie um erro `ErrPokemonFainted` e uma função `usaPokemon(nome string, hp int) (string, error)`.
Se `hp <= 0`, retorne `ErrPokemonFainted`. Na função que chama, faça wrap com `%w`
e verifique com `errors.Is`.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "errors"
    "fmt"
)

var ErrPokemonFainted = errors.New("pokemon desmaiado")

func usaPokemon(nome string, hp int) (string, error) {
    if hp <= 0 {
        return "", ErrPokemonFainted
    }
    return fmt.Sprintf("%s usou ataque!", nome), nil
}

func batalha(nome string, hp int) error {
    msg, err := usaPokemon(nome, hp)
    if err != nil {
        return fmt.Errorf("batalha: %w", err)
    }
    fmt.Println(msg)
    return nil
}

func main() {
    err := batalha("Pikachu", 0)
    if errors.Is(err, ErrPokemonFainted) {
        fmt.Println("Não pode lutar:", err)
    }
}
```
</details>

### 3. JSON marshal/unmarshal

Crie uma struct `Treinador` com campos `Nome`, `Idade`, `Pokemons []string`.
Adicione tags JSON. Escreva código que serializa para JSON e desserializa de volta.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Treinador struct {
    Nome     string   `json:"nome"`
    Idade    int      `json:"idade"`
    Pokemons []string `json:"pokemons,omitempty"`
}

func main() {
    t := Treinador{Nome: "Ash", Idade: 10, Pokemons: []string{"Pikachu", "Charizard"}}

    data, _ := json.MarshalIndent(t, "", "  ")
    fmt.Println(string(data))

    var t2 Treinador
    json.Unmarshal(data, &t2)
    fmt.Printf("%+v\n", t2)
}
```
</details>

### 4. Desafio: type switch com interfaces

Crie uma interface `Item` com método `Usar() string`. Implemente para `Pokebola`
e `Poção`. Escreva uma função que recebe `[]Item` e usa type switch para contar
quantos itens de cada tipo existem.

<details>
<summary>Gabarito</summary>

```go
package main

import "fmt"

type Item interface {
    Usar() string
}

type Pokebola struct{}

func (p Pokebola) Usar() string { return "Pokébola lançada!" }

type Pocao struct{}

func (p Pocao) Usar() string { return "Poção usada! +20 HP" }

func contaItens(itens []Item) (pokebolas, pocoes int) {
    for _, item := range itens {
        switch item.(type) {
        case Pokebola:
            pokebolas++
        case Pocao:
            pocoes++
        }
    }
    return
}

func main() {
    itens := []Item{Pokebola{}, Pocao{}, Pokebola{}, Pokebola{}}
    pb, pc := contaItens(itens)
    fmt.Printf("Pokébolas: %d, Poções: %d\n", pb, pc)
}
```
</details>

---

**Próxima fase:** [Fase 04 — Concorrência](fase-04-concorrencia.md)
