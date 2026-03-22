# Instrucoes do Copilot para este repositorio

## Idioma
- Escreva respostas e explicacoes em portugues do Brasil.
- Gere mensagens de commit em portugues do Brasil.

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

## Colaboracao com CI
- Se o workflow de validacao de commits falhar, ajuste a mensagem para o padrao Conventional Commits.
- Em revisoes, priorize padrao de commit, clareza da descricao e impacto da mudanca.
