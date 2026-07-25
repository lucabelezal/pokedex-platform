# Pipelines, Transações e Lua Scripts

## Pipelines

Pipelines agrupam múltiplos comandos Redis em uma única chamada de rede, reduzindo round-trips e melhorando significativamente a performance.

### Como funciona

Sem pipeline, cada comando é um request/response separado:

```
Cliente → SET key1 val1 → Redis → OK
Cliente → SET key2 val2 → Redis → OK
Cliente → INCR counter  → Redis → 1
```

Com pipeline, todos são enviados juntos:

```
Cliente → SET key1 val1
          SET key2 val2
          INCR counter   → Redis → OK, OK, 1
```

### Implementação em Go

```go
func pipelineExample(ctx context.Context, rdb *redis.Client) error {
    pipe := rdb.Pipeline()

    // Agrupa comandos — não são executados ainda
    set1 := pipe.Set(ctx, "pokemon:25:name", "Pikachu", 0)
    set2 := pipe.Set(ctx, "pokemon:25:type", "Electric", 0)
    incr := pipe.Incr(ctx, "stats:pokemon:updates")

    // Executa todos de uma vez
    _, err := pipe.Exec(ctx)
    if err != nil {
        return fmt.Errorf("pipeline falhou: %w", err)
    }

    // Resultados estão disponíveis nos comandos
    fmt.Printf("SET 1: %s\n", set1.Val())
    fmt.Printf("SET 2: %s\n", set2.Val())
    fmt.Printf("INCR: %d\n", incr.Val())

    return nil
}
```

### Uso no rate limiter da Pokedex Platform

O rate limiter atual já usa pipeline em `rate_limiter.go`:

```go
pipe := rdb.Pipeline()
pipe.ZRemRangeByScore(ctx, key, "0", formatMillis(windowStart))
pipe.ZCard(ctx, key)
pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
pipe.Expire(ctx, key, window)
cmds, err := pipe.Exec(ctx)
```

### Boas práticas

- Agrupe de 50 a 100 comandos por pipeline para balancear eficiência e uso de memória
- Pipelines não são transações — comandos podem ser intercalados com outros clientes
- Use `Pipelined()` para callbacks mais idiomáticos

```go
// Pipelined — estilo callback
_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
    pipe.Set(ctx, "key1", "val1", 0)
    pipe.Set(ctx, "key2", "val2", 0)
    return nil
})
```

---

## Transações (MULTI/EXEC)

Transações Redis garantem execução isolada de um grupo de comandos. Durante uma transação, outros clientes não são atendidos.

**Importante:** transações Redis **não são ACID** como em bancos relacionais. Não há rollback — se um comando falha, os demais continuam executando.

```go
func transactionExample(ctx context.Context, rdb *redis.Client) error {
    // TxPipeline combina pipeline + transação
    tx := rdb.TxPipeline()

    tx.Set(ctx, "pokemon:6:name", "Charizard", 0)
    tx.Incr(ctx, "stats:total:updates")
    tx.SAdd(ctx, "pokemon:updated:today", "6")

    _, err := tx.Exec(ctx)
    return err
}
```

### Padrão Watch (optimistic locking)

```go
func transferWithWatch(ctx context.Context, rdb *redis.Client, key string) error {
    const maxRetries = 3

    for i := 0; i < maxRetries; i++ {
        err := rdb.Watch(ctx, func(tx *redis.Tx) error {
            // Lê valor atual dentro do Watch
            val, err := tx.Get(ctx, key).Int()
            if err != nil && err != redis.Nil {
                return err
            }

            if val <= 0 {
                return fmt.Errorf("saldo insuficiente")
            }

            // Operação dentro da transação
            _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
                pipe.Decr(ctx, key)
                return nil
            })
            return err
        }, key)

        if err == nil {
            return nil
        }
        if errors.Is(err, redis.TxFailedErr) {
            continue // chave foi modificada por outro cliente, tenta de novo
        }
        return err
    }
    return fmt.Errorf("transação falhou após %d tentativas", maxRetries)
}
```

---

## Lua Scripts

Lua scripts são executados atomicamente no servidor Redis — similar a stored procedures. Nenhum outro cliente é atendido durante a execução.

### Vantagens

- **Atomicidade:** o script inteiro executa sem interrupção
- **Eficiência:** processamento acontece onde os dados estão (server-side)
- **Flexibilidade:** lógica condicional complexa dentro do Redis

### Estrutura básica

```go
// Lua script como string
var updateAndPublishScript = redis.NewScript(`
    local current = redis.call('GET', KEYS[1])
    if not current then
        current = '0'
    end

    local new_val = tonumber(current) + tonumber(ARGV[1])
    redis.call('SET', KEYS[1], new_val)
    redis.call('PUBLISH', KEYS[2], new_val)

    return new_val
`)

func runScript(ctx context.Context, rdb *redis.Client) {
    result, err := updateAndPublishScript.Run(ctx, rdb,
        []string{"counter:views", "channel:updates"}, // KEYS
        1, // ARGV[1]
    ).Int()
    if err != nil {
        log.Printf("script falhou: %v", err)
    }
    fmt.Printf("novo valor: %d\n", result)
}
```

### Exemplo: rate limiter como Lua script

O rate limiter atual usa pipeline com múltiplos comandos. Uma alternativa atômica com Lua script:

```go
var slidingWindowScript = redis.NewScript(`
    local key = KEYS[1]
    local now = tonumber(ARGV[1])
    local window = tonumber(ARGV[2])
    local max_requests = tonumber(ARGV[3])

    -- Remove entradas fora da janela
    redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

    -- Conta requisições dentro da janela
    local count = redis.call('ZCARD', key)
    if count >= max_requests then
        return 0  -- bloqueado
    end

    -- Adiciona esta requisição
    redis.call('ZADD', key, now, now .. '-' .. math.random())
    redis.call('EXPIRE', key, math.ceil(window / 1000))

    return 1  -- permitido
`)

func isAllowed(ctx context.Context, rdb *redis.Client, clientID string) (bool, error) {
    now := time.Now().UnixMilli()
    window := int64(time.Minute.Milliseconds())
    maxRequests := 100

    result, err := slidingWindowScript.Run(ctx, rdb,
        []string{"ratelimit:" + clientID},
        now, window, maxRequests,
    ).Int()
    if err != nil {
        return true, err // fail-open
    }
    return result == 1, nil
}
```

### Boas práticas com Lua

1. **Scripts determinísticos:** evite `math.random()` sem seed — o exemplo acima é didático, em produção use timestamp + contador
2. **Evite scripts longos:** Redis é single-threaded; scripts bloqueiam todos os outros clientes
3. **KEYS vs ARGV:** sempre separe chaves de argumentos para compatibilidade com Redis Cluster
4. **EVALSHA:** o go-redis faz cache automático do SHA. Use `NewScript()` + `Run()` em vez de `Eval()` para performance
5. **Teste antes:** valide scripts com `redis-cli --eval` antes de embedar no Go

---

## Comparativo: Pipeline vs. Transação vs. Lua

| Característica | Pipeline | Transação (MULTI/EXEC) | Lua Script |
|---------------|----------|------------------------|------------|
| Atomicidade | Não | Sim (isolamento, sem rollback) | Sim |
| Performance | Alta (batch de comandos) | Média (isolamento adiciona overhead) | Alta (server-side) |
| Lógica condicional | Não (client-side) | Limitada (WATCH) | Sim (Turing-complete) |
| Uso típico | Bulk writes, reads | Transferências, inventário | Rate limiting, validações complexas |
| Bloqueia outros clientes | Não | Durante EXEC | Durante toda execução |

---

> Referência: Educative — Building Practical Applications with Redis Using Go
