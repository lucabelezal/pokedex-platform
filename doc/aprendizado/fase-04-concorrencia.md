# Fase 04 — Concorrência

**Aulas do curso:** 72 a 84

## Sumário

| Passo | Recurso | Tempo est. |
|-------|---------|-----------|
| 1 | Assistir aulas do curso | ~2h30 |
| 2 | Executar exemplos Go by Example | ~40min |
| 3 | Ler Effective Go | ~30min |
| 4 | Conferir roadmap.sh | ~5min |
| 5 | Ler styleguide | ~10min |
| 6 | Fazer exercício prático | ~20min |

---

## 1. Aulas do curso

| # | Aula | Conceito |
|---|------|----------|
| 72 | Concorrência vs Paralelismo | Diferença conceitual |
| 73 | Curiosidade: Número de CPUs | `runtime.NumCPU()` |
| 74 | Conhecendo a Goroutine | `go fn()` — execução concorrente |
| 75 | Conhecendo o Channel (Canal) | `ch := make(chan T)` — comunicação entre goroutines |
| 76 | Usando Goroutine e Channel | Envio `ch <- v`, recebimento `v := <-ch` |
| 77 | Cuidado com os Deadlocks | Goroutine esperando forever |
| 78 | Channel com Buffer | `make(chan T, n)` — buffer não bloqueia até encher |
| 79 | Channel: Usando Range e Close | `for v := range ch`, `close(ch)` |
| 80 | Padrão de Concorrência: Generators | Função que retorna channel |
| 81 | Criando um Pacote Reutilizável | — |
| 82 | Padrão de Concorrência: Multiplexador | Combinar múltiplos channels |
| 83 | Estrutura de Controle: Select | `select { case <-ch1: ... case <-ch2: ... }` |
| 84 | Multiplexador com Select | Fan-in com select |

---

## 2. Go by Example — exemplos desta fase

| Exemplo | Link | O que observar |
|---------|------|---------------|
| Goroutines | https://gobyexample.com/goroutines | `go fn()` — leve, milhares simultâneas |
| Channels | https://gobyexample.com/channels | `make(chan string)`, envia e recebe — sincronização implícita |
| Channel Buffering | https://gobyexample.com/channel-buffering | `make(chan string, 2)` — não bloqueia até encher |
| Channel Synchronization | https://gobyexample.com/channel-synchronization | Usar channel como sinalizador de conclusão |
| Channel Directions | https://gobyexample.com/channel-directions | `chan<-` (send-only), `<-chan` (receive-only) |
| Select | https://gobyexample.com/select | Espera múltiplas operações de channel |
| Timeouts | https://gobyexample.com/timeouts | `time.After` com `select` — timeout pattern |
| Non-Blocking Channel Ops | https://gobyexample.com/non-blocking-channel-operations | `select` com `default` — tenta, não bloqueia |
| Closing Channels | https://gobyexample.com/closing-channels | `close(ch)` — sinaliza "acabou" |
| Range over Channels | https://gobyexample.com/range-over-channels | Itera até o channel ser fechado |
| Timers | https://gobyexample.com/timers | `time.NewTimer` — dispara uma vez no futuro |
| Tickers | https://gobyexample.com/tickers | `time.NewTicker` — dispara repetidamente em intervalo |
| Worker Pools | https://gobyexample.com/worker-pools | Múltiplas goroutines processando de um channel |
| WaitGroups | https://gobyexample.com/waitgroups | `sync.WaitGroup` — aguarda grupo de goroutines |
| Rate Limiting | https://gobyexample.com/rate-limiting | Ticker como limitador de taxa |
| Atomic Counters | https://gobyexample.com/atomic-counters | `sync/atomic` — contadores atômicos sem mutex |
| Mutexes | https://gobyexample.com/mutexes | `sync.Mutex` — exclusão mútua |
| Stateful Goroutines | https://gobyexample.com/stateful-goroutines | Encapsular estado em uma goroutine dona |
| Context | https://gobyexample.com/context | `context.WithCancel`, `WithTimeout`, `WithValue` |

---

## 3. Effective Go — seções para ler

> Leia a seção inteira de Concorrência. É a mais importante do Effective Go.

| Seção | Link | O que aprender |
|-------|------|---------------|
| **Share by communicating** | https://go.dev/doc/effective_go#sharing | "Do not communicate by sharing memory; share memory by communicating" |
| **Goroutines** | https://go.dev/doc/effective_go#goroutines | Leves, stacks pequenas, multiplexadas em threads OS |
| **Channels** | https://go.dev/doc/effective_go#channels | Sincronização + comunicação; unbuffered = sincronização, buffered = semáforo |
| **Channels of channels** | https://go.dev/doc/effective_go#chan_of_chan | Channels como first-class values — RPC pattern |
| **Parallelization** | https://go.dev/doc/effective_go#parallelization | Dividir trabalho entre CPUs, `runtime.NumCPU` |
| **A leaky buffer** | https://go.dev/doc/effective_go#leaky_buffer | Exemplo canônico: free list com channel bufferizado |

---

## 4. Roadmap.sh — tópicos desta fase

| Categoria | Tópicos |
|-----------|---------|
| Goroutines | Basics, lightweight threads |
| Channels | Buffered vs Unbuffered, Select Statement, Worker Pools |
| `sync` Package | Mutexes, WaitGroups |
| `context` Package | Deadlines & Cancellations, Common Usecases |
| Concurrency Patterns | fan-in, fan-out, pipeline |
| Race Detection | Race detector (`go test -race`) |

---

## 5. Guia de Estilo — regras desta fase

> Leia [`decisions.md#concorrência`](../guia-estilo/decisions.md#concorrência):

| Regra | Seção | Por que importa agora |
|-------|-------|----------------------|
| Aguardar goroutines finalizarem | [Aguardar goroutines](../guia-estilo/decisions.md#aguardar-goroutines-finalizarem) | Use `sync.WaitGroup` e `context.Context` para cancelamento |
| Sem goroutines em `init()` | [Sem goroutines em init()](../guia-estilo/decisions.md#sem-goroutines-em-init) | Torna o comportamento imprevisível |
| Canais: tamanho 1 ou 0 | [Canais](../guia-estilo/decisions.md#canais) | `make(chan T)` ou `make(chan T, 1)`. Sem tamanhos arbitrários. |
| Mutex zero value é válido | [Mutex com valor zero](../guia-estilo/decisions.md#mutex-com-valor-zero-é-válido) | Não use ponteiro para mutex sem necessidade |
| Evitar globais mutáveis | [Evitar globais mutáveis](../guia-estilo/decisions.md#evitar-globais-mutáveis) | Injeção de dependência, não estado global |
| Usar `time` para tempo | [Usar time](../guia-estilo/decisions.md#usar-time-para-manipular-tempo) | `time.Time` e `time.Duration`, nunca inteiros |

---

## 6. No código do projeto

### Circuit breaker com retry (concorrência segura)

```go
// Em internal/adapters/outbound/http/circuit_breaker.go:
type CircuitBreakerClient struct {
    cb    *gobreaker.CircuitBreaker[[]byte]
    retry int
}

func (c *CircuitBreakerClient) Do(req *http.Request) (*http.Response, error) {
    body, err := c.cb.Execute(func() ([]byte, error) {
        resp, err := c.client.Do(req)
        // ...
        return body, nil
    })
    if errors.Is(err, gobreaker.ErrOpenState) {
        return nil, domain.ErrServiceUnavailable
    }
    // ...
}
```

### Context propagation em todas as camadas

```go
// Em internal/service/pokemon_service.go — context é sempre o primeiro parâmetro:
func (s *PokemonService) List(ctx context.Context, params SearchParams) (*domain.PokemonPage, error) {
    pokemons, err := s.repo.List(ctx, params)
    // ...
}

// Em internal/adapters/inbound/http/pokemon_handler.go — timeout por requisição:
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()
```

### Middleware de rate limiting

```go
// Em internal/adapters/inbound/http/rate_limiter.go:
type rateLimiter interface {
    Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}
// Implementações: Redis (preferencial) e in-memory (fallback)
```

---

## 7. Exercício prático

**Objetivo:** Implementar um worker pool simples e comparar com o código do projeto.

1. Crie um arquivo de teste `doc/aprendizado/exercicios/worker_pool_test.go`:

```go
package exercicios

import (
    "sync"
    "testing"
)

func Process(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        // simula processamento
        results <- j * 2
    }
}

func TestWorkerPool(t *testing.T) {
    const numJobs = 10
    const numWorkers = 3

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go Process(jobs, results, &wg)
    }

    for j := 0; j < numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()
    close(results)

    count := 0
    for r := range results {
        t.Logf("resultado: %d", r)
        count++
    }
    if count != numJobs {
        t.Errorf("esperado %d resultados, obteve %d", numJobs, count)
    }
}
```

2. Execute `go test -race ./doc/aprendizado/exercicios/ -v`
3. Rode também com `go test -race ./...` no diretório `mobile-bff` para verificar se há data races no código do projeto
4. Compare o padrão de worker pool com o circuit breaker do projeto — ambos usam concorrência para resiliência

---

**Próxima fase:** [Fase 05 — Testes, Web & Deploy](fase-05-testes-web.md)
