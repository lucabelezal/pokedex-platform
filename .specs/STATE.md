# STATE

## Decisions

### AD-001
- **Decision**: Usar `sony/gobreaker/v2` como biblioteca de circuit breaker
- **Reason**: Leve (~500 LOC), idiomática Go, zero dependências externas, callback hooks para logging/métricas
- **Trade-off**: Menos recursos que Hystrix/Resilience4j, mas suficiente para o escopo da plataforma
- **Scope**: Todos os outbound HTTP clients no mobile-bff
- **Date**: 2026-07-25
- **Status**: active

### AD-002
- **Decision**: Usar OpenTelemetry SDK com OTLP HTTP exporter para tracing distribuído
- **Reason**: Padrão aberto da indústria, vendor-neutral, suporta W3C TraceContext, pode trocar backend (Jaeger/Zipkin/Tempo) sem mudar código
- **Trade-off**: Mais complexo que logging-only, mas essencial para diagnóstico em produção
- **Scope**: Todos os serviços (mobile-bff, auth-service, pokemon-catalog-service)
- **Date**: 2026-07-25
- **Status**: active

### AD-003
- **Decision**: BFF não acessa mais banco de dados diretamente para favoritos — usa endpoint batch do catalog-service
- **Reason**: Respeita separação de responsabilidades da arquitetura hexagonal; catalog-service é dono dos dados de Pokémon
- **Trade-off**: Escritas (add/remove favorites) migradas para o catalog-service via REST; leitura de IDs de favoritos ainda via PostgresFavoriteRepository
- **Scope**: mobile-bff, pokemon-catalog-service
- **Date**: 2026-07-25
- **Status**: active

### AD-004
- **Decision**: Extrair `DBPool` interface (Query/QueryRow/Exec/Ping/Begin) para permitir mocks pgx em testes de repository
- **Reason**: `pgxmock` não satisfaz `*pgxpool.Pool`; a interface viabiliza testes unitários com cobertura >75% em catalog e auth-service
- **Trade-off**: Pequena abstração extra; `NewXxxRepository(db *pgxpool.Pool)` mantido para produção, `NewXxxRepositoryWithPool(db DBPool)` para testes
- **Scope**: pokemon-catalog-service, auth-service
- **Date**: 2026-08-29
- **Status**: active

### AD-005
- **Decision**: `FavoriteCatalogProvider` movido para `ports/outbound/` (era interface definida em `internal/service/`)
- **Reason**: Pertence ao domínio (outbound port), não à camada de aplicação; respeita regra arquitetural de dependência
- **Trade-off**: Nenhum — movimento mecânico, testes existentes já usavam stub compatível
- **Scope**: mobile-bff
- **Date**: 2026-08-29
- **Status**: active

## Handoff

- **Feature**: production-readiness + test-coverage-and-favorites
- **Phase / Task**: CONCLUÍDO — 30/30 + 15/15 tasks (test-coverage-and-favorites-migration)
- **Completed**: Todos os tasks, todas as lessons endereçadas
- **Cobertura**: catalog-service 75.6%, auth-service 79.3% (gaps PR-25/PR-26/PR-11 resolvidos)
- **Last commit**: 40e7957 (favoritos outbound + testes integração)
- **Blockers**: none
- **Uncommitted files**: .specs/STATE.md (este)
- **Branch**: main
