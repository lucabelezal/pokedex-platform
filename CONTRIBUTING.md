# Guia de Contribuicao

## Idioma
- Escreva issues, PRs e mensagens de commit em portugues do Brasil.

## Convencao de Commit
- Use sempre Conventional Commits.
- Formato: `tipo(escopo-opcional): descricao curta em portugues`.
- Tipos: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `build`, `perf`, `revert`.

Exemplos:
- `feat(favoritos): bloquear remocao sem autenticacao`
- `fix(bff): corrigir validacao de token jwt`
- `docs(readme): atualizar fluxo de execucao local`

## Hook local de commit-msg
Para validar o formato antes do push:

```bash
git config core.hooksPath .githooks
chmod +x .githooks/commit-msg
```

## CI no GitHub Actions
Este repositorio valida automaticamente:
- Mensagens de commit (Conventional Commits)
- Titulo de PR (Conventional Commits)
- Build, test, vet e lint para modulos Go

## Checklist antes de abrir PR
- Rode build e testes dos modulos alterados
- Garanta que os commits estao no padrao
- Atualize documentacao se houver mudanca de comportamento
