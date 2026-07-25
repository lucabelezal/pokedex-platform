# Conceitos de Cache

## Definição

Cache é uma técnica de otimização que cria uma camada intermediária de dados entre dois componentes. Armazena temporariamente dados custosos ou demorados de serem recuperados da origem, funcionando também como camada de resiliência.

Os dados em cache são tipicamente:
- Resultado de uma operação anterior
- Cópias de dados armazenados em outro lugar

**Cache não é sinônimo de Redis, Memcached ou CDN.** Essas tecnologias implementam capacidades de cache, mas não definem o conceito. Um `map` em memória com `sync.RWMutex` também é cache.

## Consistência de Dados

A sincronização entre cache e armazenamento principal é crítica. Estratégias devem garantir que o cache reflita as mudanças mais recentes nos dados principais.

Exemplo: ao alterar o endereço de um usuário, a operação de escrita deve:
- Deletar a chave de cache correspondente, ou
- Atualizar o cache com o estado mais recente

## Time to Live (TTL)

TTL define um período de vida para um item no cache. Após esse período, o item é removido ou marcado como inválido.

**Por que usar TTL:**
- Previne dados desatualizados
- Recicla informações periodicamente
- Evita consumo desnecessário de memória

No Redis, o TTL é configurado com `EXPIRE` ou diretamente no `SET`:

```go
// TTL de 5 minutos
rdb.Set(ctx, "pokemon:25", data, 5*time.Minute)

// Ou após o SET
rdb.Expire(ctx, "pokemon:25", 5*time.Minute)
```

## Políticas de Evicção

Quando o cache atinge a capacidade máxima, é preciso decidir quais itens remover.

### LRU — Least Recently Used

Remove o item que não foi usado há mais tempo. Baseia-se na premissa de que itens não acessados recentemente têm menor probabilidade de serem acessados no futuro.

No Redis, configure `maxmemory-policy allkeys-lru`:

```conf
maxmemory 256mb
maxmemory-policy allkeys-lru
```

### LFU — Least Frequently Used

Remove itens menos frequentemente acessados. Mais eficiente que LRU em alguns cenários, mas requer rastreamento de frequência.

Redis suporta via `maxmemory-policy allkeys-lfu`.

### FIFO — First In, First Out

Remove itens na ordem em que foram adicionados. Simples, mas não considera frequência de uso.

### RR — Random Replacement

Remove um item aleatório. Fácil de implementar, mas não considera padrões de acesso.

### Tabela comparativa

| Política | Critério | Complexidade | Redis |
|----------|----------|-------------|-------|
| LRU | Tempo desde último acesso | Média | `allkeys-lru` |
| LFU | Frequência de acesso | Alta | `allkeys-lfu` |
| FIFO | Ordem de inserção | Baixa | Não nativa |
| RR | Aleatório | Baixa | Não nativa |

## Invalidação de Cache

Processo de remover ou marcar dados como inválidos. Pode ocorrer de três formas:

1. **Programática:** lógica da aplicação exclui itens específicos
2. **Manual:** comandos explícitos para invalidar itens individuais ou em grupo
3. **Automática:** via TTL — item expira após período definido

No Redis:

```go
// Invalidação manual de uma chave
rdb.Del(ctx, "pokemon:25")

// Invalidação por padrão (todas as chaves com prefixo)
keys, _ := rdb.Keys(ctx, "pokemon:*").Result()
for _, key := range keys {
    rdb.Del(ctx, key)
}

// Invalidação automática via TTL
rdb.Set(ctx, "pokemon:25", data, 10*time.Minute)
```

## Cache Hit, Cache Miss e Hit Rate

### Cache Hit

Ocorre quando o dado solicitado é encontrado no cache. Evita acesso à fonte original (banco de dados, API externa).

### Cache Miss

Ocorre quando o dado não está no cache. Força acesso à fonte original, resultando em maior latência.

### Hit Rate (Taxa de Acertos)

```
Hit Rate = Cache Hits / (Cache Hits + Cache Misses) × 100
```

Exemplo: 800 hits e 200 misses → Hit Rate = 80%

**Interpretação:**
- Hit rate alta (>90%): cache bem dimensionado
- Hit rate baixa (<50%): oportunidade de otimização ou cache desnecessário
- Picos de miss após limpeza de cache: normal, os itens são reconstruídos gradualmente

### Monitoramento no go-redis

O go-redis expõe estatísticas via `PoolStats()`:

```go
stats := rdb.PoolStats()
fmt.Printf("Hits: %d, Misses: %d\n", stats.Hits, stats.Misses)
```

Para métricas mais granulares, instrumente com OpenTelemetry ou exponha via `/metrics`.
