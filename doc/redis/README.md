# Redis na Pokedex Platform

Guia de referência sobre Redis e sua integração com Go, abrangendo cache, streaming, mensageria e operações distribuídas. Os conceitos apresentados conectam-se diretamente à arquitetura hexagonal da Pokedex Platform.

## Estrutura

```
doc/redis/
├── README.md                   ← este arquivo
├── cache/                      ← estratégias de cache com Redis + Go
│   ├── README.md
│   ├── 01-conceitos.md
│   ├── 02-estrategias.md
│   └── 03-integracao-go.md
├── streaming/                  ← mensageria: streams, pub/sub, pipelines
│   ├── 01-redis-streams.md
│   ├── 02-pubsub.md
│   └── 03-pipelines-lua.md
└── operacoes/                  ← arquitetura distribuída e produção
    └── 01-cluster-sentinel.md
```

## Referências

| Fonte | Conteúdo |
|-------|----------|
| [fidelissauro.dev/caching](https://fidelissauro.dev/caching/) | Conceitos e estratégias de cache sob a ótica de System Design |
| [redis.io/docs — streaming/go](https://redis.io/docs/latest/develop/use-cases/streaming/go/) | Implementação de Redis Streams com go-redis |
| Educative — Building Practical Applications with Redis Using Go | Tipos de dados, Pub/Sub, pipelines, Lua, Cluster, Sentinel |

## Contexto no projeto

O Redis já está presente na Pokedex Platform:

- **Docker Compose:** `core/docker-compose.yml` provisiona Redis 7 Alpine
- **Rate limiter:** `core/bff/mobile-bff/internal/adapters/inbound/http/rate_limiter.go` implementa sliding window com Redis Sorted Sets
- **Middleware:** `core/bff/mobile-bff/internal/adapters/inbound/http/middleware.go` instancia o cliente Redis via `REDIS_URL`
- **Config:** campo `RedisURL` declarado em `core/bff/mobile-bff/internal/config/config.go`
- **Dependência:** `github.com/redis/go-redis/v9` no `go.mod`

### AWS ElastiCache e outros provedores

O `go-redis` conecta-se com **qualquer servidor compatível com Redis** apenas trocando a URL de conexão. Não é necessário escrever um adapter específico por provedor. O que muda entre ambientes é a string de conexão:

| Ambiente | Exemplo de URL |
|----------|----------------|
| Local (Docker) | `redis://localhost:6379` |
| AWS ElastiCache | `rediss://master.xxxxx.use1.cache.amazonaws.com:6379` |
| Upstash | `rediss://:password@host:port` |

### Caminho para adapter hexagonal

Atualmente, o Redis é usado diretamente no middleware, sem passar por uma interface de port. O caminho para alinhar com a arquitetura hexagonal é:

1. Definir uma interface `CacheRepository` em `ports/outbound/`
2. Criar um adapter em `adapters/outbound/redis/` que implementa a interface
3. Injetar a dependência via `main.go` (como já é feito com PostgreSQL e clients HTTP)
4. O middleware passa a depender da interface, não do `*redis.Client` concreto

Veja `cache/03-integracao-go.md` para o detalhamento completo.
