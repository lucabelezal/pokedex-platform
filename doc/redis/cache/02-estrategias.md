# Estratégias de Cache

Existem diversas formas de implementar cache entre a aplicação e o banco de dados. As três principais estratégias são Cache-Aside, Write-Through e Write-Behind.

## Cache-Aside (Lazy Loading)

É a estratégia mais comum. A aplicação gerencia o cache sob demanda — primeiro verifica o cache e, em caso de miss, consulta o banco e popula o cache.

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  Aplicação │───▶│   Cache   │───▶│  Database │
└──────────┘    └──────────┘    └──────────┘
     │               │                │
     │  1. GET       │                │
     │──────────────▶│                │
     │               │                │
     │  2. miss      │                │
     │◀──────────────│                │
     │               │                │
     │  3. SELECT                    │
     │──────────────────────────────▶│
     │               │                │
     │  4. result                    │
     │◀──────────────────────────────│
     │               │                │
     │  5. SET       │                │
     │──────────────▶│                │
```

### Implementação em Go

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

type Pokemon struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"`
}

func getPokemon(ctx context.Context, rdb *redis.Client, id int) (*Pokemon, error) {
    key := fmt.Sprintf("pokemon:%d", id)

    // 1. Tenta buscar no cache
    data, err := rdb.Get(ctx, key).Result()
    if err == nil {
        // Cache hit
        var p Pokemon
        json.Unmarshal([]byte(data), &p)
        return &p, nil
    }

    if err != redis.Nil {
        return nil, fmt.Errorf("erro ao consultar cache: %w", err)
    }

    // 2. Cache miss — busca no banco de dados
    p, err := findPokemonInDB(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("pokemon não encontrado: %w", err)
    }

    // 3. Popula o cache para consultas futuras
    jsonData, _ := json.Marshal(p)
    rdb.Set(ctx, key, jsonData, 10*time.Minute)

    return p, nil
}

func findPokemonInDB(ctx context.Context, id int) (*Pokemon, error) {
    // Simula consulta ao banco
    return &Pokemon{ID: id, Name: "Pikachu", Type: "Electric"}, nil
}
```

### Vantagens e desvantagens

| Vantagens | Desvantagens |
|-----------|--------------|
| Simples de implementar | Primeira consulta sempre é miss (cold start) |
| Cache só contém o que é usado | Complexidade na gestão de consistência |
| Resiliência: se o cache falha, vai ao banco | Dados podem ficar desatualizados entre escritas |

---

## Write-Through (Escrita Dupla)

Os dados são escritos simultaneamente no cache e no banco de dados. O cache está sempre sincronizado com a fonte persistente.

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  Aplicação │───▶│   Cache   │    │  Database │
└──────────┘    └──────────┘    └──────────┘
     │               │                │
     │  1. SET + INSERT              │
     │──────────────▶│                │
     │               │───────────────▶│
     │               │                │
     │  2. OK        │  2. OK         │
     │◀──────────────│◀───────────────│
```

### Implementação em Go

```go
func savePokemon(ctx context.Context, rdb *redis.Client, db *sql.DB, p *Pokemon) error {
    // 1. Escreve no banco de dados
    query := "INSERT INTO pokemons (id, name, type) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET name = $2, type = $3"
    _, err := db.ExecContext(ctx, query, p.ID, p.Name, p.Type)
    if err != nil {
        return fmt.Errorf("erro ao salvar no banco: %w", err)
    }

    // 2. Atualiza o cache imediatamente
    key := fmt.Sprintf("pokemon:%d", p.ID)
    jsonData, _ := json.Marshal(p)
    if err := rdb.Set(ctx, key, jsonData, 10*time.Minute).Err(); err != nil {
        // Loga o erro, mas não falha a operação — o banco já foi atualizado
        log.Printf("erro ao atualizar cache para %s: %v", key, err)
    }

    return nil
}
```

### Vantagens e desvantagens

| Vantagens | Desvantagens |
|-----------|--------------|
| Cache sempre atualizado | Latência de escrita maior (duas operações) |
| Minimiza risco de inconsistência | Cache pode conter dados nunca lidos |
| Leitura rápida garantida | Falha no cache após escrita no banco exige rollback ou compensação |

---

## Write-Behind (Lazy Writing)

As escritas vão primeiro para o cache e são propagadas ao banco de forma assíncrona. Ideal para cenários com alta taxa de escrita.

```
┌──────────┐    ┌──────────┐         ┌──────────┐
│  Aplicação │───▶│   Cache   │  ...  ▶│  Database │
└──────────┘    └──────────┘  async  └──────────┘
     │               │
     │  1. SET       │
     │──────────────▶│
     │               │
     │  2. OK        │
     │◀──────────────│
     │               │  (depois, em lote ou agendado)
     │               │──────────────────▶│
```

### Implementação em Go (com worker de sincronização)

```go
type WriteBehindCache struct {
    rdb  *redis.Client
    db   *sql.DB
    mu   sync.Mutex
    dirty map[string]time.Time // chaves pendentes de sincronização
}

func (c *WriteBehindCache) Set(ctx context.Context, key string, value any) error {
    jsonData, _ := json.Marshal(value)
    if err := c.rdb.Set(ctx, key, jsonData, 0).Err(); err != nil {
        return err
    }

    c.mu.Lock()
    c.dirty[key] = time.Now()
    c.mu.Unlock()

    return nil
}

// Worker que executa em background, sincronizando dirty keys com o banco
func (c *WriteBehindCache) SyncWorker(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            c.flushDirty(ctx)
        }
    }
}

func (c *WriteBehindCache) flushDirty(ctx context.Context) {
    c.mu.Lock()
    keys := make([]string, 0, len(c.dirty))
    for k := range c.dirty {
        keys = append(keys, k)
    }
    c.dirty = make(map[string]time.Time)
    c.mu.Unlock()

    for _, key := range keys {
        data, err := c.rdb.Get(ctx, key).Result()
        if err != nil {
            continue
        }
        // Persiste no banco de dados
        persistToDB(ctx, c.db, key, data)
    }
}
```

### Vantagens e desvantagens

| Vantagens | Desvantagens |
|-----------|--------------|
| Escrita de baixa latência | Risco de perda de dados se o cache falhar antes da sincronização |
| Reduz carga no banco (escritas em lote) | Implementação mais complexa |
| Ideal parawrite-heavy workloads | Necessita de processo assíncrono dedicado |

---

## Comparativo entre estratégias

| Critério | Cache-Aside | Write-Through | Write-Behind |
|----------|-------------|---------------|--------------|
| Latência de leitura | Baixa (após warm) | Baixa | Baixa |
| Latência de escrita | Normal | Alta (síncrona) | Muito baixa (assíncrona) |
| Consistência | Eventual | Forte (imediata) | Fraca (assíncrona) |
| Complexidade | Baixa | Média | Alta |
| Risco de perda de dados | Baixo | Baixo | Médio/Alto |
| Uso típico | Catálogo, perfis | Dados críticos | Métricas, analytics |

## CDN Cache

Uma CDN (Content Delivery Network) é uma camada de cache geograficamente distribuída para conteúdo estático (imagens, CSS, JS, vídeos). O funcionamento segue o mesmo princípio de Cache-Aside: o primeiro acesso busca na origem, os seguintes servem do cache do edge node mais próximo do usuário.

Na Pokedex Platform, as imagens dos Pokémon (sprites) são candidatas naturais a CDN cache, já que mudam raramente e são acessadas com alta frequência.

---

> Referência: [fidelissauro.dev/caching](https://fidelissauro.dev/caching/) — System Design: Cache, por Matheus Fidelis
