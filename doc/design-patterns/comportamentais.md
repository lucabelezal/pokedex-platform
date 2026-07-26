# Padrões Comportamentais

Como Go orquestra comunicação entre objetos — delegando responsabilidades,
encapsulando algoritmos e gerenciando estado com o mínimo de cerimônia.

---

## Chain of Responsibility

### Propósito

Passar uma requisição por uma cadeia de handlers. Cada handler decide se processa
a requisição ou a passa adiante.

### Filosofia Go

Em Go, Chain of Responsibility é implementado com **interface + referência ao próximo**.
Cada handler implementa a interface e contém opcionalmente o próximo elo da cadeia.

No mundo real: **middleware HTTP é Chain of Responsibility**. Cada middleware
processa (ou não) e chama `next.ServeHTTP()`.

### Código idiomático

```go
package main

import "fmt"

type Handler interface {
    Processa(pokemon string) bool
    SetNext(Handler)
}

// Base com o próximo elo
type BaseHandler struct {
    next Handler
}

func (b *BaseHandler) SetNext(h Handler) { b.next = h }

// Handler concreto: verifica se pokémon existe
type ExisteHandler struct {
    BaseHandler
    pokedex map[string]bool
}

func (e *ExisteHandler) Processa(nome string) bool {
    if !e.pokedex[nome] {
        fmt.Println(nome, "não está na pokédex")
        return false
    }
    if e.next != nil {
        return e.next.Processa(nome)
    }
    return true
}

// Handler concreto: verifica nível
type NivelHandler struct {
    BaseHandler
    nivelMin int
}

func (n *NivelHandler) Processa(nome string) bool {
    fmt.Println(nome, "passou na verificação de nível")
    if n.next != nil {
        return n.next.Processa(nome)
    }
    return true
}

func main() {
    existe := &ExisteHandler{pokedex: map[string]bool{"Pikachu": true, "Mewtwo": true}}
    nivel := &NivelHandler{nivelMin: 20}

    existe.SetNext(nivel)

    existe.Processa("Pikachu")   // passa
    existe.Processa("Charizard") // bloqueado — não está na pokédex
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/chainOfResponsibility && go run .
```

---

## Command

### Propósito

Encapsular uma requisição como um objeto, permitindo parametrizar clientes com
diferentes requisições, enfileirar ou logar requisições.

### Filosofia Go

Go não precisa de uma classe `Command` — **funções são first-class citizens**.
Um command em Go é simplesmente:

```go
type Command func() error
```

Ou, se precisar de estado:

```go
type Command interface {
    Execute() error
}
```

### Código idiomático

```go
package main

import "fmt"

// Command como função — o jeito mais Go de fazer
type Comando func() string

// Invoker — executa comandos
type Treinador struct {
    comandos []Comando
}

func (t *Treinador) Executar(comando Comando) string {
    return comando()
}

func main() {
    ash := Treinador{}

    // Comandos são closures que capturam estado
    ataque := func() string { return "Pikachu usou Choque do Trovão!" }
    defesa := func() string { return "Pikachu usou Proteção!" }

    fmt.Println(ash.Executar(ataque))
    fmt.Println(ash.Executar(defesa))
}
```

### Command com estado via struct

```go
// Command como struct — útil para undo/redo ou filas de comandos
type CapturarCommand struct {
    pokemon string
}

func (c CapturarCommand) Execute() string {
    return fmt.Sprintf("Capturou %s!", c.pokemon)
}

func (c CapturarCommand) Undo() string {
    return fmt.Sprintf("Soltou %s!", c.pokemon)
}

type Historico struct {
    commands []CapturarCommand
}

func (h *Historico) Execute(c CapturarCommand) {
    fmt.Println(c.Execute())
    h.commands = append(h.commands, c)
}

func (h *Historico) UndoLast() {
    c := h.commands[len(h.commands)-1]
    h.commands = h.commands[:len(h.commands)-1]
    fmt.Println(c.Undo())
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/command && go run .
```

---

## Iterator

### Propósito

Fornecer uma maneira de acessar sequencialmente os elementos de uma coleção
sem expor sua representação interna.

### Filosofia Go

Iterator é **raro como padrão explícito** em Go porque `for range` resolve a
maioria dos casos nativamente. Slices, maps, channels e strings são iteráveis
sem interface nenhuma.

Só implemente Iterator explícito quando:
- A coleção é complexa (árvore, grafo)
- Você precisa de lazy evaluation (streaming de dados)
- Precisa de múltiplas formas de percorrer a mesma coleção

### Código idiomático — range nativo (95% dos casos)

```go
// Não precisa de Iterator — range já faz isso nativamente
pokemons := []string{"Pikachu", "Charizard", "Blastoise"}
for i, nome := range pokemons {
    fmt.Println(i, nome)
}

pokedex := map[string]int{"Pikachu": 25, "Charizard": 6}
for nome, level := range pokedex {
    fmt.Println(nome, level)
}
```

### Código idiomático — Iterator explícito (quando necessário)

```go
package main

import "fmt"

type Iterator interface {
    HasNext() bool
    Next() string
}

type PokemonCollection struct {
    pokemons []string
    index    int
}

func (p *PokemonCollection) HasNext() bool {
    return p.index < len(p.pokemons)
}

func (p *PokemonCollection) Next() string {
    p.index++
    return p.pokemons[p.index-1]
}

func main() {
    colecao := &PokemonCollection{pokemons: []string{"Pikachu", "Charizard"}}
    for colecao.HasNext() {
        fmt.Println(colecao.Next())
    }
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/iterator && go run .
```

---

## Mediator

### Propósito

Centralizar comunicações complexas entre objetos. Em vez de objetos se
comunicarem diretamente (N×N), eles se comunicam apenas com o mediator.

### Filosofia Go

Em Go, Mediator é uma **interface** que orquestra a comunicação. Os colegas
(Colleagues) notificam o mediator, que decide o que fazer.

Exemplos reais: um **orquestrador de saga** em microserviços, um **hub de
WebSocket** que gerencia conexões.

### Código idiomático

```go
package main

import "fmt"

// Mediator — orquestra a comunicação
type BatalhaMediator struct {
    pokemon1 *Pokemon
    pokemon2 *Pokemon
}

func (b *BatalhaMediator) Notificar(quem *Pokemon, evento string) {
    if evento == "atacou" {
        if quem == b.pokemon1 {
            b.pokemon2.ReceberDano()
        } else {
            b.pokemon1.ReceberDano()
        }
    }
}

// Colleague — se comunica apenas com o mediator
type Pokemon struct {
    nome     string
    hp       int
    mediator *BatalhaMediator
}

func (p *Pokemon) Atacar() {
    fmt.Println(p.nome, "atacou!")
    p.mediator.Notificar(p, "atacou")
}

func (p *Pokemon) ReceberDano() {
    p.hp -= 10
    fmt.Println(p.nome, "recebeu dano! HP:", p.hp)
}

func main() {
    mediator := &BatalhaMediator{}

    pikachu := &Pokemon{nome: "Pikachu", hp: 100, mediator: mediator}
    charizard := &Pokemon{nome: "Charizard", hp: 100, mediator: mediator}

    mediator.pokemon1 = pikachu
    mediator.pokemon2 = charizard

    pikachu.Atacar()
    charizard.Atacar()
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/mediator && go run .
```

---

## Memento

### Propósito

Capturar e externalizar o estado interno de um objeto para que ele possa ser
restaurado posteriormente — sem violar encapsulamento.

### Filosofia Go

Em Go, Memento é uma **struct imutável** que armazena o snapshot do estado.
O Originator cria o memento e pode restaurar a partir dele. O Caretaker gerencia
o histórico.

Útil para undo/redo e checkpoints.

### Código idiomático

```go
package main

import "fmt"

// Memento — snapshot imutável do estado
type PokemonMemento struct {
    nome  string
    level int
    hp    int
}

// Originator — cria e restaura mementos
type Pokemon struct {
    nome  string
    level int
    hp    int
}

func (p *Pokemon) SaveToMemento() *PokemonMemento {
    return &PokemonMemento{nome: p.nome, level: p.level, hp: p.hp}
}

func (p *Pokemon) RestoreFromMemento(m *PokemonMemento) {
    p.nome = m.nome
    p.level = m.level
    p.hp = m.hp
}

func (p *Pokemon) LevelUp() { p.level++; p.hp += 10 }

func main() {
    pikachu := &Pokemon{nome: "Pikachu", level: 25, hp: 80}

    snapshot := pikachu.SaveToMemento() // salva estado

    pikachu.LevelUp()
    fmt.Println("Após LevelUp:", pikachu.level, pikachu.hp) // 26, 90

    pikachu.RestoreFromMemento(snapshot) // restaura
    fmt.Println("Após restore:", pikachu.level, pikachu.hp) // 25, 80
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/memento && go run .
```

---

## Observer

### Propósito

Definir uma dependência um-para-muitos: quando um objeto (subject) muda de estado,
todos os seus dependentes (observers) são notificados automaticamente.

### Filosofia Go

Observer em Go pode ser implementado de duas formas:

1. **Interface + slice** — O subject mantém um `[]Observer` e notifica via loop
2. **Channel** — O subject envia eventos por um channel; observers leem em goroutines

A segunda forma é mais idiomática em Go — channels são nativos e thread-safe.

### Código idiomático — com channels

```go
package main

import (
    "fmt"
    "time"
)

// Subject publica eventos em um channel
type PokemonSubject struct {
    eventos chan string
}

func (s *PokemonSubject) Subscribe() <-chan string {
    return s.eventos
}

func (s *PokemonSubject) Capturado(nome string) {
    s.eventos <- fmt.Sprintf("Capturado: %s", nome)
}

func main() {
    subject := &PokemonSubject{eventos: make(chan string, 10)}

    // Observer 1
    go func() {
        for evento := range subject.Subscribe() {
            fmt.Println("Log:", evento)
        }
    }()

    // Observer 2
    go func() {
        for evento := range subject.Subscribe() {
            fmt.Println("Notificação push:", evento)
        }
    }()

    subject.Capturado("Pikachu")
    subject.Capturado("Charizard")
    time.Sleep(100 * time.Millisecond)
}
```

### Código idiomático — com interface + slice

```go
type Observer interface {
    Update(evento string)
}

type Subject struct {
    observers []Observer
}

func (s *Subject) Register(o Observer) {
    s.observers = append(s.observers, o)
}

func (s *Subject) Notify(evento string) {
    for _, o := range s.observers {
        o.Update(evento)
    }
}
```

**Channel vs slice:** Use channel quando os observers são goroutines independentes
(desacoplados). Use slice quando precisa de controle síncrono sobre a notificação.

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/observer && go run .
```

---

## State

### Propósito

Permitir que um objeto altere seu comportamento quando seu estado interno muda.
O objeto parece mudar de classe.

### Filosofia Go

State é implementado com uma **interface State** e structs concretas para cada
estado. O contexto (máquina de estados) mantém o estado atual e delega chamadas
a ele. Cada estado decide a transição.

Go não tem classes abstratas — usa-se interface + struct simples.

### Código idiomático

```go
package main

import "fmt"

// Interface State
type PokemonState interface {
    Atacar() string
    Defender() string
    ProximoEstado() PokemonState
}

// Estado concreto: Normal
type NormalState struct{}

func (NormalState) Atacar() string  { return "Ataque normal" }
func (NormalState) Defender() string { return "Defesa normal" }
func (NormalState) ProximoEstado() PokemonState { return FuryState{} }

// Estado concreto: Fury
type FuryState struct{}

func (FuryState) Atacar() string  { return "ATAQUE FURIOSO! Dano x2!" }
func (FuryState) Defender() string { return "Defesa reduzida na fúria" }
func (FuryState) ProximoEstado() PokemonState { return ExaustoState{} }

// Estado concreto: Exausto
type ExaustoState struct{}

func (ExaustoState) Atacar() string  { return "Sem energia para atacar..." }
func (ExaustoState) Defender() string { return "Sem energia para defender..." }
func (ExaustoState) ProximoEstado() PokemonState { return NormalState{} }

// Contexto
type Pokemon struct {
    nome  string
    state PokemonState
}

func (p *Pokemon) AvancarEstado() {
    p.state = p.state.ProximoEstado()
    fmt.Println(p.nome, "mudou para", fmt.Sprintf("%T", p.state))
}

func main() {
    pikachu := &Pokemon{nome: "Pikachu", state: NormalState{}}

    fmt.Println(pikachu.state.Atacar())  // Ataque normal
    pikachu.AvancarEstado()
    fmt.Println(pikachu.state.Atacar())  // ATAQUE FURIOSO!
    pikachu.AvancarEstado()
    fmt.Println(pikachu.state.Atacar())  // Sem energia...
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/state && go run .
```

---

## Strategy

### Propósito

Definir uma família de algoritmos, encapsular cada um e torná-los intercambiáveis.
O Strategy permite variar o algoritmo independentemente dos clientes.

### Filosofia Go

**Strategy é o padrão mais importante em Go.** Toda vez que você injeta uma
interface para permitir múltiplas implementações, está usando Strategy.

É a base da arquitetura hexagonal: o Service depende de uma interface
`PokemonRepository`. Em runtime, você injeta a implementação concreta
(Postgres, HTTP, Memory). O Service não sabe qual — só conhece a interface.

### Código idiomático — exatamente como no projeto real

```go
package main

import (
    "context"
    "fmt"
)

// Strategy — interface
type CapturaStrategy interface {
    CalcularChance(nivelPokemon int) float64
}

// Estratégia concreta 1
type PokebolaStrategy struct{}

func (PokebolaStrategy) CalcularChance(nivel int) float64 {
    return 1.0 - float64(nivel)/200.0
}

// Estratégia concreta 2
type UltraBallStrategy struct{}

func (UltraBallStrategy) CalcularChance(nivel int) float64 {
    return 1.0 - float64(nivel)/400.0
}

// Estratégia concreta 3
type MasterBallStrategy struct{}

func (MasterBallStrategy) CalcularChance(_ int) float64 {
    return 1.0
}

// Contexto — não sabe qual estratégia está usando
type Capturador struct {
    strategy CapturaStrategy
}

func (c *Capturador) SetEstrategia(s CapturaStrategy) {
    c.strategy = s
}

func (c *Capturador) TentarCaptura(ctx context.Context, pokemon string, nivel int) string {
    chance := c.strategy.CalcularChance(nivel)
    return fmt.Sprintf("Tentando capturar %s (chance: %.0f%%)", pokemon, chance*100)
}

func main() {
    capturador := &Capturador{strategy: PokebolaStrategy{}}

    fmt.Println(capturador.TentarCaptura(nil, "Pikachu", 25))    // 87%
    capturador.SetEstrategia(UltraBallStrategy{})
    fmt.Println(capturador.TentarCaptura(nil, "Charizard", 50))  // 87%
    capturador.SetEstrategia(MasterBallStrategy{})
    fmt.Println(capturador.TentarCaptura(nil, "Mewtwo", 100))    // 100%
}
```

### Por que Strategy é tão importante em Go

```
Handler ──→ inbound.UseCase (interface)    ← Strategy
                ↑
           Service (struct)
                ↓
           outbound.Repository (interface)  ← Strategy
                ↑
           PostgresRepo / HTTPClient        ← implementações concretas
```

Cada seta é um Strategy. A arquitetura hexagonal inteira é Strategy aplicado
em cascata. É por isso que Go dispensa frameworks de injeção de dependência —
o wiring manual no `main.go` é a aplicação do padrão Strategy.

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/strategy && go run .
```

---

## Template Method

### Propósito

Definir o esqueleto de um algoritmo, delegando alguns passos para subclasses.
A estrutura do algoritmo permanece fixa; os detalhes variam.

### Filosofia Go

Go não tem herança, então o Template Method usa **embedding de struct + interface**.
A struct base contém o método template (`genAndSendOTP()`). Os passos variáveis
são chamadas à interface `IOtp`, que cada subtipo implementa.

### Código idiomático

```go
package main

import "fmt"

// Interface com os passos variáveis
type BatalhaSteps interface {
    AntesDaBatalha()
    DepoisDaBatalha()
}

// Template — o esqueleto do algoritmo
type Batalha struct {
    steps BatalhaSteps
}

func (b *Batalha) Executar() {
    b.steps.AntesDaBatalha()
    fmt.Println("⚔️ Batalha iniciada!")
    // ... lógica fixa de batalha ...
    fmt.Println("🏆 Batalha concluída!")
    b.steps.DepoisDaBatalha()
}

// Subtipo 1 — implementa os passos
type BatalhaGinasio struct {
    Batalha  // embedding
}

func (BatalhaGinasio) AntesDaBatalha() {
    fmt.Println("🎵 Toca música de ginásio")
}

func (BatalhaGinasio) DepoisDaBatalha() {
    fmt.Println("🏅 Entrega insígnia ao treinador")
}

// Subtipo 2
type BatalhaSelvagem struct {
    Batalha
}

func (BatalhaSelvagem) AntesDaBatalha() {
    fmt.Println("🌿 Um Pokémon selvagem apareceu!")
}

func (BatalhaSelvagem) DepoisDaBatalha() {
    fmt.Println("💨 O Pokémon selvagem fugiu!")
}

func main() {
    ginasio := BatalhaGinasio{Batalha: Batalha{}}
    ginasio.steps = ginasio
    ginasio.Executar()

    selvagem := BatalhaSelvagem{Batalha: Batalha{}}
    selvagem.steps = selvagem
    selvagem.Executar()
}
```

**Atenção:** Esse padrão é mais verboso em Go do que em linguagens com herança.
Só use Template Method quando a estrutura do algoritmo é realmente fixa e
complexa. Na maioria dos casos, Strategy (injetar uma interface) é mais limpo.

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/template && go run .
```

---

## Visitor

### Propósito

Separar um algoritmo da estrutura de objetos sobre a qual ele opera.
Permite adicionar novas operações sem modificar as classes dos elementos.

### Filosofia Go

Visitor é **raro em Go**. Para a maioria dos casos, um **type switch** resolve
com menos código e sem a necessidade de double dispatch.

Só use Visitor quando:
- Você tem uma hierarquia de tipos estável
- Precisa adicionar muitas operações diferentes sobre essa hierarquia
- Type switch espalhado em muitos lugares se torna problema de manutenção

### Alternativa idiomática: type switch (95% dos casos)

```go
// Em vez de Visitor, use type switch
func calculaDano(movimento interface{}) int {
    switch m := movimento.(type) {
    case AtaqueFisico:
        return m.Poder
    case AtaqueEspecial:
        return m.Poder * 2
    case AtaqueStatus:
        return 0
    default:
        return 0
    }
}
```

### Quando usar Visitor (5% dos casos)

```go
package main

import "fmt"

// Element interface
type Movimento interface {
    Aceitar(v Visitor)
}

type Visitor interface {
    VisitarAtaqueFisico(m *AtaqueFisico)
    VisitarAtaqueEspecial(m *AtaqueEspecial)
}

// Elementos concretos
type AtaqueFisico struct {
    Nome  string
    Poder int
}
func (a *AtaqueFisico) Aceitar(v Visitor) { v.VisitarAtaqueFisico(a) }

type AtaqueEspecial struct {
    Nome  string
    Poder int
}
func (a *AtaqueEspecial) Aceitar(v Visitor) { v.VisitarAtaqueEspecial(a) }

// Visitor concreto: calcula dano
type CalculadoraDano struct {
    DanoTotal int
}

func (c *CalculadoraDano) VisitarAtaqueFisico(a *AtaqueFisico) {
    c.DanoTotal += a.Poder
}

func (c *CalculadoraDano) VisitarAtaqueEspecial(a *AtaqueEspecial) {
    c.DanoTotal += a.Poder * 2
}

func main() {
    movimentos := []Movimento{
        &AtaqueFisico{"Soco", 40},
        &AtaqueEspecial{"Choque do Trovão", 90},
    }

    calc := &CalculadoraDano{}
    for _, m := range movimentos {
        m.Aceitar(calc)
    }
    fmt.Println("Dano total:", calc.DanoTotal) // 40 + 180 = 220
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/visitor && go run .
```
