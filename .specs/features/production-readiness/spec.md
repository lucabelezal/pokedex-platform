# Production Readiness - Specification

## Problem Statement

A plataforma Pokedex possui uma arquitetura sólida conceitualmente (hexagonal, BFF, separação de domínio), mas está incompleta em aspectos de produção: zero resiliência (sem circuit breaker, retry ou graceful degradation), zero observabilidade (sem métricas ou tracing), violações arquiteturais (BFF acessa banco diretamente, código de produção importa package de testes), bugs de lógica (ordenação quebrada, race condition, N+1 queries, erros mascarados) e cobertura de testes inconsistente (catalog-service com 0% de cobertura). Além disso, práticas de segurança nos Dockerfiles estão ausentes (containers rodam como root, credenciais hardcoded, sem security headers).

## Goals

- [ ] Todos os outbound HTTP clients com circuit breaker e retries com backoff
- [ ] Dockerfiles seguros (non-root, pinned images, HEALTHCHECK)
- [ ] Violações arquiteturais corrigidas (BFF sem acesso direto ao DB, mocks isolados do main.go)
- [ ] Mapa de type colors centralizado em um único ponto (DRY)
- [ ] Security headers ativos em todas as respostas do BFF
- [ ] OpenTelemetry tracing nos 3 serviços com propagação W3C TraceContext
- [ ] Métricas Prometheus (request count, latency, errors) nos 3 serviços
- [ ] Logging unificado com `log/slog` em todos os serviços
- [ ] Bugs conhecidos corrigidos (ordenação, race condition, N+1, erros mascarados, loadHome)
- [ ] Cobertura de testes >75% no catalog-service e auth-service
- [ ] Rate limiter distribuído com Redis no BFF

## Out of Scope

| Feature | Reason |
|---------|--------|
| Migração para gRPC | Comunicação HTTP é suficiente para o escopo atual da plataforma |
| Event sourcing / CQRS | Complexidade desnecessária para um catálogo de Pokémon |
| Frontend | O frontend está em outro repositório |
| CI/CD pipeline changes | Fora do escopo desta feature; CI atual cobre build/test/lint |
| Migração para Kubernetes | Docker Compose atual atende o ambiente de desenvolvimento |
| Substituir Kong por Envoy | Kong atende os requisitos atuais de API Gateway |

---

## Assumptions & Open Questions

| Assumption / decision | Chosen default | Rationale | Confirmed? |
|-----------------------|---------------|-----------|------------|
| Biblioteca de circuit breaker | `sony/gobreaker` | Leve, idiomática Go, sem dependências externas, usada nos projetos de referência msfidelis | n |
| Backend de tracing | Jaeger (via OTel) | Mais leve que Zipkin para dev, suporta OTLP nativo | n |
| Endpoint de favorites no catalog-service | `GET /v1/pokemons/favorites?ids=id1,id2,...` | Simples, RESTful, aceita batch por query param | n |
| Location do type colors centralizado | `core/bff/mobile-bff/internal/domain/type_colors.go` | Domínio puro, importável por todas as camadas sem criar dependência circular | n |
| Rate limiter Redis | Sliding window com sorted sets (ZSET) | Padrão já implementado no `go-ratelimiter-sliding-window` do msfidelis | n |
| Cobertura de testes mínima | 75% (convenção do projeto) | Definido no AGENTS.md e copilot-instructions.md | y |

**Open questions:** todas as ambiguidades estão resolvidas ou registradas acima.

---

## User Stories

### P1: Resiliência nos Outbound HTTP Clients ⭐ MVP

**User Story**: Como operador da plataforma, quero que falhas transitórias nos serviços internos (auth-service, pokemon-catalog-service) não causem erros em cascata no BFF, para que a plataforma seja tolerante a falhas e se recupere automaticamente.

**Why P1**: Sem resiliência, qualquer falha temporária nos serviços internos causa erro 500 para o cliente. É o alicerce de produção mais crítico atualmente ausente.

**Acceptance Criteria**:

1. WHEN o BFF faz uma chamada HTTP ao pokemon-catalog-service E o serviço falha (5xx, timeout, conexão recusada) THEN o circuit breaker SHALL abrir após N falhas consecutivas E retornar erro rápido sem tentar nova chamada
2. WHEN o circuit breaker está aberto THEN após um período de half-open o BFF SHALL permitir uma tentativa de teste E fechar o circuito se a chamada for bem-sucedida
3. WHEN o BFF faz uma chamada HTTP que falha com erro transitório (5xx, timeout) THEN o client SHALL retentar até 3 vezes com backoff exponencial (1s, 3s, 10s)
4. WHEN o circuit breaker está aberto E o BFF não pode chamar o serviço THEN o handler SHALL retornar resposta de degradação controlada (status 503 com mensagem clara) em vez de erro genérico
5. WHEN o pokemon-catalog-service está indisponível THEN o BFF SHALL servir respostas com cache ou fallback explícito (indicando degradação) em vez de 2 Pokémon hardcoded silenciosamente

**Independent Test**: Simular falha no pokemon-catalog-service (parar o container) e verificar que o BFF retorna 503 com mensagem de degradação, e que após o serviço voltar o circuit breaker fecha automaticamente.

---

### P1: Segurança — Dockerfiles e Credenciais ⭐ MVP

**User Story**: Como operador de infraestrutura, quero que os containers da plataforma executem com práticas seguras (non-root user, imagens pinadas, sem credenciais expostas) para que o ambiente esteja em conformidade com boas práticas de segurança.

**Why P1**: Containers rodando como root e credenciais hardcoded são vulnerabilidades críticas que não podem permanecer em produção.

**Acceptance Criteria**:

1. WHEN qualquer Dockerfile é construído THEN o estágio de runtime SHALL conter diretiva `USER 1000` (ou UID não-root equivalente)
2. WHEN qualquer Dockerfile referencia imagem base THEN a imagem SHALL ser pinada por SHA digest (não `alpine:latest`)
3. WHEN o código de produção inicializa conexão com banco de dados THEN NÃO SHALL conter credenciais hardcoded como fallback (a string `user:password@localhost:5432/pokedex` deve ser removida)
4. WHEN o mobile-bff responde a requisições HTTP THEN os headers de segurança SHALL estar presentes: `Strict-Transport-Security`, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`
5. WHEN o auth-service e pokemon-catalog-service são construídos com Dockerfile THEN SHALL conter diretiva `HEALTHCHECK`

**Independent Test**: Inspecionar Dockerfiles para `USER 1000`, SHA digests, HEALTHCHECK. Fazer scan de código por strings de credenciais. Fazer curl no BFF e verificar headers de segurança.

---

### P1: Arquitetura — Isolar DB do BFF e Centralizar Type Colors ⭐ MVP

**User Story**: Como desenvolvedor da plataforma, quero que o BFF não acesse o banco de dados diretamente e que o mapa de cores dos tipos Pokémon esteja centralizado, para que a arquitetura hexagonal seja respeitada e a manutenção seja simplificada.

**Why P1**: O BFF acessando o banco diretamente quebra a separação de responsabilidades (o catalog-service deveria ser o dono dos dados). O mapa de type colors duplicado em 5 locais com valores inconsistentes causa bugs sutis.

**Acceptance Criteria**:

1. WHEN o BFF precisa consultar ou modificar favoritos THEN a operação SHALL ser feita via endpoint REST no pokemon-catalog-service, não via acesso direto ao PostgreSQL
2. WHEN o pokemon-catalog-service recebe uma requisição de favoritos THEN SHALL existir endpoint `GET /v1/pokemons/favorites?ids=id1,id2,...` que retorna detalhes dos Pokémon em batch
3. WHEN qualquer camada do sistema precisa mapear tipo de Pokémon para cor THEN SHALL consultar um único mapa centralizado em `internal/domain/type_colors.go`
4. WHEN o `cmd/server/main.go` inicializa o BFF THEN NÃO SHALL importar o package `tests/mocks` (mocks devem estar em `adapters/outbound/memory/`)
5. WHEN o pokemon-catalog-service não consegue conectar ao PostgreSQL THEN SHALL retornar erro explícito em vez de servir silenciosamente apenas 2 Pokémon hardcoded

**Independent Test**: Remover credenciais de banco do BFF e verificar que favoritos continuam funcionando via catalog-service. Verificar que `grep -r "typeColor\|type_color\|getTypeColor"` retorna resultado apenas no arquivo centralizado. Verificar que `grep "tests/mocks" cmd/server/main.go` não retorna resultado.

---

### P2: Observabilidade — Tracing, Métricas e Logging Unificado

**User Story**: Como operador da plataforma, quero visibilidade completa sobre o comportamento dos serviços (tracing distribuído, métricas de latency/erros, logs estruturados consistentes) para que eu possa diagnosticar problemas rapidamente em produção.

**Why P2**: Observabilidade é essencial para operação, mas a plataforma funciona sem ela (embora com diagnóstico prejudicado).

**Acceptance Criteria**:

1. WHEN uma requisição HTTP atravessa o Kong → BFF → catalog-service ou auth-service THEN um trace context (W3C TraceContext) SHALL ser propagado entre todos os serviços
2. WHEN o BFF, auth-service ou catalog-service processam requisições THEN métricas Prometheus SHALL ser expostas no endpoint `/metrics` incluindo: request count (por método, path, status), request latency (histograma), e erros (por tipo)
3. WHEN o auth-service registra eventos (login, logout, erros) THEN SHALL usar `log/slog` estruturado em vez de `log.Printf`
4. WHEN qualquer serviço é iniciado THEN SHALL expor endpoint `/health` para liveness E endpoint `/ready` para readiness (verificando dependências downstream)
5. WHEN o tracing está configurado THEN spans SHALL incluir atributos de: método HTTP, path, status code, e duração de chamadas externas

**Independent Test**: Subir Jaeger, fazer requisição ao BFF que chama catalog-service, verificar que uma trace completa aparece no Jaeger com spans de ambos os serviços. Fazer curl em `/metrics` e verificar que métricas são expostas.

---

### P2: Correção de Bugs — Ordenação, Race Condition, N+1, Erros Mascarados

**User Story**: Como usuário da API, quero que a listagem de Pokémon seja ordenada numericamente, que favoritos não tenham condição de corrida, que a home page não faça dezenas de chamadas HTTP desnecessárias, e que erros reais do banco de dados não sejam mascarados como "pokemon não encontrado".

**Why P2**: Bugs que afetam a experiência do usuário e a confiabilidade do sistema. Não bloqueiam funcionalidade, mas degradam qualidade.

**Acceptance Criteria**:

1. WHEN a home page ordena Pokémon por menor número THEN a ordenação SHALL ser numérica (`"25" < "100"` resulta em Charizard depois de Pikachu) e NÃO lexicográfica
2. WHEN um usuário adiciona um favorito que já existe THEN a operação SHALL ser idempotente via `INSERT ... ON CONFLICT DO NOTHING` sem race condition de check-then-act
3. WHEN o endpoint de favoritos retorna detalhes de múltiplos Pokémon THEN uma única chamada batch SHALL ser feita ao catalog-service em vez de N chamadas individuais
4. WHEN um erro de banco de dados ocorre no `PostgresPokemonRepository.GetByID` (timeout, conexão recusada) THEN o erro SHALL ser propagado como erro interno e NÃO convertido para `ErrPokemonNotFound`
5. WHEN a home page carrega Pokémon com filtro de tipo ou região THEN o filtro e ordenação SHALL ser delegados ao catalog-service (query params) em vez de carregar todos os Pokémon em memória

**Independent Test**: Testar ordenação com Pokémon nº 25 (Pikachu) e nº 100 (Voltorb) e verificar ordem correta. Testar adição concorrente do mesmo favorito e verificar que apenas uma linha é inserida.

---

### P2: Testes — Cobertura no Catalog-Service e Auth-Service

**User Story**: Como desenvolvedor da plataforma, quero que todos os serviços tenham cobertura de testes adequada (>75%) para que mudanças possam ser feitas com confiança.

**Why P2**: Testes são infraestrutura de qualidade, mas não afetam o comportamento em runtime.

**Acceptance Criteria**:

1. WHEN os testes do pokemon-catalog-service são executados THEN SHALL cobrir >75% do código incluindo handlers HTTP, repository e lógica de negócio
2. WHEN os testes do auth-service são executados THEN SHALL cobrir handlers de login e signup (atualmente não testados)
3. WHEN os testes de integração do mobile-bff são executados sem PostgreSQL disponível THEN SHALL falhar explicitamente em vez de pular silenciosamente (`t.Skipf`)
4. WHEN os testes do pokemon-catalog-service são executados THEN SHALL usar mocks de repository para isolar os handlers HTTP

**Independent Test**: Executar `go test -coverprofile=coverage.out ./...` no catalog-service e verificar cobertura >75%.

---

### P3: Rate Limiter Distribuído com Redis

**User Story**: Como operador da plataforma, quero que o rate limiting funcione de forma distribuída entre múltiplas instâncias do BFF, para que o limite seja consistente independentemente de qual instância recebe a requisição.

**Why P3**: O rate limiter in-memory atual funciona para uma única instância. Múltiplas instâncias exigiriam backend compartilhado, mas a plataforma atualmente roda com instância única.

**Acceptance Criteria**:

1. WHEN múltiplas instâncias do BFF estão ativas THEN o rate limiting SHALL ser consistente usando Redis como backend compartilhado
2. WHEN o Redis está indisponível THEN o rate limiter SHALL operar em modo degradado (in-memory fallback) com log de aviso

**Independent Test**: Subir 2 instâncias do BFF com Redis, fazer requisições em ambas e verificar que o contador de rate limit é compartilhado.

---

### P3: Infraestrutura — Pool Config e Health Probes

**User Story**: Como operador da plataforma, quero que as conexões de banco de dados tenham pool sizes explícitos e que os serviços tenham health probes separados (liveness vs readiness), para melhor controle operacional.

**Why P3**: Otimizações de infraestrutura. O sistema funciona com defaults, mas com configurações explícitas ganha-se previsibilidade.

**Acceptance Criteria**:

1. WHEN qualquer serviço conecta ao PostgreSQL THEN o pool SHALL ser configurado com `MaxConns`, `MinConns`, `MaxConnLifetime` e `HealthCheckPeriod` explícitos
2. WHEN o auth-service e pokemon-catalog-service expõem health checks THEN SHALL ter endpoints separados: `/health` (liveness — sempre ok se processo está vivo) e `/ready` (readiness — verifica conexão com banco)

**Independent Test**: Verificar código de configuração do pgxpool para valores explícitos. Fazer curl em `/ready` com banco parado e verificar status 503.

---

### P3: Qualidade de Código — Separar Arquivos Monolíticos e Remover Código Morto

**User Story**: Como desenvolvedor da plataforma, quero que arquivos com mais de 500 linhas sejam divididos e que código morto seja removido, para facilitar navegação e manutenção.

**Why P3**: Qualidade de código afeta velocidade de desenvolvimento, mas não o comportamento em runtime.

**Acceptance Criteria**:

1. WHEN inspecionamos o código da plataforma THEN nenhum arquivo SHALL ter mais de 500 linhas (dividindo responsabilidades quando necessário)
2. WHEN inspecionamos o código da plataforma THEN métodos `Validate()` não utilizados, erros nunca referenciados (`ErrUserNotFound`, `ErrInvalidPagination`), e stubs vazios (`GetFavorites` que retorna `[]string{}, nil`) SHALL ser removidos
3. WHEN inspecionamos o response_builder.go THEN responsabilidades misturadas (type colors, number formatting, JSON helpers, cookies) SHALL estar em arquivos separados

**Independent Test**: `find . -name "*.go" -exec wc -l {} + | sort -rn | head` não deve mostrar arquivos >500 linhas. `grep -r "ErrUserNotFound\|ErrInvalidPagination" --include="*.go"` deve retornar apenas definição (se for usada) ou zero resultados.

---

## Edge Cases

- WHEN o circuit breaker abre para o auth-service E o usuário tenta login THEN o BFF SHALL retornar 503 com mensagem "serviço de autenticação temporariamente indisponível" e NÃO tentar fallback que exponha dados sem autenticação
- WHEN o catalog-service recebe `GET /v1/pokemons/favorites?ids=` com lista vazia THEN SHALL retornar array vazio `[]` e não erro
- WHEN o catalog-service recebe `GET /v1/pokemons/favorites?ids=` com mais de 100 IDs THEN SHALL retornar erro 400 "máximo de 100 IDs por requisição"
- WHEN o OTel exporter (Jaeger) está indisponível THEN os serviços SHALL continuar operando normalmente com log de aviso (tracing não é crítico para operação)
- WHEN o endpoint `/metrics` é acessado sem autenticação THEN SHALL expor métricas (endpoint público, padrão Prometheus)

---

## Requirement Traceability

| Requirement ID | Story | Phase | Status |
|---------------|-------|-------|--------|
| PR-01 | P1: Resiliência — circuit breaker abre após N falhas | Design | Pending |
| PR-02 | P1: Resiliência — half-open e fechamento automático | Design | Pending |
| PR-03 | P1: Resiliência — retries com backoff exponencial (1s, 3s, 10s) | Design | Pending |
| PR-04 | P1: Resiliência — resposta de degradação controlada (503) | Design | Pending |
| PR-05 | P1: Resiliência — fallback explícito em vez de 2 Pokémon hardcoded | Design | Pending |
| PR-06 | P1: Segurança — USER 1000 nos Dockerfiles | Design | Pending |
| PR-07 | P1: Segurança — imagens pinadas por SHA digest | Design | Pending |
| PR-08 | P1: Segurança — remover credenciais hardcoded do código | Design | Pending |
| PR-09 | P1: Segurança — security headers (HSTS, X-Content-Type-Options, X-Frame-Options) | Design | Pending |
| PR-10 | P1: Segurança — HEALTHCHECK nos Dockerfiles de auth e catalog | Design | Pending |
| PR-11 | P1: Arquitetura — BFF sem acesso direto ao DB para favoritos | Design | Pending |
| PR-12 | P1: Arquitetura — endpoint batch `GET /v1/pokemons/favorites?ids=` | Design | Pending |
| PR-13 | P1: Arquitetura — mapa de type colors centralizado | Design | Pending |
| PR-14 | P1: Arquitetura — main.go não importa tests/mocks | Design | Pending |
| PR-15 | P1: Arquitetura — catalog-service não serve 2 Pokémon hardcoded como fallback | Design | Pending |
| PR-16 | P2: Observabilidade — tracing com W3C TraceContext nos 3 serviços | - | Pending |
| PR-17 | P2: Observabilidade — métricas Prometheus `/metrics` | - | Pending |
| PR-18 | P2: Observabilidade — logging unificado com slog no auth-service | - | Pending |
| PR-19 | P2: Observabilidade — health probes separados `/health` + `/ready` | - | Pending |
| PR-20 | P2: Bugs — ordenação numérica na home page | - | Pending |
| PR-21 | P2: Bugs — race condition TOCTOU no AddFavorite | - | Pending |
| PR-22 | P2: Bugs — batch no endpoint de favoritos (substitui N+1) | - | Pending |
| PR-23 | P2: Bugs — erros de DB propagados corretamente (não mascarados) | - | Pending |
| PR-24 | P2: Bugs — filtro/ordenação delegados ao catalog-service | - | Pending |
| PR-25 | P2: Testes — cobertura >75% no catalog-service | - | Pending |
| PR-26 | P2: Testes — testes de login/signup no auth-service | - | Pending |
| PR-27 | P2: Testes — integração não pula silenciosamente | - | Pending |
| PR-28 | P3: Rate limiter distribuído com Redis | - | Pending |
| PR-29 | P3: Pool config explícito no pgxpool | - | Pending |
| PR-30 | P3: Arquivos monolíticos separados em múltiplos arquivos | - | Pending |
| PR-31 | P3: Código morto removido | - | Pending |

**Coverage:** 31 total, 0 mapped to tasks, 31 unmapped

---

## Success Criteria

- [ ] Nenhum outbound HTTP client sem circuit breaker e retries
- [ ] Todos os Dockerfiles com USER não-root e HEALTHCHECK
- [ ] Zero credenciais hardcoded no código de produção
- [ ] BFF sem acesso direto ao banco para operações de favoritos
- [ ] Type colors centralizado em um único arquivo
- [ ] Tracing distribuído funcional (trace completa visível no Jaeger)
- [ ] Métricas Prometheus expostas em todos os serviços
- [ ] Todos os serviços usando `log/slog` estruturado
- [ ] Ordenação numérica correta na home page
- [ ] Cobertura de testes >75% em todos os serviços
