# Integração com Go

Guia prático de instalação, conexão e uso do `go-redis` na Pokedex Platform, seguindo o padrão de Ports & Adapters (Hexagonal).

## Instalação

```bash
go get github.com/redis/go-redis/v9
```

A dependência já está presente no `go.mod` do `mobile-bff`:

```
github.com/redis/go-redis/v9 v9.21.0
```

## Conexão

### Local (Docker)

```go
import "github.com/redis/go-redis/v9"

rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})
```

### Via URL (recomendado para flexibilidade)

```go
opts, err := redis.ParseURL("redis://localhost:6379")
if err != nil {
    log.Fatal(err)
}
rdb := redis.NewClient(opts)
```

### AWS ElastiCache

```go
// ElastiCache com TLS
opts, err := redis.ParseURL("rediss://master.xxxxx.use1.cache.amazonaws.com:6379")
```

### Upstash Redis

```go
opts, err := redis.ParseURL("rediss://:password@host:port")
```

### Verificação de saúde

```go
ctx := context.Background()
if err := rdb.Ping(ctx).Err(); err != nil {
    log.Fatalf("redis indisponível: %v", err)
}
```

## Data Types e Comandos

### String

```go
// SET com TTL
rdb.Set(ctx, "pokemon:25", `{"name":"Pikachu"}`, 10*time.Minute)

// GET
val, err := rdb.Get(ctx, "pokemon:25").Result()
if err == redis.Nil {
    // chave não existe
}

// INCR — contador atômico
count, _ := rdb.Incr(ctx, "api:requests").Result()

// DEL
rdb.Del(ctx, "pokemon:25")
```

### Hash

```go
// HMSet
rdb.HSet(ctx, "pokemon:25", map[string]any{
    "name": "Pikachu",
    "type": "Electric",
    "hp":   35,
})

// HGet / HGetAll
name, _ := rdb.HGet(ctx, "pokemon:25", "name").Result()
all, _ := rdb.HGetAll(ctx, "pokemon:25").Result()
```

### Set

```go
// SAdd
rdb.SAdd(ctx, "favorites:user:123", "pokemon:25", "pokemon:6")

// SIsMember
isFav, _ := rdb.SIsMember(ctx, "favorites:user:123", "pokemon:25").Result()

// SMembers
allFavs, _ := rdb.SMembers(ctx, "favorites:user:123").Result()
```

### Sorted Set

```go
// ZAdd — leaderboard
rdb.ZAdd(ctx, "leaderboard:views", redis.Z{Score: 1500, Member: "Pikachu"})
rdb.ZAdd(ctx, "leaderboard:views", redis.Z{Score: 1200, Member: "Charizard"})

// ZRevRangeWithScores — top 10
top, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard:views", 0, 9).Result()
```

## Padrão de Adapter Hexagonal

### Situação atual

O Redis é usado diretamente no middleware, sem interface:

```go
// middleware.go — abordagem atual (acoplamento direto)
client := redis.NewClient(opts)
limiter := newRedisRateLimiter(client, maxRequests, window)
```

### Caminho proposto: alinhar com a arquitetura

#### 1. Definir interface de port outbound

```go
// core/bff/mobile-bff/internal/ports/outbound/cache_repository.go
package outbound

import (
    "context"
    "time"
)

type CacheRepository interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value any, ttl time.Duration) error
    Del(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Allow(ctx context.Context, key string, maxRequests int, window time.Duration) (bool, error)
}
```

#### 2. Implementar adapter Redis

```go
// core/bff/mobile-bff/internal/adapters/outbound/redis/cache_repository.go
package redis

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
    "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

type CacheRepository struct {
    client *redis.Client
}

func NewCacheRepository(redisURL string) (*CacheRepository, error) {
    opts, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, fmt.Errorf("redis URL inválida: %w", err)
    }
    client := redis.NewClient(opts)
    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, fmt.Errorf("redis indisponível: %w", err)
    }
    return &CacheRepository{client: client}, nil
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
    val, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", nil // chave não existe não é erro de sistema
    }
    return val, err
}

func (r *CacheRepository) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
    return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *CacheRepository) Del(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}

func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
    n, err := r.client.Exists(ctx, key).Result()
    return n > 0, err
}

func (r *CacheRepository) Allow(ctx context.Context, key string, maxRequests int, window time.Duration) (bool, error) {
    // Implementação do sliding window rate limiter existente,
    // agora encapsulada no adapter
    // ...
}

// Compile-time interface check
var _ outbound.CacheRepository = (*CacheRepository)(nil)
```

#### 3. Implementar adapter in-memory (fallback)

```go
// core/bff/mobile-bff/internal/adapters/outbound/memory/cache_repository.go
package memory

import (
    "context"
    "sync"
    "time"

    "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)

type CacheRepository struct {
    mu   sync.RWMutex
    data map[string]cacheEntry
}

type cacheEntry struct {
    value   string
    expires time.Time
}

func NewCacheRepository() *CacheRepository {
    return &CacheRepository{data: make(map[string]cacheEntry)}
}

// ... implementações de Get, Set, Del, Exists, Allow ...

var _ outbound.CacheRepository = (*CacheRepository)(nil)
```

#### 4. Configuração via Config struct

```go
// config.go — campo já existente
type Config struct {
    // ...
    RedisURL string // lido de REDIS_URL
}
```

#### 5. Wiring no main.go

```go
// main.go
var cacheRepo outbound.CacheRepository

if cfg.RedisURL != "" {
    cacheRepo, err = redis.NewCacheRepository(cfg.RedisURL)
    if err != nil {
        slog.Warn("redis indisponível, usando cache in-memory", "error", err)
    }
}
if cacheRepo == nil {
    cacheRepo = memory.NewCacheRepository()
}

// Inject no middleware
authMiddleware := middleware.NewAuthMiddleware(authService, cacheRepo)
```

## Pipeline

Redis pipelines agrupam múltiplos comandos em uma única chamada de rede, reduzindo round-trips. Já usado no rate limiter atual:

```go
pipe := rdb.Pipeline()
pipe.ZRemRangeByScore(ctx, key, "0", formatMillis(windowStart))
pipe.ZCard(ctx, key)
pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
pipe.Expire(ctx, key, window)
cmds, err := pipe.Exec(ctx)
```

## Boas práticas

1. **Fail-open no cache:** se o Redis falhar, a aplicação continua funcionando (vai ao banco direto), logando o erro
2. **TTL sempre:** nunca crie chaves sem expiração em cache
3. **Prefixos para namespacing:** `pokemon:`, `user:`, `ratelimit:` evitam colisões
4. **Nunca use `KEYS *` em produção:** prefira `SCAN` para iterar chaves
5. **Conexão com pool:** o `go-redis` gerencia pool automaticamente; ajuste `PoolSize` se necessário
6. **Context propagation:** sempre passe `context.Context` para suportar cancelamento e tracing

---

> Referências:
> - [redis.io/docs — go-redis](https://redis.io/docs/latest/develop/clients/go-redis/)
> - [go-redis GitHub](https://github.com/redis/go-redis)
> - [fidelissauro.dev/caching](https://fidelissauro.dev/caching/)
