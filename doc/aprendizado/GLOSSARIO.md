# Glossário Go → Swift

Termos e conceitos de Go com seus equivalentes em Swift. Use como dicionário durante o aprendizado.

## A

### array
**Go:** `[3]int` — coleção de tamanho fixo. O tamanho faz parte do tipo.
**Swift:** `Array` quando você sabe o tamanho exato, mas Swift não distingue tamanho no tipo.
**Nota:** Em Go quase sempre se usa slice no lugar de array.

### atribuição curta (short variable declaration)
**Go:** `x := 42` — declara e infere o tipo. Só funciona dentro de funções.
**Swift:** `let x = 42` (constante) ou `var x = 42` (variável).
**Nota:** `:=` em Go sempre cria variável mutável. Use `const` para constantes.

## B

### blank identifier
**Go:** `_` — descarta um valor. Usado quando uma função retorna algo que você não precisa.
**Swift:** `_` — mesmo conceito, usado em tuplas e switches.

### buffer (channel)
**Go:** `make(chan int, 10)` — channel que armazena até N valores antes de bloquear.
**Swift:** Similar a um `AsyncStream` com buffer, mas mais explícito.

## C

### channel
**Go:** `chan T` — tubo tipado para comunicação entre goroutines.
**Swift:** Não tem equivalente direto. O mais próximo seria `AsyncStream<T>` ou `PassthroughSubject<T, Never>` do Combine.

### comma ok idiom
**Go:** `v, ok := m["chave"]` — o segundo valor booleano indica se a chave existe no map.
**Swift:** Similar ao retorno de `dict[key]` que é `T?` (Optional), mas em Go é `(T, bool)`.

### composição (embedding)
**Go:** `type C struct { B }` — C "embute" B, herdando seus métodos.
**Swift:** Não tem equivalente sintático. O mais próximo seria composição manual + delegação.
**Nota:** Go não tem herança de classes. Composição é a ferramenta de reuso.

### const
**Go:** `const X = 42` — valor constante em tempo de compilação. Diferente de `var`.
**Swift:** `let x = 42` — mas `let` em Swift permite inicialização em runtime; `const` em Go é estritamente compile-time.

### context
**Go:** `context.Context` — carrega deadlines, cancelamento e valores entre chamadas.
**Swift:** Similar a `Task` + `TaskGroup` do Swift Concurrency, mas mais explícito.

## D

### defer
**Go:** `defer f.Close()` — agenda a execução para o final da função (LIFO).
**Swift:** `defer { ... }` — idêntico. Ambos executam ao sair do escopo.

## E

### embedding
Veja [composição](#composição).

### error (tipo)
**Go:** `error` é uma interface com um método: `Error() string`.
**Swift:** `Error` é um protocol com `localizedDescription`. Conceito similar.
**Nota:** Em Go você retorna `error` como segundo valor. Swift usa `throw`.

### errors.Is / errors.As
**Go:** `errors.Is(err, ErrNotFound)` — verifica se um erro é (ou envolve) um erro sentinela. `errors.As` faz type assertion em erros.
**Swift:** Equivalente a `catch let error as MeuErro` ou `switch error { case .notFound: ... }`.

### exported / unexported
**Go:** Nomes começando com maiúscula são públicos (exported); minúscula são privados (unexported).
**Swift:** `public` / `internal` / `private` — controle explícito com palavras-chave. Go usa capitalização.

## F

### função variádica
**Go:** `func soma(nums ...int) int` — recebe N argumentos do mesmo tipo.
**Swift:** `func soma(_ nums: Int...) -> Int` — mesma sintaxe conceitual.

## G

### garbage collector (GC)
**Go:** Coleta de lixo automática. Libera memória não referenciada.
**Swift:** ARC (Automatic Reference Counting) — deterministic, baseado em contagem. Go é não-determinístico.
**Nota:** Em Go você **não** precisa de `weak`/`unowned`. O GC resolve ciclos.

### generics
**Go:** `func Map[T any](slice []T, fn func(T) T) []T` — disponível desde Go 1.18.
**Swift:** `func map<T>(_ transform: (Element) -> T) -> [T]` — Swift tinha generics desde o início.

### goroutine
**Go:** `go fn()` — inicia execução concorrente. Extremamente leve (poucos KB de stack).
**Swift:** Equivalente a `Task { await fn() }` do Swift Concurrency. Mas goroutines são mais leves e você pode ter milhares.

## I

### indent error flow
**Go:** Padrão de tratar o erro primeiro e retornar cedo, mantendo o "happy path" sem indentação.
**Swift:** `guard let ... else { return }` — mesmo princípio.

### initialisms
**Go:** Regra de nomenclatura: siglas são todas maiúsculas ou todas minúsculas. `userID` (não `userId`), `HTTPServer` (não `HttpServer`).
**Swift:** Segue o mesmo padrão na maioria dos guias de estilo.

### interface
**Go:** `type Leitor interface { Read([]byte) (int, error) }` — conjunto de métodos. Implementação **implícita**.
**Swift:** `protocol Leitor { func read(into: Data) throws -> Int }` — implementação **explícita** (`: Leitor`).
**Nota:** Esta é a diferença mais importante. Em Go um tipo implementa uma interface automaticamente se tiver os métodos certos.

### iota
**Go:** Identificador que incrementa automaticamente em blocos `const`. Usado para criar enumerados.
**Swift:** Não tem equivalente. O mais próximo é `enum { case a = 0; case b = 1; ... }`.

## M

### method (método)
**Go:** `func (r Receptor) Nome()` — função com receiver. Pode ser value (`r T`) ou pointer (`r *T`).
**Swift:** `func nome()` dentro de `struct`/`class`/`enum`. `mutating func` quando modifica self em struct.

### módulo (module)
**Go:** Unidade de versionamento. Definido por `go.mod`.
**Swift:** Similar a um `Package.swift` (Swift Package Manager).

## P

### package
**Go:** Unidade de organização de código. Todo arquivo `.go` pertence a um pacote.
**Swift:** Similar a um módulo ou target.

### panic
**Go:** `panic("algo deu errado")` — interrompe a execução. Usado para erros irrecuperáveis.
**Swift:** `fatalError("...")` — mesma semântica.

### pointer
**Go:** `*int` — endereço de um valor na memória. `&v` pega o endereço, `*p` desreferencia.
**Swift:** `UnsafePointer<Int>` — existe mas é raro no dia a dia. Go usa ponteiros extensivamente.
**Nota:** Go não tem aritmética de ponteiros (sem `p++`). Ponteiros são seguros.

### port (arquitetura hexagonal)
**Go:** Interface que define um contrato (inbound = use cases, outbound = repositórios).
**Swift:** Igual a um protocol que define o contrato de um serviço.

## R

### range
**Go:** `for i, v := range slice` — itera sobre arrays, slices, maps, strings, channels.
**Swift:** Equivalente a `for (index, value) in array.enumerated()` ou `for value in array`.

### receiver
**Go:** Parâmetro antes do nome do método: `func (u *User) Nome()`. O "self" do Go.
**Swift:** O `self` implícito dentro de um método.

### recover
**Go:** `recover()` — captura um panic dentro de um `defer`.
**Swift:** Não tem equivalente. Panic em Go não deve ser usado como exceção.

### rune
**Go:** `rune` é um alias para `int32`. Representa um code point Unicode.
**Swift:** Equivalente a `Character` (embora `Character` possa ser um cluster de code points).

## S

### slice
**Go:** `[]int` — visão dinâmica sobre um array subjacente. Similar a um array redimensionável.
**Swift:** `Array<Int>` — praticamente o mesmo conceito.
**Nota:** Slice é um tipo valor (struct internamente) mas aponta para um array compartilhado. Cuidado ao copiar.

### struct
**Go:** `type User struct { Nome string }` — tipo valor, copiado na atribuição.
**Swift:** `struct User { var nome: String }` — idêntico. Ambos são value types.

## T

### type assertion
**Go:** `v, ok := x.(T)` — verifica se x é do tipo T. "Comma ok" para segurança.
**Swift:** `if let v = x as? T` — mesmo padrão de segurança.

### type switch
**Go:** `switch v := x.(type) { case int: ... case string: ... }` — switch no tipo dinâmico.
**Swift:** `switch x { case let v as Int: ... case let v as String: ... }`

## Z

### zero value
**Go:** Todo tipo tem um valor padrão: `0` para números, `""` para strings, `nil` para ponteiros/slices/maps/channels/interfaces.
**Swift:** Não tem zero value universal. Opcionais começam como `nil`, mas não-opcionais precisam ser inicializados.
**Nota:** Zero value é uma das ideias mais poderosas de Go. Um `sync.Mutex` com zero value já está pronto para uso.
