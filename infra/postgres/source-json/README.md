# Data Source

Este diretório contém a fonte de verdade dos dados da Pokedex.

## Estrutura

```text
infra/postgres/source-json/
├── 01_region.json
├── 02_type.json
├── 03_egg_group.json
├── 04_generation.json
├── 05_ability.json
├── 06_species.json
├── 07_stats.json
├── 08_evolution_chains.json
├── 09_pokemon.json
└── 10_weaknesses.json
```

## Ordem de Processamento

Os arquivos são numerados para respeitar dependências de chaves estrangeiras.

1. Tabelas base: regions, types, egg_groups, generations
2. Tabelas intermediárias: abilities, species, stats, evolution_chains
3. Tabela principal: pokemons
4. Tabelas de relacionamento: pokemon_types, pokemon_abilities, pokemon_egg_groups, pokemon_weaknesses

## Fluxo Operacional

1. Edite os arquivos em `infra/postgres/source-json`.
2. Rode o gerador JSON -> SQL.
3. Gere o seed em `infra/postgres/seeds/init-data.sql`.
4. Recrie o banco para reaplicar schema e seed quando necessário.

## Regras

- Nao alterar a numeração dos arquivos.
- Manter campos e tipos consistentes.
- Evitar editar SQL gerado manualmente.
