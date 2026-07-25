# Redis Pub/Sub

Redis Pub/Sub implementa o padrão Publish/Subscribe: produtores publicam mensagens em canais, e um ou mais consumidores (subscribers) as recebem em tempo real.

## Características

- **Fire-and-forget:** mensagens não são persistidas — se não houver subscriber ativo, a mensagem é perdida
- **At-most-once:** sem confirmação de entrega (sem equivalente ao XACK)
- **Baixa latência:** ideal para notificações em tempo real onde perda ocasional é aceitável
- **Broadcast:** todos os subscribers de um canal recebem a mesma mensagem

## Quando usar Pub/Sub vs. Streams

| Cenário | Escolha |
|---------|---------|
| Notificações em tempo real (chat, alerts) | Pub/Sub |
| Eventos que precisam ser persistidos e reprocessados | Streams |
| Múltiplos consumidores independentes processando o mesmo fluxo | Streams com consumer groups |
| Broadcast para todos os ouvintes ativos | Pub/Sub |
| Garantia de entrega (at-least-once) | Streams |

---

## Implementação em Go

### Publisher

```go
func publishEvent(ctx context.Context, rdb *redis.Client, channel string, message any) error {
    jsonData, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("erro ao serializar: %w", err)
    }

    return rdb.Publish(ctx, channel, jsonData).Err()
}

// Uso
err := publishEvent(ctx, rdb, "pokemon:updates", map[string]any{
    "event": "pokemon.favorited",
    "pokemon_id": 25,
    "user_id": "user-123",
})
```

### Subscriber

```go
func subscribeChannel(ctx context.Context, rdb *redis.Client, channel string) {
    pubsub := rdb.Subscribe(ctx, channel)
    defer pubsub.Close()

    ch := pubsub.Channel()

    for {
        select {
        case <-ctx.Done():
            return
        case msg, ok := <-ch:
            if !ok {
                return
            }
            handleMessage(msg)
        }
    }
}

func handleMessage(msg *redis.Message) {
    log.Printf("[%s] recebido: %s", msg.Channel, msg.Payload)

    var event map[string]any
    if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
        log.Printf("erro ao deserializar mensagem: %v", err)
        return
    }

    // Processa o evento (ex: atualizar WebSocket, enviar push notification)
}
```

### Pattern Subscriptions

Redis suporta pattern matching no subscribe:

```go
// Subscreve em todos os canais que começam com "pokemon:"
pubsub := rdb.PSubscribe(ctx, "pokemon:*")
defer pubsub.Close()

ch := pubsub.Channel()
for msg := range ch {
    fmt.Printf("canal: %s, pattern: %s, payload: %s\n",
        msg.Channel, msg.Pattern, msg.Payload)
}
```

### Exemplo completo: notificações da Pokedex

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "os/signal"

    "github.com/redis/go-redis/v9"
)

type PokemonEvent struct {
    Type      string `json:"type"`
    PokemonID int    `json:"pokemon_id"`
    UserID    string `json:"user_id,omitempty"`
}

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    // Subscriber: escuta notificações de favoritos
    go func() {
        pubsub := rdb.Subscribe(ctx, "pokemon:favorites")
        defer pubsub.Close()

        for msg := range pubsub.Channel() {
            var ev PokemonEvent
            json.Unmarshal([]byte(msg.Payload), &ev)
            log.Printf("pokemon %d favoritado por %s", ev.PokemonID, ev.UserID)
        }
    }()

    // Publisher: emite evento quando pokemon é favoritado
    func favoritePokemon(userID string, pokemonID int) {
        ev := PokemonEvent{
            Type:      "pokemon.favorited",
            PokemonID: pokemonID,
            UserID:    userID,
        }
        data, _ := json.Marshal(ev)
        rdb.Publish(ctx, "pokemon:favorites", data)
    }

    // Simula favoritar alguns pokemons
    favoritePokemon("user-1", 25)
    favoritePokemon("user-1", 6)

    <-ctx.Done()
    log.Println("encerrando...")
}
```

## Limitações

- **Sem persistência:** se o subscriber desconectar, perde as mensagens do período offline
- **Sem confirmação:** não há ACK — não é possível saber se a mensagem foi processada
- **Sem replay:** não é possível "voltar no tempo" e ler mensagens antigas
- **Fire-and-forget:** o publisher não sabe quantos (se algum) subscriber recebeu

Para cenários que exigem garantia de entrega, use **Redis Streams** (`streaming/01-redis-streams.md`).

---

> Referência: [redis.io/docs — Pub/Sub](https://redis.io/docs/latest/develop/interact/pubsub/)
