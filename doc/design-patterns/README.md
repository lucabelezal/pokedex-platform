# Design Patterns em Go

Este diretório explica os 22 padrões de design GoF adaptados à filosofia Go.
Cada padrão é apresentado com **código idiomático**, a **mentalidade Go** por
trás dele, e referência ao [repositório de exemplos executáveis](../../design-patterns/design-patterns-go/).

## Por que estudar padrões em Go?

Go não é orientado a objetos no sentido clássico. Não tem classes, herança, nem
construtores. Mas os **problemas** que os padrões resolvem — acoplamento, extensibilidade,
criação controlada de objetos — existem em qualquer linguagem.

A diferença é que Go resolve esses problemas com **menos cerimônia**: interfaces
implícitas, composição de structs, funções como first-class citizens e channels.

## Como Go simplifica (ou elimina) padrões

| Padrão GoF | Em Go... |
|-----------|---------|
| **Abstract Factory** | Prefira Factory Method + injeção de dependência. Raramente necessário. |
| **Singleton** | `sync.Once` resolve em 3 linhas. Sem double-checked locking manual. |
| **Iterator** | `for range` sobre slices, maps e channels é nativo. A interface `Iterator` raramente é necessária. |
| **Observer** | Channels + goroutines fazem o trabalho de pub/sub sem hierarquia de classes. |
| **Command** | Funções são first-class. `type Command func() error` já é um command. |
| **Template Method** | Embedding de struct + interface para os passos variáveis. Sem classes abstratas. |
| **Visitor** | Type switch substitui double dispatch na maioria dos casos. Raro em Go. |

## Índice

| Categoria | Arquivo | Padrões |
|-----------|---------|---------|
| Criação | [criacao.md](criacao.md) | Singleton, Factory Method, Abstract Factory, Builder, Prototype |
| Estruturais | [estruturais.md](estruturais.md) | Adapter, Bridge, Composite, Decorator, Facade, Flyweight, Proxy |
| Comportamentais | [comportamentais.md](comportamentais.md) | Chain of Responsibility, Command, Iterator, Mediator, Memento, Observer, State, Strategy, Template Method, Visitor |

## Como usar

1. Leia a explicação do padrão com a filosofia Go
2. Estude o código exemplo embutido no documento
3. Execute o exemplo real: `cd ../../design-patterns/design-patterns-go/<padrao> && go run .`
4. Compare com como você faria em Swift — as diferenças revelam a mentalidade Go

## Filosofia comum a todos os padrões em Go

Independente do padrão, estas regras se aplicam sempre:

- **Interfaces definem contratos, não hierarquias.** Quem consome define a interface.
- **Structs carregam estado.** Métodos adicionam comportamento.
- **Composição, nunca herança.** Embedding de struct e de interface.
- **Funções são valores.** Callbacks, estratégias e comandos são closures, não classes.
- **Zero values são úteis.** `sync.Mutex` e `bytes.Buffer` estão prontos sem inicialização.
- **Retorne structs concretas, aceite interfaces.** `func NewService(repo Repositorio) *Service`.
- **Erro é valor, não exceção.** Toda operação que pode falhar retorna `error`.

## Swift vs Go nos padrões

| Padrão em Swift | Em Go |
|----------------|-------|
| `protocol` + class | `interface` + struct |
| `Delegate` (weak var) | Interface injetada (GC resolve ownership) |
| Closure callback | Função como campo (`type Callback func(T) error`) |
| `static let shared` | `sync.Once` + var de pacote |
| `NotificationCenter` | Channel + goroutine |
| `Combine` Publisher | `<-chan T` |
| Extension em protocol | Método no mesmo pacote (não existe extension em Go) |
