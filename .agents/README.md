# .agents - Governanca Local

Este diretorio contem customizacoes locais de agente para o workspace da Plataforma Pokedex.

## Precedencia De Regras

1. Fonte principal de regras do projeto: `.github/copilot-instructions.md`.
2. Regras em `.agents/rules` sao complementares e nao devem contradizer a fonte principal.
3. Skills em `.agents/skills` sao carregadas sob demanda; nao sao regras always-on.

## Politica De Conflito

- Em caso de conflito, prevalece `.github/copilot-instructions.md`.
- Se uma skill sugerir padrao diferente do projeto, seguir o padrao do projeto.
- Ajustes de skill devem ser feitos no proprio workspace para preservar compatibilidade com o repositorio.

## Escopo Atual

- `rules/`: regra local de Conventional Commits e idioma PT-BR.
- `skills/`: lote piloto de skills Go selecionadas para avaliacao incremental.
- `workflows/`: reservado para fluxos de execucao do time.

## Objetivo Do Piloto

Avaliar ganho de qualidade e consistencia com baixo risco de ruido de contexto antes de adotar o pacote completo de skills Go.
