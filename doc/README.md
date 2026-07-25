# Documentação

Este diretório centraliza a documentação arquitetural da Plataforma Pokedex.

## Objetivos

- Explicar o papel de cada parte principal do repositório.
- Documentar as decisões arquiteturais atuais.
- Explicar como o `mobile-bff` aplica arquitetura hexagonal.
- Tornar futuras refatorações mais fáceis de discutir e revisar.

## Documentos

### Plataforma

- [SYSTEM-OVERVIEW.md](SYSTEM-OVERVIEW.md): visão geral da plataforma.
- [BFF.md](BFF.md): papel e responsabilidades do Backend for Frontend.
- [GATEWAY.md](GATEWAY.md): papel do Kong como ponto de entrada.
- [INFRA.md](INFRA.md): infraestrutura, fluxo de dados e limites operacionais.
- [DECISIONS.md](DECISIONS.md): decisões arquiteturais atuais e trade-offs.
- [SOLID-AND-PATTERNS.md](SOLID-AND-PATTERNS.md): princípios SOLID conectados a patterns, com exemplos em Go.
- [architecture/hexagonal.md](architecture/hexagonal.md): guia detalhado da arquitetura hexagonal aplicada ao `mobile-bff`.

### Guias temáticos

- [bff/](bff/): padrão arquitetural BFF — API Composition, segregação de canais, microfrontends, versionamento, resiliência, métricas.
- [redis/](redis/): guia completo de Redis com Go — cache (conceitos e estratégias), streaming (Streams, Pub/Sub, pipelines/Lua) e operações (Cluster, Sentinel, AWS ElastiCache).
- [referencias/](referencias/): índice curado de artigos de System Design do Matheus Fidelis, mapeados por relevância para a Pokedex Platform.

## Convenção De Nomes

- Nomes de pastas em minúsculo.
- Nomes de arquivos Markdown em inglês e em maiúsculo.
