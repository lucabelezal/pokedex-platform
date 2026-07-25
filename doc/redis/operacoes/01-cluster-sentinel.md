# Redis Cluster, Sentinel e Provedores

Guia sobre arquiteturas distribuídas de Redis e integração com provedores cloud, especialmente AWS ElastiCache.

## Redis Cluster

Redis Cluster particiona (shard) dados automaticamente entre múltiplos nós. Cada nó é responsável por um subconjunto das keys (hash slots).

### Arquitetura

```
┌─────────────────────────────────────────────┐
│                 Redis Cluster                │
│                                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐      │
│  │ Node A  │  │ Node B  │  │ Node C  │      │
│  │ Primary │  │ Primary │  │ Primary │      │
│  │ 0-5460  │  │5461-10922│ │10923-   │      │
│  │          │  │          │  │  16383  │      │
│  └────┬────┘  └────┬────┘  └────┬────┘      │
│       │            │            │             │
│  ┌────┴────┐  ┌────┴────┐  ┌────┴────┐      │
│  │ Replica │  │ Replica │  │ Replica │      │
│  │  (B)    │  │  (C)    │  │  (A)    │      │
│  └─────────┘  └─────────┘  └─────────┘      │
└─────────────────────────────────────────────┘
```

- Mínimo de 3 nós primários
- Cada primário pode ter 1+ réplicas
- 16.384 hash slots distribuídos entre os primários
- Failover automático: se Node A cair, sua réplica (em outro nó) é promovida

### Conexão com go-redis

```go
rdb := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: []string{
        "node1:6379",
        "node2:6379",
        "node3:6379",
        "node4:6379",
        "node5:6379",
        "node6:6379",
    },
    Password: os.Getenv("REDIS_PASSWORD"),
})

// O go-redis automaticamente:
// - Descobre a topologia do cluster
// - Roteia comandos para o nó correto (baseado no hash slot da key)
// - Lida com MOVED/ASK redirections
// - Atualiza a topologia quando nós entram/saem
```

### Hash tags

Para garantir que keys relacionadas fiquem no mesmo nó (necessário para operações multi-key):

```go
// Keys com {tag} vão para o mesmo hash slot
rdb.Set(ctx, "user:{123}:name", "Alice", 0)
rdb.Set(ctx, "user:{123}:email", "alice@example.com", 0)

// Agora é seguro usar em transações ou Lua scripts
rdb.MGet(ctx, "user:{123}:name", "user:{123}:email")
```

### Limitações

- Operações multi-key só funcionam se todas as keys estão no mesmo hash slot
- Lua scripts devem acessar keys do mesmo slot (use hash tags)
- Transações são limitadas ao escopo de um nó
- Pub/Sub é broadcast para todos os nós (não sharded)

---

## Redis Sentinel

Sentinel é um sistema de monitoramento e failover para Redis. Diferente do Cluster, **não particiona dados** — é alta disponibilidade para instâncias standalone ou master-replica.

### Arquitetura

```
┌──────────────────────────────┐
│        Sentinel Nodes        │
│   ┌──────┐ ┌──────┐ ┌──────┐│
│   │ S1   │ │ S2   │ │ S3   ││
│   └──┬───┘ └──┬───┘ └──┬───┘│
│      │        │        │     │
│      └────────┼────────┘     │
│               │ monitoram    │
└───────────────┼──────────────┘
                │
    ┌───────────┴───────────┐
    │     Redis Servers      │
    │  ┌────────┐ ┌────────┐ │
    │  │ Master │ │Replica │ │
    │  └────────┘ └────────┘ │
    └────────────────────────┘
```

- Mínimo de 3 Sentinel nodes (quorum)
- Sentinel elege novo master se o atual falhar
- Cliente consulta Sentinel para descobrir qual é o master atual

### Conexão com go-redis

```go
rdb := redis.NewFailoverClient(&redis.FailoverOptions{
    MasterName:    "mymaster",
    SentinelAddrs: []string{
        "sentinel1:26379",
        "sentinel2:26379",
        "sentinel3:26379",
    },
    Password: os.Getenv("REDIS_PASSWORD"),
})

// go-redis automaticamente:
// - Consulta Sentinel pelo endereço do master atual
// - Reconecta em caso de failover
// - Gerencia troca de master transparentemente
```

### Cluster vs. Sentinel

| Critério | Cluster | Sentinel |
|----------|---------|----------|
| Particionamento (sharding) | Sim (automático) | Não |
| Alta disponibilidade | Sim (failover automático) | Sim (failover automático) |
| Escala horizontal | Sim (adiciona nós) | Apenas leitura com réplicas |
| Complexidade | Alta | Média |
| Mínimo de nós | 6 (3 primary + 3 replica) | 5 (3 sentinel + 2 redis) |
| Uso típico | Grandes volumes de dados | HA para workloads moderados |

---

## AWS ElastiCache

AWS ElastiCache é o serviço gerenciado de Redis da AWS. Suporta tanto Cluster Mode quanto Sentinel (chamado de "Replication Group").

### Modos disponíveis

| Modo ElastiCache | Equivalente Redis | Quando usar |
|------------------|-------------------|-------------|
| Serverless | Gerenciado automático | Cargas variáveis, sem administração |
| Cluster Mode Enabled | Redis Cluster | Escala horizontal, sharding |
| Cluster Mode Disabled | Redis Sentinel | HA simples, sem sharding |

### Conexão via go-redis

**Cluster Mode Enabled:**

```go
rdb := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: []string{
        os.Getenv("REDIS_PRIMARY_ENDPOINT"),
    },
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
    },
})
```

**Cluster Mode Disabled (com TLS):**

```go
rdb := redis.NewClient(&redis.Options{
    Addr:      os.Getenv("REDIS_PRIMARY_ENDPOINT"),
    TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
})
```

### Configuração de segurança no ElastiCache

- **Encryption in transit:** obrigatório para conexões TLS. Use `rediss://` ou configure `TLSConfig`
- **Encryption at rest:** habilitado no lado AWS, transparente para o cliente
- **Auth token:** senha obrigatória, configurada no parâmetro `Password`
- **Security Groups:** controle de acesso por IP/security group

### Autenticação com IAM (ElastiCache Serverless)

Para ElastiCache Serverless com autenticação IAM:

```go
import "github.com/aws/aws-sdk-go-v2/service/elasticache"

// Gera token de autenticação IAM
func generateAuthToken(endpoint, username string) (string, error) {
    // O token IAM é válido por 15 minutos
    // Use AWS SDK para gerar
    // ...
}

rdb := redis.NewClient(&redis.Options{
    Addr:      endpoint,
    Username:  username,
    Password:  token, // rotacionar a cada 15 min
    TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
})
```

---

## Client-side vs. Server-side Proxy

### Client-side proxy (go-redis)

O próprio go-redis implementa a lógica de descoberta e roteamento:
- Cluster: descobre topologia, roteia por hash slot
- Sentinel: consulta sentinel, descobre master
- Ring: anel de hash consistente para sharding manual

**Vantagem:** sem ponto único de falha adicional.

### Server-side proxy (Twemproxy, Envoy)

Um proxy dedicado fica entre a aplicação e os servidores Redis:
- Aplicação conecta no proxy como se fosse um Redis único
- Proxy roteia comandos para os nós corretos

**Vantagem:** clientes mais simples, compatível com qualquer linguagem.

---

## Docker Compose (Pokedex Platform)

O Redis local da Pokedex Platform está em `core/docker-compose.yml`:

```yaml
redis:
  image: redis:7-alpine
  command: ["redis-server", "/usr/local/etc/redis/redis.conf"]
  ports:
    - "${REDIS_PORT:-6379}:6379"
  volumes:
    - ./infra/redis/redis.conf:/usr/local/etc/redis/redis.conf:ro
  healthcheck:
    test: ["CMD", "redis-cli", "ping"]
```

Para usar um ElastiCache em produção, basta:
1. Remover o serviço `redis` do compose de produção
2. Apontar `REDIS_URL` para o endpoint do ElastiCache
3. Configurar TLS se necessário

Nenhuma alteração de código é necessária — o `go-redis` lida com a diferença transparentemente via `redis.ParseURL()`.

---

> Referência: Educative — Building Practical Applications with Redis Using Go
