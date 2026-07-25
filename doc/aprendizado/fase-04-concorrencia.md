# Fase 04 — Concorrência

A concorrência é onde Go brilha. Enquanto Swift tem `async/await` (introduzido em
2021), Go nasceu com concorrência no DNA — goroutines, channels e `select` são
primitivas da linguagem desde a versão 1.0.

Ao final desta fase, você vai entender o mantra de Go:
**"Do not communicate by sharing memory; share memory by communicating."**

---

## Goroutines

Uma goroutine é uma **thread leve** gerenciada pelo runtime de Go. Você a inicia
com a palavra-chave `go`:

```go
package main

import (
    "fmt"
    "time"
)

func imprimir(msg string) {
    for i := 0; i < 3; i++ {
        fmt.Println(msg, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    go imprimir("goroutine")  // executa concorrentemente
    imprimir("main")          // executa na goroutine principal

    // Sem time.Sleep, a main terminaria antes da goroutine rodar
    time.Sleep(1 * time.Second)
}
```

Saída (ordem varia a cada execução):

```
main 0
goroutine 0
goroutine 1
main 1
main 2
goroutine 2
```

Goroutines são **extremamente leves** — você pode ter dezenas de milhares rodando
simultaneamente. Cada uma começa com ~2 KB de stack, que cresce dinamicamente.

### Swift vs Go

```swift
// Swift — async/await + Task
func imprimir(_ msg: String) async {
    for i in 0..<3 {
        print(msg, i)
        try? await Task.sleep(nanoseconds: 100_000_000)
    }
}

Task { await imprimir("task") }
await imprimir("main")
```

```go
// Go — goroutines
func imprimir(msg string) {
    for i := 0; i < 3; i++ {
        fmt.Println(msg, i)
        time.Sleep(100 * time.Millisecond)
    }
}

go imprimir("goroutine")
imprimir("main")
```

**Atenção:** Se a goroutine principal (`main`) terminar, **todas** as outras
goroutines são encerradas imediatamente. Use `sync.WaitGroup` ou channels para
esperar a conclusão — nunca use `time.Sleep` em código real.

---

## Channels

Channels são **tubos tipados** que conectam goroutines. Você envia e recebe valores
através deles. A operação é **bloqueante** por padrão: um envio bloqueia até que
alguém receba; um recebimento bloqueia até que alguém envie.

```go
package main

import "fmt"

func main() {
    ch := make(chan string)  // channel unbuffered de strings

    // goroutine envia
    go func() {
        ch <- "Pikachu"       // bloqueia até alguém receber
    }()

    // main recebe
    msg := <-ch               // bloqueia até alguém enviar
    fmt.Println(msg)          // "Pikachu"
}
```

### Channels buffered

```go
ch := make(chan int, 3)  // buffer de 3 inteiros

ch <- 1   // não bloqueia (buffer tem espaço)
ch <- 2   // não bloqueia
ch <- 3   // não bloqueia (buffer cheio)
ch <- 4   // BLOQUEIA! buffer lotado — espera alguém consumir
```

### Channel directions

Você pode restringir a direção de um channel em parâmetros de função:

```go
// send-only — só pode enviar
func produz(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}

// receive-only — só pode receber
func consome(ch <-chan int) {
    for v := range ch {
        fmt.Println(v)
    }
}

func main() {
    ch := make(chan int)
    go produz(ch)
    consome(ch)
}
```

### `close` e `range` sobre channels

```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
ch <- 3
close(ch)  // sinaliza "não haverá mais envios"

// range itera até o channel ser fechado
for v := range ch {
    fmt.Println(v)  // 1 2 3
}

// ler de um channel fechado retorna zero value
v, ok := <-ch
fmt.Println(v, ok)  // 0 false
```

**Atenção:** Só o **sender** deve fechar o channel. Enviar em um channel fechado
causa **panic**. Receber de um channel fechado retorna zero value + `ok=false`.

---

## `select`

`select` permite esperar por **múltiplas operações** de channel simultaneamente:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "Pikachu"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Charizard"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Recebido do ch1:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Recebido do ch2:", msg2)
        case <-time.After(3 * time.Second):
            fmt.Println("Timeout!")
            return
        }
    }
}
```

### `select` com `default` — não bloqueante

```go
select {
case msg := <-ch:
    fmt.Println("Recebido:", msg)
default:
    fmt.Println("Nenhuma mensagem disponível")
}
```

### Swift vs Go — channels + select

```swift
// Swift — AsyncStream com TaskGroup (não tem equivalente direto a select)
// O mais próximo seria usar múltiplas Tasks e um contador
```

```go
// Go — select é primitiva nativa
select {
case v := <-ch1:  // espera ch1
case v := <-ch2:  // OU ch2
case <-timeout:    // OU timeout
}
```

Não existe equivalente direto a `select` em Swift. É uma das features que tornam
Go tão poderoso para sistemas concorrentes.

---

## `sync.WaitGroup`

WaitGroup espera que um grupo de goroutines termine:

```go
package main

import (
    "fmt"
    "sync"
)

func processa(id int, wg *sync.WaitGroup) {
    defer wg.Done()          // decrementa o contador ao sair
    fmt.Printf("Worker %d processando\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)            // incrementa antes de lançar goroutine
        go processa(i, &wg)
    }

    wg.Wait()                // bloqueia até o contador chegar a zero
    fmt.Println("Todos terminaram!")
}
```

**Regra:** `Add` deve ser chamado **antes** de `go`, não dentro da goroutine.

---

## `sync.Mutex`

Mutex protege acesso a dados compartilhados entre goroutines:

```go
package main

import (
    "fmt"
    "sync"
)

type Contador struct {
    mu    sync.Mutex
    valor int
}

func (c *Contador) Incrementa() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.valor++
}

func (c *Contador) Valor() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.valor
}

func main() {
    var c Contador
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Incrementa()
        }()
    }

    wg.Wait()
    fmt.Println("Valor final:", c.Valor())  // 1000
}
```

**Atenção:** O zero value de `sync.Mutex` é um mutex pronto para uso.
`var mu sync.Mutex` — sem `make()`, sem inicialização.

### `sync/atomic` — contadores sem mutex

Para operações simples de incremento/decremento, use `atomic` (mais rápido):

```go
import "sync/atomic"

var contador int64
atomic.AddInt64(&contador, 1)
v := atomic.LoadInt64(&contador)
```

---

## `context.Context`

Context carrega **deadlines, cancelamento e valores** através das chamadas de
função. É o mecanismo padrão para propagar cancelamento em Go.

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func buscaPokemon(ctx context.Context, nome string) (string, error) {
    // simula uma operação lenta (ex: chamada HTTP)
    select {
    case <-time.After(2 * time.Second):
        return fmt.Sprintf("%s encontrado!", nome), nil
    case <-ctx.Done():
        return "", ctx.Err()  // retorna o erro de cancelamento/timeout
    }
}

func main() {
    // context com timeout de 1 segundo
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()  // libera recursos se a operação terminar antes do timeout

    resultado, err := buscaPokemon(ctx, "Pikachu")
    if err != nil {
        fmt.Println("Erro:", err)  // context deadline exceeded
        return
    }
    fmt.Println(resultado)
}
```

### Tipos de context

```go
ctx := context.Background()                         // contexto raiz (geralmente em main)
ctx := context.TODO()                               // placeholder (quando não sabe qual usar)

ctx, cancel := context.WithCancel(parent)           // cancelamento manual
ctx, cancel := context.WithTimeout(parent, 5*time.Second)  // timeout
ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second)) // deadline
ctx := context.WithValue(parent, key, value)        // valores (use com moderação)
```

### Swift vs Go — context

```swift
// Swift — Task cancellation + Task.sleep
let task = Task {
    try await Task.sleep(nanoseconds: 2_000_000_000)
    return "Pikachu encontrado!"
}
task.cancel()  // cancela manualmente
let resultado = try await task.value
```

```go
// Go — context com cancelamento/timeout explícito
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
resultado, err := buscaPokemon(ctx, "Pikachu")
```

**Regra de ouro:** Context deve ser o **primeiro parâmetro** da função:

```go
func Busca(ctx context.Context, nome string) (string, error) { ... }
```

**Nunca** armazene context em um campo de struct:

```go
type Service struct {
    ctx context.Context  // NUNCA faça isso
}

// Correto: passe context como parâmetro
func (s *Service) Busca(ctx context.Context, nome string) (string, error) { ... }
```

---

## Padrões de concorrência

### Worker pool

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        fmt.Printf("Worker %d processando job %d\n", id, job)
        results <- job * 2
    }
}

func main() {
    const numJobs = 10
    const numWorkers = 3

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    var wg sync.WaitGroup
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()
    close(results)

    for r := range results {
        fmt.Println("Resultado:", r)
    }
}
```

### Fan-in (multiplexador)

Combina múltiplos channels em um só:

```go
func fanIn(ch1, ch2 <-chan string) <-chan string {
    out := make(chan string)
    go func() {
        for {
            select {
            case v := <-ch1:
                out <- v
            case v := <-ch2:
                out <- v
            }
        }
    }()
    return out
}
```

### Pipeline

```go
func gera(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func dobra(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * 2
        }
        close(out)
    }()
    return out
}

func main() {
    numeros := gera(1, 2, 3, 4, 5)
    dobrados := dobra(numeros)
    for v := range dobrados {
        fmt.Println(v)  // 2 4 6 8 10
    }
}
```

---

## Race detector

Go inclui um detector de data races. Use-o sempre ao testar código concorrente:

```bash
go test -race ./...
go run -race main.go
```

Exemplo de código com data race:

```go
func main() {
    contador := 0
    for i := 0; i < 1000; i++ {
        go func() {
            contador++   // DATA RACE!
        }()
    }
}
```

Execute `go run -race main.go` e o race detector apontará o problema.

---

## Exercícios da Fase 04

### 1. Goroutine com WaitGroup

Escreva uma função que lança 5 goroutines, cada uma imprimindo seu ID e dormindo
por um tempo aleatório. Use `sync.WaitGroup` para esperar todas terminarem.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func trabalhador(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    duracao := time.Duration(rand.Intn(500)) * time.Millisecond
    time.Sleep(duracao)
    fmt.Printf("Trabalhador %d terminou em %v\n", id, duracao)
}

func main() {
    var wg sync.WaitGroup
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go trabalhador(i, &wg)
    }
    wg.Wait()
    fmt.Println("Todos terminaram!")
}
```
</details>

### 2. Channel com timeout

Escreva uma função `buscaComTimeout(nome string, timeout time.Duration) (string, error)`
que usa um channel e `select` com `time.After` para simular uma busca. Se demorar
mais que o timeout, retorne erro.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "errors"
    "fmt"
    "time"
)

func buscaComTimeout(nome string, timeout time.Duration) (string, error) {
    resultado := make(chan string)

    go func() {
        time.Sleep(300 * time.Millisecond) // simula busca lenta
        resultado <- nome + " encontrado!"
    }()

    select {
    case res := <-resultado:
        return res, nil
    case <-time.After(timeout):
        return "", errors.New("timeout ao buscar")
    }
}

func main() {
    fmt.Println(buscaComTimeout("Pikachu", 500*time.Millisecond))  // sucesso
    fmt.Println(buscaComTimeout("Pikachu", 100*time.Millisecond))  // timeout
}
```
</details>

### 3. Context com cancelamento

Escreva uma função `processaComCancelamento(ctx context.Context, id int)` que
simula processamento longo. Lance várias goroutines. Cancele o context após 500ms
e verifique que todas as goroutines param.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func processaComCancelamento(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Goroutine %d cancelada\n", id)
            return
        default:
            fmt.Printf("Goroutine %d trabalhando...\n", id)
            time.Sleep(200 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    for i := 1; i <= 3; i++ {
        go processaComCancelamento(ctx, i)
    }

    time.Sleep(500 * time.Millisecond)
    cancel()
    time.Sleep(200 * time.Millisecond) // aguarda goroutines finalizarem
    fmt.Println("Programa finalizado")
}
```
</details>

### 4. Worker pool simples

Implemente um worker pool que processa 20 jobs com 4 workers. Cada job é um número
inteiro. O worker calcula o quadrado e envia o resultado.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        fmt.Printf("Worker %d: job %d\n", id, j)
        results <- j * j
    }
}

func main() {
    const totalJobs = 20
    const numWorkers = 4

    jobs := make(chan int, totalJobs)
    results := make(chan int, totalJobs)

    var wg sync.WaitGroup
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for j := 1; j <= totalJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()
    close(results)

    for r := range results {
        fmt.Println("Resultado:", r)
    }
}
```
</details>

---

**Próxima fase:** [Fase 05 — Testes, Web & Ferramentas](fase-05-testes-web.md)
