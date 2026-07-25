# Padrão BFF — Backend for Frontend

## Definição

BFF (Backend for Frontend) é um padrão arquitetural que cria backends especializados para cada tipo de frontend. Em vez de um único backend monolítico atendendo web, mobile, IoT e APIs públicas, cada canal recebe um backend dedicado com contratos, regras de cache, segurança e performance otimizados para suas necessidades específicas.

**BFF não é um componente de infraestrutura.** Diferente de API Gateways e proxies reversos — que são componentes de infraestrutura para unificar e rotear — os BFFs são **aplicações completas**, com deployment, escalabilidade e segurança independentes. Eles podem inclusive atuar como backends de API Gateways.

```
┌──────────┐     ┌──────────────┐     ┌─────────────────┐
│  Mobile  │────▶│  mobile-bff  │────▶│                 │
└──────────┘     └──────────────┘     │                 │
                                      │  microsserviços  │
┌──────────┐     ┌──────────────┐     │  (catálogo,      │
│   Web    │────▶│   web-bff    │────▶│   auth, etc.)    │
└──────────┘     └──────────────┘     │                 │
                                      │                 │
┌──────────┐     ┌──────────────┐     │                 │
│   IoT    │────▶│   iot-bff    │────▶│                 │
└──────────┘     └──────────────┘     └─────────────────┘
```

---

## Responsabilidades Arquiteturais

O BFF atua como intermediário entre o frontend e os microsserviços, assumindo responsabilidades que antes ficavam no cliente:

| Responsabilidade | Descrição |
|-----------------|-----------|
| **Autenticação e autorização** | Propagação de tokens, validação de sessão, escopos por canal |
| **API Composition** | Agregar múltiplas chamadas a serviços distintos em uma única resposta |
| **Adaptação de contratos** | Renomear campos, formatar dados, filtrar campos sensíveis |
| **Cache** | Estratégias de cache específicas por canal (TTLs diferentes para mobile vs web) |
| **Ordenação e filtros** | Lógica de apresentação que não pertence aos serviços core |
| **Fallbacks** | Respostas parciais ou degradadas quando serviços downstream falham |
| **Gestão de estado de tela** | Estados como `empty`, `loading`, `unauthenticated`, `error` |

### O que NÃO é responsabilidade do BFF

- Regras de negócio canônicas (pertencem aos serviços core)
- Persistência de dados de domínio (o BFF pode ter cache, não fonte de verdade)
- Lógica central de autenticação (delega ao `auth-service`)

---

## API Composition Pattern

O padrão de composição de APIs no BFF consiste em consolidar múltiplas chamadas a serviços backend em **um único ponto de entrada** para o frontend.

```
┌──────────┐      ┌──────────────────────────────────────┐
│  Cliente │      │              mobile-bff               │
│          │      │                                       │
│  GET     │      │  ① GET /pokemons ──▶ catálogo-service │
│  /home ──┼─────▶│  ② GET /favorites ──▶ PostgreSQL      │
│          │      │  ③ GET /regions ────▶ catálogo-service │
│          │      │                                       │
│          │      │  Compõe resposta unificada            │
│          │◀─────┼── payload final                       │
└──────────┘      └──────────────────────────────────────┘
```

### Benefícios

- **Redução de latência:** uma chamada do cliente vs. N chamadas individuais
- **Código de frontend mais simples:** lógica de orquestração fica no servidor
- **Flexibilidade:** transformações, enriquecimento e filtragem sem expor complexidade

### Exemplo na Pokedex Platform

A home (`GET /v1/home`) compõe dados de múltiplas fontes em uma única resposta:

```go
// PokemonService — camada de aplicação no mobile-bff
func (s *PokemonService) GetHome(ctx context.Context, params HomeParams) (*domain.HomeScreen, error) {
    // ① Busca pokemons do catálogo
    page, err := s.pokemonRepo.GetAll(ctx, params.Page, params.Size)
    
    // ② Se usuário autenticado, busca favoritos
    var favorites map[string]bool
    if userID := GetUserID(ctx); userID != "" {
        favIDs, _ := s.favoriteRepo.GetUserFavorites(ctx, userID)
        favorites = toSet(favIDs)
    }
    
    // ③ Compõe resposta UI-oriented
    return s.buildHomeResponse(page, favorites, params.Filters), nil
}
```

---

## Segregação de Canais

Cada canal (web, mobile, IoT) tem requisitos distintos. BFFs segregados evitam que um único backend acumule condicionais para cada tipo de cliente.

| Canal | Características típicas do BFF |
|-------|-------------------------------|
| **Web** | SSR, sessão em cache centralizado, JWT, mais funcionalidades |
| **Mobile** | Payloads enxutos, compressão intensa, suporte offline, sincronização |
| **IoT / Eletrodomésticos** | Conexões intermitentes, MQTT/WebSockets, segurança rigorosa, dados mínimos |

### O anti-padrão

```go
// ❌ Anti-padrão: condicionais por canal dentro do mesmo BFF
if channel == "mobile" {
    return compactResponse(data)
} else if channel == "web" {
    return fullResponse(data)
} else if channel == "iot" {
    return minimalResponse(data)
}
```

### O padrão correto

```
mobile/                        web/                        iot/
  mobile-bff (deploy dedicado)   web-bff (deploy dedicado)   iot-bff (deploy dedicado)
  ├─ payloads enxutos            ├─ SSR + cache              ├─ MQTT + WebSocket
  ├─ compressão gzip/brotli     ├─ JWT + sessão             ├─ segurança por dispositivo
  └─ offline sync               └─ feature flags            └─ telemetria mínima
```

Cada BFF conhece apenas as necessidades do seu canal. Sem condicionais, sem acoplamento entre canais.

---

## Segregação de Microfrontends

Em projetos com microfrontends, cada módulo de UI tem seu próprio time e seu próprio BFF:

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  dashboard-mf    │     │  notifications-mf  │     │  profile-mf      │
│  (time A)        │     │  (time B)          │     │  (time C)        │
└────────┬─────────┘     └────────┬─────────┘     └────────┬─────────┘
         │                        │                        │
┌────────┴─────────┐     ┌────────┴─────────┐     ┌────────┴─────────┐
│  dashboard-bff   │     │ notifications-bff │     │   profile-bff    │
│  (time A)        │     │  (time B)         │     │   (time C)       │
└────────┬─────────┘     └────────┬─────────┘     └────────┬─────────┘
         │                        │                        │
         └────────────────────────┼────────────────────────┘
                                  │
                    ┌─────────────┴──────────────┐
                    │     microsserviços core      │
                    └────────────────────────────┘
```

**Vantagens:**

- **Lei de Conway aplicada:** time dono do microfrontend é dono do BFF — alinhamento total UI ↔ API
- **Deployment independente:** rollback de um módulo não afeta os demais
- **Isolamento de falhas:** bug no BFF de dashboard não derruba notificações
- **Unidade de deployment coesa:** frontend + BFF podem ser deployados juntos

---

## Versionamento de Interfaces

BFFs desacoplados facilitam versionamento e experimentação controlada:

```
┌──────────┐     ┌─────────────────┐
│ Cliente  │     │     API Gateway │
│ v1       │────▶│                 │──▶ mobile-bff v1 (estável)
└──────────┘     │  roteamento     │
                 │  por versão     │
┌──────────┐     │  ou feature     │──▶ mobile-bff v2 (canário — 5% do tráfego)
│ Cliente  │     │  flag           │
│ v2       │────▶│                 │──▶ mobile-bff v2 (estável)
└──────────┘     └─────────────────┘
```

**Estratégias:**

| Estratégia | Como funciona |
|------------|---------------|
| **Feature toggles** | Nova versão do BFF ativada por flag para % do tráfego |
| **Canary deployment** | BFF v2 recebe 5% → 25% → 100% progressivamente |
| **Blue-green** | BFF v1 (azul) e v2 (verde) coexistem; switch instantâneo |
| **BFF por versão de API** | `/v1/*` → bff-legacy, `/v2/*` → bff-current |

---

## Resiliência e Blast Radius

Cada BFF é responsável por sua própria resiliência contra falhas de serviços downstream.

### Circuit Breaker

```go
// Padrão já implementado na Pokedex Platform (circuit_breaker.go)
type CircuitBreakerClient struct {
    client           *http.Client
    failureThreshold int
    resetTimeout     time.Duration
    retryBackoff     []time.Duration
    // ...
}

// Estados: Closed → Open → HalfOpen → Closed
```

### Bulkhead

Isolar grupos de chamadas para que uma falha não contamine todo o BFF:

```
┌────────────────────────────────────────────┐
│                mobile-bff                   │
│                                             │
│  ┌─────────────┐  ┌─────────────┐          │
│  │ catálogo    │  │ auth        │          │
│  │ (pool: 10)  │  │ (pool: 5)   │          │
│  │ timeout: 5s │  │ timeout: 3s │          │
│  └─────────────┘  └─────────────┘          │
│                                             │
│  ┌─────────────┐  ┌─────────────┐          │
│  │ favoritos   │  │ cache       │          │
│  │ (pool: 5)   │  │ (pool: 10)  │          │
│  │ timeout: 2s │  │ timeout: 1s │          │
│  └─────────────┘  └─────────────┘          │
└────────────────────────────────────────────┘
```

Cada dependência tem seu próprio pool de conexões e timeout — se o catálogo ficar lento, autenticação e favoritos continuam funcionando.

### Blast Radius

Blast radius é o alcance de uma falha. Para reduzi-lo:

- **Unidades de deployment independentes:** cada BFF é deployado separadamente
- **Bulkheads lógicos:** grupos isolados de chamadas dentro do mesmo BFF
- **Canary/Blue-green:** validar novas versões em pequenos segmentos antes do rollout total
- **Fallbacks:** respostas degradadas mantêm UX funcional mesmo com dependências indisponíveis

---

## Desacoplamento de Métricas

BFFs segregados permitem analisar a experiência de cada canal de forma independente:

```
┌─────────────────────────────────────────────────────┐
│                   Observabilidade                    │
│                                                      │
│  mobile-bff ──▶ latency p95: 120ms, error rate: 0.1% │
│  web-bff    ──▶ latency p95: 80ms,  error rate: 0.3% │
│  iot-bff    ──▶ latency p95: 500ms, error rate: 1.2% │
│                                                      │
│  Serviços downstream (compartilhados):               │
│  ──▶ catálogo-service: 45ms p95 (igual para todos)   │
│  ──▶ auth-service: 30ms p95                          │
└──────────────────────────────────────────────────────┘
```

**O que monitorar por BFF:**

| Métrica | Por que importa |
|---------|----------------|
| **Latência (p50, p95, p99)** | Mobile pode ter latência maior que web (rede celular) |
| **Error rate** | Identificar se um canal específico está com problemas |
| **Throughput (rps)** | Dimensionamento independente por canal |
| **Cache hit rate** | Cada canal pode ter política de cache diferente |
| **Circuit breaker state** | Qual dependência está degradada para qual canal |

### SLOs por canal

Com métricas segregadas, é possível definir SLOs diferentes:

| Canal | Latência SLO | Disponibilidade SLO |
|-------|-------------|---------------------|
| Web | p95 < 100ms | 99.9% |
| Mobile | p95 < 300ms | 99.5% |
| IoT | p95 < 1000ms | 99.9% |

---

> Referência: [fidelissauro.dev/bffs](https://fidelissauro.dev/bffs/) — System Design: Backend for Frontend, por Matheus Fidelis (Jun/2025)
