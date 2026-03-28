# Infraestrutura

## Objetivo

O diretório `core/infra/` contém os ativos técnicos compartilhados necessários para executar a plataforma localmente e em ambientes containerizados.

## PostgreSQL

O PostgreSQL é usado como principal armazenamento persistente.

### Arquivos

- `core/infra/postgres/schema/schema.sql`: schema do banco de dados.
- `core/infra/postgres/seeds/init-data.sql`: seed gerada.
- `core/infra/postgres/source-json/`: fonte de verdade do conteúdo do catálogo.
- `core/infra/postgres/json2sql/`: ferramenta que converte arquivos JSON de origem em seed SQL.

### Pipeline Atual De Dados

```text
source-json/*.json
  -> json2sql
    -> seeds/init-data.sql
      -> inicialização do container PostgreSQL
```

Esse fluxo é adequado para ambientes locais determinísticos e projetos de estudo porque mantém o dataset versionado e reproduzível.

## Redis

O Redis está provisionado no Docker Compose, mas seu papel arquitetural ainda não aparece com tanta força no código da aplicação.

Isso normalmente indica uma de duas situações:

- uso planejado para o futuro
- infraestrutura pronta, mas ainda não integrada aos fluxos centrais

Isso deveria continuar sendo documentado de forma explícita conforme o projeto evolui.

## Docker Compose

O `core/docker-compose.yml` descreve a topologia local completa:

- PostgreSQL
- Redis
- `pokemon-catalog-service`
- `auth-service`
- `mobile-bff`
- Kong

Esse é um dos documentos mais claros da arquitetura real de execução hoje.

## Oportunidades De Melhoria

- Documentar qual serviço é dono de quais tabelas.
- Esclarecer se o BFF pode persistir seus próprios dados no longo prazo ou apenas temporariamente.
- Decidir se o Redis será usado para cache, sessões, rate limiting, ou se não será usado.
- Adicionar documentação de variáveis de ambiente por serviço.
