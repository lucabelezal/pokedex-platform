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
- **Trade-off**: Escritas (add/remove favorites) ainda usam PostgresFavoriteRepository até que um endpoint de escrita seja criado no catalog-service
- **Scope**: mobile-bff, pokemon-catalog-service
- **Date**: 2026-07-25
- **Status**: active

## Handoff

- **Feature**: production-readiness
- **Phase / Task**: CONCLUÍDO — 30/30 tasks implementados
- **Completed**: Todos os 30 tasks dos 3 batches
- **Last commit**: 39ee1ee (rate limiter Redis)
- **Blockers**: none
- **Uncommitted files**: none
- **Branch**: main
