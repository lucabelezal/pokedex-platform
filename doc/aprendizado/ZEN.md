# ZEN — Os Princípios do Go

Este documento reúne a filosofia de design da linguagem, a engenharia de valores
que guia o ecossistema, e o mapeamento de como os 22 padrões GoF se expressam em
Go idiomático. É a sua referência para **pensar em Go**, não apenas escrever em Go.

**Inspiração:** The Zen of Go (Dave Cheney), Effective Go, Package Names (blog oficial),
o repositório [design-patterns-go](../../design-patterns/design-patterns-go/) e a
[documentação completa de padrões em Go](../design-patterns/README.md).

---

## Os 10 valores da engenharia Go

Derivados do Zen of Go de Dave Cheney. Não são regras de sintaxe — são **valores**
que guiam decisões de design quando você não sabe qual caminho tomar.

### 1. Um bom pacote começa com um bom nome

O nome do pacote é o resumo de uma palavra do que ele provê. `package user`,
não `package utils`. Quando o nome não descreve mais o conteúdo, é hora de
dividir o pacote.

```go
// Ruim — nome genérico, acumula de tudo
package common

// Bom — nome específico, propósito claro
package pokemon
```

Um pacote com nome bom torna o código cliente legível naturalmente:
`pokemon.Service`, `pokemon.Repository`, `pokemon.NewService()`.

**Regras da documentação oficial:**
- minúsculo, singular, sem underscores ou mixedCaps
- evite `util`, `common`, `misc`, `helpers`, `api`, `types`
- evite nomes que colidem com pacotes padrão (`http`, `io`, `json`)
- o último elemento do path de importação deve ser o nome do pacote
- não roube bons nomes de variável do usuário (`bufio`, não `buf`)

[Leia o artigo original: Package Names — go.dev](https://go.dev/blog/package-names)

### 2. Simplicidade importa

*"Simples é melhor que complexo."* — PEP-20

Simple não significa fácil. Significa legível, mantível, previsível. Código
simples é código que outra pessoa consegue entender sem você explicar.

Go não otimiza para o menor número de linhas, nem para one-liners engenhosos.
Otimiza para clareza.

```go
// Engenhoso — difícil de entender
x := []int{}
for _, v := range data { if v > 10 { x = append(x, v*2) } }

// Simples — óbvio na primeira leitura
var result []int
for _, v := range data {
    if v <= 10 {
        continue
    }
    result = append(result, v*2)
}
```

### 3. Evite estado no nível do pacote

Estado global torna o código impossível de testar em isolamento e em paralelo.
Prefira encapsular estado em structs e injetar dependências.

```go
// Ruim — estado global, impossível testar em paralelo
package counter
var count int
func Increment() int { count++; return count }

// Bom — estado encapsulado, testável
package counter
type Counter struct { count int }
func (c *Counter) Increment() int { c.count++; return c.count }
```

### 4. Planeje para falha, não para sucesso

*"Erros nunca devem passar silenciosamente."* — PEP-20

Programadores Go pensam no caso de falha primeiro. O happy path vem depois.
Isso produz código que falha visivelmente no desenvolvimento, não silenciosamente
em produção.

```go
// O fluxo de leitura natural de uma função Go:
// 1. valide precondições (erro → return)
// 2. execute operação (erro → return)
// 3. happy path no final
func (s *Service) LevelUp(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("id obrigatório")      // falha primeiro
    }
    p, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return fmt.Errorf("level up: %w", err)   // falha
    }
    p.Level++
    return s.repo.Save(ctx, p)                   // happy path
}
```

### 5. Retorne cedo, evite aninhamento profundo

*"Plano é melhor que aninhado."* — PEP-20 reinterpretado

Cada nível de indentação consome memória de trabalho do leitor. Mantenha o
happy path próximo à margem esquerda. Trate erros e saia cedo.

```go
// Ruim — 3 níveis de indentação
if user != nil {
    if user.IsActive {
        if user.HasPermission {
            doSomething()
        }
    }
}

// Bom — guard clauses, happy path na margem
if user == nil { return }
if !user.IsActive { return }
if !user.HasPermission { return }
doSomething()
```

### 6. Antes de lançar uma goroutine, saiba quando ela vai parar

Três letras — `g`, `o`, espaço — e você criou um processo concorrente. Mas
Go não tem `stop()` ou `kill()`. Você precisa saber:

1. **Sob que condição a goroutine vai parar?** (channel fechado, context cancelado)
2. **O que precisa acontecer para essa condição surgir?** (quem fecha o channel?)
3. **Que sinal você usa para saber que ela parou?** (waitgroup, channel de resposta)

```go
// Sempre tenha um caminho de saída claro
func worker(ctx context.Context, jobs <-chan Job, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case job, ok := <-jobs:
            if !ok { return }       // channel fechado → para
            process(job)
        case <-ctx.Done():
            return                   // contexto cancelado → para
        }
    }
}
```

### 7. Deixe a concorrência para quem chama

Se seu código pode ser executado concorrentemente, deixe a decisão de **como**
para quem chama. Bibliotecas devem ser síncronas por padrão.

```go
// Ruim — a biblioteca decide ser concorrente
func Process(data []Item) {
    go func() { /* ... */ }()  // quem chamou não controla isso
}

// Bom — síncrono. Se o chamador quiser concorrência, ele faz:
func Process(data []Item) { /* ... */ }
// Uso: go Process(data)
```

### 8. Escreva testes para travar o comportamento da API

Testes são o contrato executável do que seu código faz. Se cada comportamento
tem um teste, você sabe com `go test` se sua mudança quebrou algo.

Mudanças na API pública devem vir acompanhadas de mudanças nos testes.

### 9. Moderação é virtude

Go tem só 25 palavras-chave. As features que existem são tentadoras de usar
em todo lugar — channels pra tudo, embedding pra tudo, interfaces pra tudo.

Resista. Comece simples. Adicione complexidade só quando necessário.

### 10. Manutenibilidade é o que importa

*"Legibilidade conta."* — PEP-20

O objetivo real não é legibilidade — é **manutenibilidade**. Código que sobrevive
ao autor original. Código que gera valor no futuro, não só no presente.

Se seu código não puder ser mantido, ele será reescrito — e talvez não em Go.

---

## Go e os 22 padrões GoF

Os padrões de design clássicos se expressam de forma natural em Go usando
interfaces, composição de structs e funções. Abaixo, cada padrão com sua
tradução idiomática.

> **Repositório de exemplos:** [`../../design-patterns/design-patterns-go/`](../../design-patterns/design-patterns-go/)
> Execute qualquer padrão com `cd <padrão> && go run .`

### Padrões de Criação

| Padrão | Como Go implementa | Diretório |
|--------|-------------------|-----------|
| **Singleton** | `sync.Once` + variável no nível do pacote. Zero boilerplate. | `singleton/syncOnce/` |
| **Factory Method** | Função `NewTipo()` que retorna interface. Sem classes abstratas. | `factory/` |
| **Abstract Factory** | Interface de factory com métodos `CreateX()` retornando interfaces. | `abstractFactory/` |
| **Builder** | Struct builder com métodos encadeados ou Director + interface `IBuilder`. | `builder/` |
| **Prototype** | Método `Clone()` na interface. Deep copy manual ou via serialização. | `prototype/` |

**O mais idiomático:** Singleton com `sync.Once` é o padrão canônico Go — seguro,
simples e lazy por natureza.

### Padrões Estruturais

| Padrão | Como Go implementa | Diretório |
|--------|-------------------|-----------|
| **Adapter** | Struct que wrapped o adaptee e implementa a interface alvo. Delegação pura. | `adapter/` |
| **Bridge** | Separa abstração (`Computer`) de implementação (`Printer`) via interfaces. Troca em runtime. | `bridge/` |
| **Composite** | Interface comum (`Component`) + struct que contém `[]Component`. Recursão natural. | `composite/` |
| **Decorator** | Struct que wrapped outro e implementa a mesma interface. Adiciona comportamento. | `decorator/` |
| **Facade** | Struct que agrega múltiplos subsistemas. Métodos que orquestram chamadas. | `facade/` |
| **Flyweight** | Factory com cache (`map[string]Dress`). Estado intrínseco compartilhado. | `flyweight/` |
| **Proxy** | Struct que implementa a mesma interface do real e adiciona controle de acesso. | `proxy/` |

**O mais idiomático:** Composite e Decorator são naturais em Go porque interfaces
de 1-2 métodos são a norma. Um `io.Reader` que wrapped outro `io.Reader` é um
decorator — e está em toda a biblioteca padrão.

### Padrões Comportamentais

| Padrão | Como Go implementa | Diretório |
|--------|-------------------|-----------|
| **Chain of Resp.** | Interface `Department` com `execute()` e `setNext()`. Embedding de base struct. | `chainOfResponsibility/` |
| **Command** | Interface `Command` com `execute()`. Structs concretas encapsulam receiver + parâmetros. | `command/` |
| **Iterator** | Interface `Iterator` com `hasNext()`/`next()`. Slice interno + cursor. | `iterator/` |
| **Mediator** | Interface `Mediator` centraliza comunicação. Structs notificam o mediator, não entre si. | `mediator/` |
| **Memento** | Struct `Memento` imutável. `Originator` cria e restaura. `Caretaker` guarda histórico. | `memento/` |
| **Observer** | Interface `Observer` registrada no `Subject`. Notificação via loop sobre slice. | `observer/` |
| **State** | Interface `State` com métodos por transição. Cada estado é um struct que referencia a máquina. | `state/` |
| **Strategy** | Interface `Strategy` injetada no contexto. Troca de algoritmo em runtime. | `strategy/` |
| **Template Method** | Struct base com método template que chama steps via interface. Subtipos implementam steps. | `template/` |
| **Visitor** | Interface `Visitor` com método por tipo. Double dispatch via `accept(Visitor)`. | `visitor/` |

**O mais idiomático:** Strategy é onipresente em Go. Toda vez que você injeta
uma interface em uma struct para permitir múltiplas implementações, você está
usando Strategy. É o fundamento da arquitetura hexagonal.

---

## Como executar os exemplos

O repositório está em `design-patterns/design-patterns-go/`. Cada padrão é um
pacote `main` independente:

```bash
# Navegue até o padrão desejado
cd ../../design-patterns/design-patterns-go/strategy
go run .

# Singleton — duas implementações
cd ../../design-patterns/design-patterns-go/singleton/default
go run .
cd ../../design-patterns/design-patterns-go/singleton/syncOnce
go run .
```

## Os 22 padrões em 3 categorias de relevância para backend Go

### Uso diário (alta relevância)

**Strategy, Decorator, Factory Method, Singleton (sync.Once), Adapter, Facade, Proxy, Observer, Command, Chain of Responsibility**

São os padrões que você encontra em TODO código Go bem estruturado.
Strategy = injeção de interface. Decorator = middleware HTTP.
Factory Method = `NewService(repo)`. Adapter = tradução entre camadas.

### Uso pontual (média relevância)

**Builder, Composite, Iterator, Mediator, Prototype, State, Template Method, Bridge**

Builder é útil para configs complexas. Composite aparece em estruturas de árvore
(menus, diretórios). State é valioso para máquinas de estado (workflows).

### Uso raro em Go (baixa relevância)

**Abstract Factory, Flyweight, Memento, Visitor**

Abstract Factory é pesado para Go — prefira Factory Method + injeção.
Visitor é verboso e raramente necessário.
Flyweight é útil para jogos, não para APIs.

---

## Swift vs Go: como os padrões se traduzem

| Padrão Swift/iOS | Equivalente Go | Nota |
|-----------------|----------------|------|
| `Delegate` pattern | Interface com 1 método injetada | Não precisa de `weak` — GC resolve |
| `Closure` callback | Função como campo de struct | `type Callback func(T) error` |
| `Singleton` via `static let` | `sync.Once` + var de pacote | Go é lazy por padrão com `sync.Once` |
| `NotificationCenter` | Channel + Observer pattern | `chan Event` + goroutine listener |
| `MVVM` / `VIPER` | Hexagonal (ports & adapters) | Interfaces substituem protocolos |
| `Combine` / Publisher | Channel + goroutine | `<-chan T` é um stream nativo |
| `@Published` / `@State` | Não existe. Estado é explícito. | Mutex, channel, ou atomic |
| `Coordinator` pattern | Middleware chain + context | `http.Handler` wrapping |

---

## O "Go way" em uma frase

> Go não é sobre escrever menos código. É sobre escrever código que **faz menos
> coisas inesperadas**.

Cada decisão de design da linguagem — interfaces implícitas, erro como valor,
goroutines leves, zero values — converge para um objetivo: **código previsível
e mantível**.

Isso não é um paradigma. É uma disciplina.

---

**Documentos relacionados:**
- [ROTEIRO.md](ROTEIRO.md) — roteiro completo de aprendizado
- [Fase 07 — Testes Avançados & Filosofia de Design](fase-07-testes-avancados-design.md) — filosofia e testes
- [GLOSSARIO.md](GLOSSARIO.md) — dicionário Swift ↔ Go
- [CHEATSHEET.md](CHEATSHEET.md) — sintaxe rápida
- [Design Patterns em Go](../design-patterns/README.md) — documentação completa dos 22 padrões GoF em Go idiomático
