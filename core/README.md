# Core

Este diretório concentra o código executável e os ativos de infraestrutura da Plataforma Pokedex.

## Estrutura

```text
core/
├── app/           # Serviços internos de backend
│   ├── auth-service/
│   └── pokemon-catalog-service/
├── bff/           # Aplicações Backend for Frontend
│   └── mobile-bff/
├── gateway/       # Configuração do API Gateway
├── infra/         # Ativos compartilhados de infraestrutura
├── observability/ # Ativos de observabilidade e operação
└── docker-compose.yml
```

## Objetivo

A pasta `core/` mantém os artefatos de implementação reunidos para que a raiz do repositório permaneça focada em materiais transversais, como documentação e arquivos de colaboração.

## Áreas Principais

- `app/`: capacidades de negócio encapsuladas em serviços: `auth-service` (autenticação e ciclo de vida de token) e `pokemon-catalog-service` (catálogo canônico de Pokémons).
- `bff/`: camada de composição voltada ao cliente, atualmente o `mobile-bff`.
- `gateway/`: configuração declarativa do Kong.
- `infra/`: schema do PostgreSQL, seeds, arquivos-fonte em JSON, configuração do Redis e ferramentas de dados.

## Execução Local

A partir da raiz do repositório:

```bash
docker compose -p pokedex -f core/docker-compose.yml up --build
```

## Documentação Relacionada

- [../README.md](../README.md)
- [../doc/SYSTEM-OVERVIEW.md](../doc/SYSTEM-OVERVIEW.md)
- [../doc/BFF.md](../doc/BFF.md)
- [../doc/INFRA.md](../doc/INFRA.md)
