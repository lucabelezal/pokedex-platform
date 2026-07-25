# Redis Streams

Redis Streams é uma estrutura de dados append-only que implementa um log de eventos durável com suporte a consumer groups, processamento at-least-once e recuperação de falhas.

## Conceitos fundamentais

Um stream é um log ordenado por tempo. Cada entrada tem um ID auto-gerado no formato `{milliseconds}-{sequence}` e carrega pares field/value.

```
demo:events:orders
  1716998413541-0   type=order.placed   order_id=o-1234  customer=alice  amount=49.50
  1716998413542-0   type=order.paid     order_id=o-1234  customer=alice  amount=49.50
  1716998413542-1   type=order.shipped  order_id=o-1235  customer=bob    amount=12.00
```

### Diferenciais vs. Pub/Sub

| Característica | Streams | Pub/Sub |
|---------------|---------|---------|
| Persistência | Mensagens retidas após consumo | Efêmero — perdido se não houver subscriber |
| Replay | Suportado via `XRANGE` | Não suportado |
| Consumer groups | Sim, com PEL e XACK | Não |
| At-least-once | Sim, via XACK + XAUTOCLAIM | At-most-once |
| Ordenação | Por ID de entrada | Por ordem de chegada |

---

## Comandos básicos

### XADD — produzir evento

```go
func produceEvent(ctx context.Context, rdb *redis.Client) {
    id, err := rdb.XAdd(ctx, &redis.XAddArgs{
        Stream: "demo:events:orders",
        MaxLen: 2000,        // retenção aproximada
        Approx: true,
        Values: map[string]any{
            "type":     "order.placed",
            "order_id": "o-1234",
            "customer": "alice",
            "amount":   "49.50",
        },
    }).Result()
    if err != nil {
        log.Printf("erro no XADD: %v", err)
    }
    fmt.Printf("evento produzido: %s\n", id)
}
```

### XADD em lote com pipeline

```go
func produceBatch(ctx context.Context, rdb *redis.Client, events []map[string]any) ([]string, error) {
    pipe := rdb.Pipeline()
    cmds := make([]*redis.StringCmd, len(events))

    for i, ev := range events {
        cmds[i] = pipe.XAdd(ctx, &redis.XAddArgs{
            Stream: "demo:events:orders",
            MaxLen: 2000,
            Approx: true,
            Values: ev,
        })
    }

    if _, err := pipe.Exec(ctx); err != nil {
        return nil, err
    }

    ids := make([]string, len(cmds))
    for i, cmd := range cmds {
        ids[i] = cmd.Val()
    }
    return ids, nil
}
```

### MAXLEN ~ (aproximado)

O `~` (Approx) permite que o Redis faça trimming em blocos (macro-nodes), muito mais eficiente que o trimming exato. Com MAXLEN ~ 50 e 300 eventos produzidos, o stream pode ficar com ~100 entradas — o Redis liberou o macro-node mais antigo inteiro e parou.

Se precisar de limite exato (raro), remova `Approx: true`. A diferença de performance é significativa em streams com alto volume.

---

## Consumer Groups

Consumer groups permitem que múltiplos consumidores independentes leiam o mesmo stream, cada um no seu ritmo, com seu próprio cursor e PEL (Pending Entries List).

### Criar consumer group

```go
// Cria o grupo "notifications" começando do início do stream
rdb.XGroupCreate(ctx, "demo:events:orders", "notifications", "0-0").Err()

// Ou a partir de agora ($), ignorando histórico
rdb.XGroupCreate(ctx, "demo:events:orders", "notifications", "$").Err()

// XGroupCreateMkStream cria o stream se não existir
rdb.XGroupCreateMkStream(ctx, "demo:events:orders", "notifications", "0-0").Err()
```

### XREADGROUP — consumir eventos

O ID especial `>` significa "entregue entradas que este grupo ainda não entregou a ninguém":

```go
func consumeEvents(ctx context.Context, rdb *redis.Client) {
    for {
        streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    "notifications",
            Consumer: "worker-a",
            Streams:  []string{"demo:events:orders", ">"},
            Count:    10,
            Block:    500 * time.Millisecond, // bloqueia até 500ms esperando eventos
        }).Result()

        if err == redis.Nil {
            continue // timeout, nenhum evento novo
        }
        if err != nil {
            log.Printf("erro no XREADGROUP: %v", err)
            time.Sleep(time.Second)
            continue
        }

        for _, stream := range streams {
            for _, msg := range stream.Messages {
                processEvent(ctx, rdb, msg)
            }
        }
    }
}
```

### XACK — confirmar processamento

```go
func processEvent(ctx context.Context, rdb *redis.Client, msg redis.XMessage) {
    // Processa o evento (envia notificação, atualiza dashboard, etc.)
    handleEvent(msg.Values)

    // Confirma o processamento — remove da PEL
    _, err := rdb.XAck(ctx, "demo:events:orders", "notifications", msg.ID).Result()
    if err != nil {
        log.Printf("erro no XACK %s: %v", msg.ID, err)
    }
}
```

**Crucial:** um evento nunca "ackado" permanece na PEL até ser reivindicado por XAUTOCLAIM. Isso garante entrega at-least-once.

---

## XAUTOCLAIM — recuperação de consumidores crashados

Se um consumidor morre sem dar XACK, seus eventos ficam na PEL. O XAUTOCLAIM permite que outro consumidor saudável resgate esses eventos.

```go
func reapIdlePel(ctx context.Context, rdb *redis.Client) {
    var cursor string = "0-0"
    for {
        claimed, nextCursor, err := rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
            Stream:   "demo:events:orders",
            Group:    "notifications",
            Consumer: "worker-b",  // consumidor saudável que receberá os eventos
            MinIdle:  5 * time.Second,
            Start:    cursor,
            Count:    10,
        }).Result()
        if err != nil {
            log.Printf("erro no XAUTOCLAIM: %v", err)
            return
        }

        for _, msg := range claimed {
            processEvent(ctx, rdb, msg)
        }

        if nextCursor == "0-0" {
            break
        }
        cursor = nextCursor
    }
}
```

**Nota:** `go-redis` v9 descarta o terceiro elemento do retorno de `XAUTOCLAIM` (lista de IDs deletados). Se precisar dessa informação (Redis 7.0+), use `rdb.Do(ctx, "XAUTOCLAIM", ...)` para parse manual.

---

## Consumidor como goroutine (loop completo)

```go
type ConsumerWorker struct {
    rdb    *redis.Client
    stream string
    group  string
    name   string
    cancel context.CancelFunc
}

func (w *ConsumerWorker) Start(ctx context.Context) {
    ctx, w.cancel = context.WithCancel(ctx)
    go w.run(ctx)
}

func (w *ConsumerWorker) run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }

        streams, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    w.group,
            Consumer: w.name,
            Streams:  []string{w.stream, ">"},
            Count:    10,
            Block:    500 * time.Millisecond,
        }).Result()

        if err == redis.Nil || errors.Is(err, context.Canceled) {
            continue
        }
        if err != nil {
            log.Printf("[%s/%s] erro: %v", w.group, w.name, err)
            time.Sleep(time.Second)
            continue
        }

        for _, stream := range streams {
            for _, msg := range stream.Messages {
                w.handleEntry(ctx, msg)
            }
        }
    }
}

func (w *ConsumerWorker) handleEntry(ctx context.Context, msg redis.XMessage) {
    // Lógica de negócio
    log.Printf("[%s/%s] processando %s: %v", w.group, w.name, msg.ID, msg.Values)

    // XACK
    w.rdb.XAck(ctx, w.stream, w.group, msg.ID)
}

func (w *ConsumerWorker) Stop() {
    w.cancel()
}
```

---

## Replay com XRANGE

XRANGE lê uma fatia do histórico sem afetar o estado de nenhum consumer group:

```go
// Ler todos os eventos do stream
msgs, _ := rdb.XRange(ctx, "demo:events:orders", "-", "+").Result()

// Ler os últimos 50 eventos
msgs, _ = rdb.XRangeN(ctx, "demo:events:orders", "-", "+", 50).Result()

// Ler eventos a partir de um timestamp específico
msgs, _ = rdb.XRange(ctx, "demo:events:orders", "1716998413541", "+").Result()
```

## XTRIM — retenção explícita

```go
// Manter apenas os ~1000 eventos mais recentes (aproximado)
rdb.XTrimMaxLenApprox(ctx, "demo:events:orders", 1000, 0).Result()

// Manter eventos mais novos que um ID específico
rdb.XTrimMinID(ctx, "demo:events:orders", "1716998413541-0").Result()
```

## Observabilidade

```go
// Tamanho do stream
length, _ := rdb.XLen(ctx, "demo:events:orders").Result()

// Info do stream
info, _ := rdb.XInfoStream(ctx, "demo:events:orders").Result()

// Info dos grupos
groups, _ := rdb.XInfoGroups(ctx, "demo:events:orders").Result()

// Info dos consumidores de um grupo
consumers, _ := rdb.XInfoConsumers(ctx, "demo:events:orders", "notifications").Result()

// Pending entries com idle time e delivery count
pending, _ := rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
    Stream: "demo:events:orders",
    Group:  "notifications",
    Start:  "-",
    End:    "+",
    Count:  20,
}).Result()
```

## Boas práticas com Streams

1. **Idempotência:** XAUTOCLAIM pode reentregar eventos — seu processamento deve ser idempotente
2. **Poison pill detection:** se um evento foi entregue 5+ vezes sem sucesso, mova para uma dead-letter stream
3. **Não use o mesmo consumer name após restart:** reaproveite o nome para drenar a PEL com `XREADGROUP 0` antes de ler novos eventos (`>`)
4. **Shutdown via context.Context:** cancele o contexto para desbloquear `XREADGROUP` com Block
5. **Particiamento:** em Redis Cluster, uma stream é uma única key em um único shard. Para escala, particione por tenant (`events:orders:{tenant_a}`)
6. **Métrica de lag:** monitore `XINFO GROUPS` para detectar grupos atrasados
7. **Não use XREAD (sem group) para at-least-once:** XREAD não tem PEL, não tem XACK — use XREADGROUP

---

> Referência: [redis.io/docs/latest/develop/use-cases/streaming/go/](https://redis.io/docs/latest/develop/use-cases/streaming/go/)
