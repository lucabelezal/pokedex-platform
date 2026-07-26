# Padrões de Criação

Como Go resolve a criação controlada de objetos — sem construtores, sem `new`,
sem fábricas abstratas pesadas.

---

## Singleton

### Propósito

Garantir que uma classe tenha **uma única instância** e fornecer um ponto global
de acesso a ela.

### Filosofia Go

Em Go, singleton é trivial graças a `sync.Once`. Não há double-checked locking
manual, nem static initializers. `sync.Once` garante que o código de inicialização
execute **exatamente uma vez**, mesmo sob concorrência.

Além disso, Go não tem classes — o singleton é uma **variável no nível do pacote**
com inicialização lazy.

### Código idiomático

```go
package db

import (
    "database/sql"
    "sync"
)

var (
    instance *sql.DB
    once     sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        // Inicialização pesada — executa uma única vez
        db, err := sql.Open("postgres", "host=localhost dbname=pokedex")
        if err != nil {
            panic(err) // falha catastrófica na inicialização
        }
        instance = db
    })
    return instance
}
```

**Por que `sync.Once`, não `sync.Mutex` com `if instance == nil`?**

```go
// ERRADO — race condition entre a verificação e a atribuição
var mu sync.Mutex
func GetDB() *sql.DB {
    if instance == nil {       // thread A e B veem nil
        mu.Lock()              // ambas entram em sequência
        instance = createDB()  // createDB() executa DUAS vezes
        mu.Unlock()
    }
    return instance
}
```

`sync.Once` resolve isso atomicamente, sem lock manual.

### Quando usar em Go

- Conexão de banco de dados (exemplo acima)
- Configuração carregada uma vez (`config.Load()`)
- Logger global
- Client HTTP reutilizável (pool de conexões)

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/singleton/syncOnce && go run .
```

---

## Factory Method

### Propósito

Definir uma interface para criar um objeto, mas deixar as subclasses decidirem
qual classe instanciar. O Factory Method delega a criação para subtipos.

### Filosofia Go

Go não tem subclasses. Em vez disso, usamos uma **função construtora** que retorna
uma **interface**. Quem chama não sabe (nem precisa saber) qual struct concreta
está recebendo — só conhece a interface.

É o padrão mais comum em Go: `func NewX(params) Interface`.

### Código idiomático

```go
package main

import "fmt"

// Interface — o contrato
type Pokemon interface {
    Nome() string
    Tipo() string
}

// Structs concretas — detalhes internos
type Pikachu struct{}

func (p Pikachu) Nome() string { return "Pikachu" }
func (p Pikachu) Tipo() string { return "Elétrico" }

type Charizard struct{}

func (c Charizard) Nome() string { return "Charizard" }
func (c Charizard) Tipo() string { return "Fogo" }

// Factory — o Factory Method em Go é uma função, não uma classe
func NovoPokemon(nome string) Pokemon {
    switch nome {
    case "Pikachu":
        return Pikachu{}
    case "Charizard":
        return Charizard{}
    default:
        return nil
    }
}

func main() {
    p := NovoPokemon("Pikachu")  // retorna Pokemon (interface)
    fmt.Println(p.Nome(), "-", p.Tipo())
}
```

**Por que funciona sem subclasses?** Porque `Pikachu` e `Charizard` implementam
`Pokemon` automaticamente (interfaces implícitas). Não precisam declarar herança.

### Quando usar em Go

- Criar diferentes implementações de repository (Postgres, Memory, HTTP)
- Criar diferentes estratégias de cache (Redis, in-memory)
- Criar diferentes codecs (JSON, XML, Protobuf)

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/factory && go run .
```

---

## Abstract Factory

### Propósito

Fornecer uma interface para criar **famílias** de objetos relacionados sem
especificar suas classes concretas.

### Filosofia Go

Abstract Factory é **raro em Go**. O padrão adiciona uma camada extra de
indireção que raramente se justifica. Prefira Factory Method + injeção de
dependência.

Quando necessário, implemente com interface de factory e funções que retornam
interfaces de produto.

### Código idiomático (se precisar)

```go
package main

import "fmt"

// Produtos — interfaces
type Pokemon interface {
    Nome() string
}
type Treinador interface {
    Nome() string
}

// Fábrica abstrata — interface
type RegiaoFactory interface {
    CriarPokemon() Pokemon
    CriarTreinador() Treinador
}

// Fábrica concreta: Kanto
type KantoFactory struct{}

type Pikachu struct{}
func (Pikachu) Nome() string { return "Pikachu" }

type Ash struct{}
func (Ash) Nome() string { return "Ash" }

func (KantoFactory) CriarPokemon() Pokemon  { return Pikachu{} }
func (KantoFactory) CriarTreinador() Treinador { return Ash{} }

func main() {
    var fabrica RegiaoFactory = KantoFactory{}
    p := fabrica.CriarPokemon()
    t := fabrica.CriarTreinador()
    fmt.Println(t.Nome(), "capturou", p.Nome())
}
```

### Por que evitar em Go

Na prática, se você tem duas implementações de repository (Postgres e Memory),
não precisa de Abstract Factory. Basta:

```go
repo := NewPostgresPokemonRepo(db)   // ou NewMemoryPokemonRepo()
svc := pokemon.NewService(repo)       // Service não sabe qual repo recebeu
```

A injeção de dependência manual substitui a Abstract Factory na maioria dos casos.

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/abstractFactory && go run .
```

---

## Builder

### Propósito

Separar a construção de um objeto complexo da sua representação, permitindo o
mesmo processo de construção criar diferentes representações.

### Filosofia Go

Em Go, o Builder aparece de duas formas:

1. **Builder encadeado** — struct com métodos que retornam `*Builder` para
   encadeamento fluente. Útil para objetos com muitos campos opcionais.
2. **Director + Builder interface** — quando a construção tem passos fixos
   mas o resultado varia. Menos comum.

O padrão **Functional Options** (ver Fase 07) é a alternativa idiomática ao
Builder na maioria dos casos.

### Código idiomático — Builder encadeado

```go
package main

import "fmt"

type Pokemon struct {
    Nome  string
    Level int
    Shiny bool
    Item  string
}

type PokemonBuilder struct {
    pokemon Pokemon
}

func NovoPokemonBuilder(nome string) *PokemonBuilder {
    return &PokemonBuilder{pokemon: Pokemon{Nome: nome, Level: 1}}
}

func (b *PokemonBuilder) Level(l int) *PokemonBuilder {
    b.pokemon.Level = l
    return b
}

func (b *PokemonBuilder) Shiny() *PokemonBuilder {
    b.pokemon.Shiny = true
    return b
}

func (b *PokemonBuilder) ComItem(item string) *PokemonBuilder {
    b.pokemon.Item = item
    return b
}

func (b *PokemonBuilder) Build() Pokemon {
    return b.pokemon
}

func main() {
    pikachu := NovoPokemonBuilder("Pikachu").
        Level(50).
        Shiny().
        ComItem("Light Ball").
        Build()

    fmt.Printf("%+v\n", pikachu)
    // {Nome:Pikachu Level:50 Shiny:true Item:Light Ball}
}
```

### Builder vs Functional Options

```go
// Builder — verboso para quem consome, flexível
p := NovoPokemonBuilder("Pikachu").Level(50).Shiny().Build()

// Functional Options — idiomático, menos verboso
p := NovoPokemon("Pikachu", ComLevel(50), ComShiny())
```

Use Functional Options para APIs públicas. Use Builder quando a validação
entre campos é complexa e precisa de estado intermediário.

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/builder && go run .
```

---

## Prototype

### Propósito

Criar novos objetos copiando um objeto existente (protótipo), em vez de
instanciar do zero.

### Filosofia Go

Em Go, prototype é útil quando a criação é cara ou complexa, e você quer
clonar um objeto pré-configurado.

Para tipos simples, Go não tem `clone()` nativo. Você faz:
- Cópia por valor (`novo := original` para structs sem ponteiros)
- Método `Clone()` manual
- Serialização/desserialização (JSON, gob)

### Código idiomático

```go
package main

import "fmt"

type Pokemon struct {
    Nome     string
    Level    int
    Ataques  []string
    Treinador *string
}

// Clone profundo — copia slices e desreferencia ponteiros
func (p Pokemon) Clone() Pokemon {
    clone := p

    // slice precisa de cópia profunda
    clone.Ataques = make([]string, len(p.Ataques))
    copy(clone.Ataques, p.Ataques)

    // ponteiro precisa de desreferenciação
    if p.Treinador != nil {
        treinador := *p.Treinador
        clone.Treinador = &treinador
    }

    return clone
}

func main() {
    treinador := "Ash"
    original := Pokemon{
        Nome: "Pikachu", Level: 25,
        Ataques:  []string{"Choque do Trovão", "Cauda de Ferro"},
        Treinador: &treinador,
    }

    clone := original.Clone()
    clone.Level = 99
    clone.Ataques[0] = "Raio"

    fmt.Println("Original:", original.Level, original.Ataques[0]) // 25, Choque do Trovão
    fmt.Println("Clone:", clone.Level, clone.Ataques[0])          // 99, Raio
}
```

### Quando usar em Go

- Objetos de configuração com defaults complexos
- Templates de resposta HTTP
- Estruturas de dados imutáveis (copy-on-write)
- Test fixtures que variam pouco entre casos de teste

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/prototype && go run .
```
