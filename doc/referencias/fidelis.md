# Referências — System Design (Matheus Fidelis)

Índice curado dos artigos de System Design do [fidelissauro.dev](https://fidelissauro.dev/), mapeados por relevância para a Pokedex Platform. Os artigos já adaptados para o projeto possuem links para a documentação local.

> Fonte: [fidelissauro.dev/system-design-sumario](https://fidelissauro.dev/system-design-sumario/)

---

## Já documentado no projeto

| # | Artigo | Doc local | O que cobre |
|---|--------|-----------|-------------|
| 6 | [Estratégias de Cache](https://fidelissauro.dev/caching/) | [doc/redis/cache/](../redis/cache/) | Cache-Aside, Write-Through, Write-Behind, TTL, evicção, hit rate, Redis + Go |
| 10 | [Backend for Frontends](https://fidelissauro.dev/bffs/) | [doc/bff/](../bff/) + [doc/BFF.md](../BFF.md) | API Composition, segregação de canais, microfrontends, resiliência, métricas |
| 5 | Databases, Modelos de Dados e Indexação | Parcial — [doc/INFRA.md](../INFRA.md) | PostgreSQL schema, seeds, migrations |

---

## Alta prioridade — aplicação direta no código atual

| # | Artigo | Link | Por que importa |
|---|--------|------|-----------------|
| 22 | [Patterns de Resiliência](https://fidelissauro.dev/resilience/) | fidelissauro.dev/resilience | Já temos circuit breaker (`circuit_breaker.go`). Cobrir idempotência, retries, fallbacks, timeouts, bulkhead |
| 14 | [Mensageria, Eventos, Streaming e Arquitetura Assíncrona](https://fidelissauro.dev/mensageria/) | fidelissauro.dev/mensageria | Redis Streams já documentado em `doc/redis/streaming/`. Casos reais: notificações de favoritos, eventos de auth |
| 13 | [Padrões de Comunicação Síncronos](https://fidelissauro.dev/comunicacao-sincrona/) | fidelissauro.dev/comunicacao-sincrona | Base do que já existe: BFF → catalog-service e auth-service via HTTP. Fundamenta REST, gRPC, contratos |
| 23 | [Monitoramento e Observabilidade](https://fidelissauro.dev/observabilidade/) | fidelissauro.dev/observabilidade | OpenTelemetry já está no código. Métricas, tracing, logging, SLOs |
| 24 | [Bulkhead Pattern](https://fidelissauro.dev/bulkheads/) | fidelissauro.dev/bulkheads | Isolar falhas entre dependências do BFF (catálogo, auth, favoritos, cache) |
| 29 | [Single Point of Failure](https://fidelissauro.dev/single-point-of-failure/) | fidelissauro.dev/single-point-of-failure | Identificar SPOFs na plataforma: PostgreSQL, Redis, Kong, serviços |

---

## Média prioridade — evolução futura

| # | Artigo | Link | Por que importa |
|---|--------|------|-----------------|
| 7 | [Microsserviços, Monolitos e Domínios](https://fidelissauro.dev/microsservicos/) | fidelissauro.dev/microsservicos | Formalizar fronteiras de domínio: catálogo, auth, favoritos. A plataforma já é distribuída |
| 9 | [API Gateways](https://fidelissauro.dev/api-gateways/) | fidelissauro.dev/api-gateways | Kong já está no compose. Enriquecer `doc/GATEWAY.md` com rate limiting, auth, routing |
| 8 | [Load Balancers e Proxies Reversos](https://fidelissauro.dev/load-balancers/) | fidelissauro.dev/load-balancers | Complementa API Gateway. Estratégias de distribuição de carga |
| 15 | [Performance, Capacidade e Escalabilidade](https://fidelissauro.dev/performance/) | fidelissauro.dev/performance | Métricas de capacidade dos serviços, dimensionamento |
| 19 | [CQRS](https://fidelissauro.dev/cqrs/) | fidelissauro.dev/cqrs | Separar leitura (catálogo) de escrita (favoritos). Modelo de queries otimizadas |
| 20 | [Saga Pattern](https://fidelissauro.dev/saga/) | fidelissauro.dev/saga | Coordenar transações distribuídas: favoritar Pokémon envolve auth + catálogo + PostgreSQL |
| 21 | [Event Sourcing](https://fidelissauro.dev/event-sourcing/) | fidelissauro.dev/event-sourcing | Log de eventos para favoritos, histórico de ações do usuário |
| 26 | [Modelos de Deployment](https://fidelissauro.dev/deployment/) | fidelissauro.dev/deployment | Blue/green, canary, feature toggles. Aplicável ao CI/CD do projeto |
| 12 | [Princípios de Concorrência e Paralelismo](https://fidelissauro.dev/concorrencia/) | fidelissauro.dev/concorrencia | Go routines, channels, padrões de concorrência no BFF |
| 16 | [Scale Cube](https://fidelissauro.dev/scale-cube/) | fidelissauro.dev/scale-cube | Estratégias de escala para cada serviço da plataforma |

---

## Baixa prioridade — avançado ou escopo diferente

| # | Artigo | Link | Por que importa |
|---|--------|------|-----------------|
| 17 | [Sharding e Particionamento](https://fidelissauro.dev/sharding/) | fidelissauro.dev/sharding | Relevante se a plataforma crescer para múltiplos shards de PostgreSQL |
| 18 | [Replicação de Dados](https://fidelissauro.dev/replicacao/) | fidelissauro.dev/replicacao | Read replicas para o catálogo se o volume de leitura crescer |
| 25 | [Cell-Based Architecture](https://fidelissauro.dev/cell-based/) | fidelissauro.dev/cell-based | Arquitetura avançada para multitenancy ou isolamento regional |
| 27 | [Capacity Planning e Teoria das Filas](https://fidelissauro.dev/capacity-planning/) | fidelissauro.dev/capacity-planning | Modelagem matemática de capacidade. Útil para testes de carga |
| 28 | [Testes de Carga e Estresse](https://fidelissauro.dev/testes-de-carga/) | fidelissauro.dev/testes-de-carga | Complementa capacity planning com execução prática |
| 11 | [Service Mesh Pattern](https://fidelissauro.dev/service-mesh/) | fidelissauro.dev/service-mesh | Istio, Linkerd — overkill para o estágio atual da plataforma |

---

## Fundamentos — base teórica

| # | Artigo | Link | Por que importa |
|---|--------|------|-----------------|
| 0 | [Teoria das Janelas Quebradas](https://fidelissauro.dev/janelas-quebradas/) | fidelissauro.dev/janelas-quebradas | Cultura de engenharia, qualidade de código, disciplina técnica |
| 1 | [Protocolos de Rede](https://fidelissauro.dev/protocolos-de-rede/) | fidelissauro.dev/protocolos-de-rede | TCP, UDP, HTTP, DNS — base de toda comunicação entre serviços |
| 2 | [Storage, RAID e I/O](https://fidelissauro.dev/storage/) | fidelissauro.dev/storage | Fundamentos de persistência. Menos crítico com PostgreSQL gerenciado |
| 3 | [Teorema CAP, ACID e BASE](https://fidelissauro.dev/cap/) | fidelissauro.dev/cap | Consistência vs. disponibilidade no PostgreSQL e Redis |
| 4 | [Teorema PACELC](https://fidelissauro.dev/pacelc/) | fidelissauro.dev/pacelc | Extensão do CAP para sistemas com particionamento |

---

## Fora do escopo do projeto

| # | Artigo | Link | Categoria |
|---|--------|------|-----------|
| — | [Staff Framework — STAR](https://fidelissauro.dev/star/) | fidelissauro.dev/star | Carreira / Staff+ |
| — | [Navalha de Ockham](https://fidelissauro.dev/navalha-de-ockham/) | fidelissauro.dev/navalha-de-ockham | Carreira / Staff+ |
| — | [SMART Method](https://fidelissauro.dev/staff-smart-methods/) | fidelissauro.dev/staff-smart-methods | Carreira / Staff+ |
| — | [Disaster Recovery na AWS](https://fidelissauro.dev/disaster-recovery/) | fidelissauro.dev/disaster-recovery | Infra AWS específica |
| — | [Rate Limit com Istio](https://fidelissauro.dev/istio-rate-limit/) | fidelissauro.dev/istio-rate-limit | Infra Kubernetes/Istio |

---

## Progresso de documentação

| Status | Artigos |
|--------|---------|
| ✅ Documentado | Cache (6), BFF (10) |
| 🔜 Próximos | Resiliência (22), Comunicação Síncrona (13), Mensageria (14) |
| 📋 Planejado | Observabilidade (23), Bulkhead (24), SPOF (29), Microsserviços (7), API Gateway (9) |
| ⬜ Backlog | Demais artigos de média e baixa prioridade |

---

> **Prefácio do autor:** *"Este livro é uma tentativa de compilar o conhecimento de uma vida. Não como um repositório definitivo, mas como um recorte sincero do que foi possível compreender. Conhecimento precisa ser questionado, debatido e evolutivo."* — Matheus Fidelis
