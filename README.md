# Plataforma Pokédex

Projeto de estudo de arquitetura backend moderna usando Go, BFF, API Gateway, Redis, PostgreSQL e Docker.

## Arquitetura

Cliente -> Kong -> BFF -> Serviço -> PostgreSQL/Redis

## Stack Tecnológico

- Go (net/http)
- Kong Gateway Open Source
- PostgreSQL
- Redis
- Docker Compose

## Estrutura do Repositório

```text
.
├── core/
│   ├── app/
│   │   └── pokemon-catalog-service/
│   ├── bff/
│   │   └── mobile-bff/
│   ├── gateway/
│   │   └── kong/
│   ├── infra/
│   │   ├── postgres/
│   │   │   ├── schema/
│   │   │   ├── seeds/
│   │   │   ├── source-json/
│   │   │   └── json2sql/
│   │   └── redis/
│   └── docker-compose.yml
└── doc/
```

## Documentacao

A documentacao arquitetural do projeto fica em `doc/`.
A implementacao executavel da plataforma fica em `core/`.

- Visao geral: `doc/SYSTEM-OVERVIEW.md`
- BFF: `doc/BFF.md`
- Gateway: `doc/GATEWAY.md`
- Infraestrutura: `doc/INFRA.md`
- Decisoes arquiteturais: `doc/DECISIONS.md`
- SOLID e patterns: `doc/SOLID-AND-PATTERNS.md`
- Visao do runtime: `core/README.md`

## UI/UX Design

### Figma Community

- [Pokémon App By Junior Saraiva](https://www.figma.com/pt-br/comunidade/file/1202971127473077147/pokedex-pokemon-app)

## Como Executar

```bash
docker compose -p pokedex -f core/docker-compose.yml up --build
```

### Endpoints

- Proxy Kong: `http://localhost:8000`
- Admin Kong: `http://localhost:8001`
- Saúde do BFF via Kong: `http://localhost:8000/bff/health`
- Rota Ping via Kong: `http://localhost:8000/v1/pokemon/ping`

Endpoints de autenticacao (via BFF):
- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

Endpoints autenticados (via BFF):
- `GET /api/v1/me`
- `GET /api/v1/me/favorites`
- `POST /api/v1/pokemons/{id}/favorite`
- `DELETE /api/v1/pokemons/{id}/favorite`

## Pipeline de Dados

- Arquivos JSON de origem: `core/infra/postgres/source-json/*.json`
- Schema do banco de dados: `core/infra/postgres/schema/schema.sql`
- Seed gerada: `core/infra/postgres/seeds/init-data.sql`
- CLI do gerador: `core/infra/postgres/json2sql/cmd/json2sql`

### Gerador JSON para SQL

A CLI `json2sql` lê dados de Pokémon de 10 arquivos JSON e gera comandos SQL INSERT completos.

#### Executando o Gerador

A partir do diretório do módulo:
```bash
cd core/infra/postgres/json2sql
go run ./cmd/json2sql/ --input ../source-json --output ../seeds/init-data.sql
```

Ou com validação rigorosa (sai na primeira inconsistência de FK):
```bash
go run ./cmd/json2sql/ --input ../source-json --output ../seeds/init-data.sql --strict
```

#### Flags

- `--input` (padrão: `core/infra/postgres/source-json`) — diretório com 10 arquivos JSON de origem
- `--output` (padrão: `core/infra/postgres/seeds/init-data.sql`) — caminho para escrever o SQL gerado
- `--strict` (padrão: false) — falhar em avisos de integridade referencial

#### Mapeamento de Arquivos JSON

| Arquivo | Tabela | Propósito |
|---------|--------|----------|
| 01_region.json | regions | Regiões da Pokédex |
| 02_type.json | types | Tipos de batalha Pokémon e cores |
| 03_egg_group.json | egg_groups | Categorias de reprodução |
| 04_generation.json | generations | Gerações de jogos |
| 05_ability.json | abilities | Habilidades Pokémon |
| 06_species.json | species | Metadados de espécies |
| 07_stats.json | stats | Stats base de cada Pokémon |
| 08_evolution_chains.json | evolution_chains | Árvores de evolução (armazenadas como JSONB) |
| 09_pokemon.json | pokemons | Dados principais do Pokémon + referências muitos-para-muitos |
| 10_weaknesses.json | pokemon_weaknesses | Efetividade de tipo (busca de nome→id) |

#### Tratamento de Tabelas Especiais

**evolution_chains**: O campo `chain` do JSON é armazenado como JSONB na coluna `chain_data`.

**pokemons**: Também gera:
- `pokemon_types` – uma linha por type_id
- `pokemon_abilities` – uma linha por habilidade com flag is_hidden
- `pokemon_egg_groups` – uma linha por egg_group_id

**pokemon_weaknesses**: Nomes de tipos são resolvidos para IDs usando o mapeamento canônico de `02_type.json` no tempo de geração.

#### Validação de Pré-voo

O gerador valida integridade referencial antes da geração:
- Todas as referências de FK são verificadas contra IDs reais nos dados de origem.
- Avisos são impressos para IDs órfãos (ex: pokemon_id 99 referenciado mas não presente).
- Com `--strict`, a ferramenta sai imediatamente no primeiro aviso.
- Sem `--strict`, avisos são registrados e linhas inválidas são puladas.

## Estratégia de Banco de Dados

Este projeto intencionalmente não usa ferramentas de migração na v1.

- Mudanças estruturais são feitas diretamente em `core/infra/postgres/schema/schema.sql`.
- Quando a estrutura muda, recrie o banco de dados do zero.
- Carregamento de seed é completo e determinístico de `core/infra/postgres/seeds/init-data.sql`.

## Segurança

- Não faça commit de segredos reais.
- Mantenha credenciais locais apenas em `.env`.
- Use `.env.example` como template público.
