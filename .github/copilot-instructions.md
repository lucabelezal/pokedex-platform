# Instrucoes do Copilot para este repositorio

## Idioma
- Escreva respostas e explicacoes em portugues do Brasil.
- Gere mensagens de commit em portugues do Brasil.

## Estrutura do repositorio
- A raiz do repositorio concentra materiais transversais, como `doc/`, `.github/` e arquivos de colaboracao.
- A implementacao executavel da plataforma fica em `core/`.
- Dentro de `core/`, use a organizacao:
  - `core/app/` para servicos internos
  - `core/bff/` para BFFs
  - `core/gateway/` para configuracao do gateway
  - `core/infra/` para ativos de infraestrutura e dados
  - `core/docker-compose.yml` como compose principal da plataforma

## Nomes oficiais dos servicos
- Use os nomes oficiais atuais da plataforma:
  - `mobile-bff`
  - `pokemon-catalog-service`
  - `auth-service`
- Nao introduza novamente o nome legado `pokedex-service`.
- Ao falar da plataforma como um todo, `Pokedex Platform` ou `Plataforma Pokedex` continua fazendo sentido.

## Responsabilidades arquiteturais
- `mobile-bff` e o servico voltado para experiencia do cliente.
- `pokemon-catalog-service` e a fonte canonica do catalogo de Pokemon.
- `auth-service` concentra autenticacao e ciclo de vida de token.
- Evite deslocar regras canonicas de catalogo para o BFF.

## Hexagonal no mobile-bff
- Preserve a direcao das dependencias:
  - adapters de entrada dependem de use cases
  - use cases dependem de ports
  - adapters externos implementam ports
- Nao faca handlers HTTP dependerem diretamente de clients concretos de infraestrutura.
- Nao use `tests/` como dependencia de runtime.
- Ao introduzir integracoes externas, normalize erros no adapter externo ou na camada de aplicacao, nao no handler HTTP.

## Convencao de commits (obrigatorio)
- Use sempre Conventional Commits.
- Formato: `tipo(escopo-opcional): descricao curta em portugues`.
- Tipos permitidos: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `build`, `perf`, `revert`.
- Exemplos validos:
  - `feat(favoritos): bloquear remocao sem autenticacao`
  - `fix(bff): corrigir validacao de token JWT`
  - `docs(readme): atualizar fluxo de execucao local`

## Qualidade das mensagens
- Evite mensagens genericas como `update`, `ajustes`, `wip`.
- Descreva o impacto funcional de forma objetiva.
- Quando houver mudanca relevante, inclua corpo explicando contexto e motivacao.

## Documentacao
- Atualize `doc/` quando houver mudanca relevante de arquitetura, nomenclatura ou responsabilidade entre servicos.
- Se a mudanca afetar a estrutura do runtime, revise tambem `README.md` e `core/README.md`.

## Colaboracao com CI
- Se o workflow de validacao de commits falhar, ajuste a mensagem para o padrao Conventional Commits.
- Em revisoes, priorize padrao de commit, clareza da descricao e impacto da mudanca.
