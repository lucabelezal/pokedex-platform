# Ferramenta json2sql

Ferramenta para gerar seed SQL a partir dos JSONs em `core/infra/postgres/source-json`.

## Estrutura

```text
core/infra/postgres/json2sql/
├── cmd/json2sql/main.go
├── go.mod
└── internal/
```

## Saída Esperada

- Arquivo gerado: `core/infra/postgres/seeds/init-data.sql`

## Estado Atual

- O CLI Go já existe como ponto de entrada.
- Não há scripts legados neste projeto.

## Estratégia

1. Editar JSONs em `core/infra/postgres/source-json`.
2. Validar consistência referencial.
3. Gerar SQL determinístico.
4. Recriar banco e aplicar schema + seed.
