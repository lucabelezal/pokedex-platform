---
description: Planeja e executa refactor Go seguro (Plan → Stage → Land) com blast-radius via gopls, PRs atomicos e safety-net de testes
---

Alvo do refactor: $ARGUMENTS

> Se $ARGUMENTS estiver vazio, pergunte qual codigo/area refatorar e qual objetivo (ex.: "extrair validacao de ProcessOrder", "quebrar ciclo billing↔orders").

## Skills principais

Carregue nesta ordem:

1. **golang-refactoring** (`.claude/skills/golang-refactoring/SKILL.md` ou `.agents/skills/golang-refactoring/SKILL.md`) — workflow canonico Plan→Stage→Land
2. **golang-gopls** — mecanica de gopls (references, call hierarchy, rename/inline/extract) — unica fonte de como obter gopls/MCP
3. **go-style-combined** — estilo canonico do projeto (receivers, initialisms, errors %w, make([]T,0), time.Duration, etc.)
4. Conforme o caso: `golang-naming` (renames), `golang-project-layout` (splits), `golang-modernize` (idioms), `golang-code-style` (control flow), `golang-design-patterns` (DI/patterns), `golang-safety`/`golang-security` (mudancas com risco logico)

Referencias obrigatorias: `references/workflow.md`, `references/safety-net.md`, `references/go-tooling.md`, `references/catalog.md` da skill golang-refactoring.

## Fluxo

### Fase 1 — Plan (gate obrigatorio, sem tocar codigo)

1. **Mapeie blast radius com gopls** antes de qualquer edicao:
   - `gopls references`, call hierarchy (ida e volta), workspace symbol, package API
   - Verifique exported API que modulos externos possam depender
2. **Monte inventario de refactor** — uma linha por mudanca atomica:

   | Transform | Arquivos / call sites | Risco | S/B |
   |-----------|----------------------|-------|-----|
   | ex.: Rename `Client.Send`→`Publish` | `pkg/client/*.go`, 14 sites, 3 pkgs | Low | S |
   | ex.: Extrair `validateOrder` | `internal/orders/process.go` | Low | S |
   | ... | ... | ... | ... |

   - Risco: Low/Medium/High (tabela da skill)
   - S/B: Structural vs Behavioral (Beck) — **nunca misture S e B no mesmo PR/commit**
   - Mesmo dentro de S: move + otimizacao no mesmo arquivo = 2 PRs sequenciais (verificacoes diferentes)
3. **Ordene em 3 dimensoes:**
   - (a) Beck: structural primeiro, behavioral por ultimo
   - (b) Conflito: arquivos/simbolos compartilhados → sequencial; file-disjoint → pode paralelizar
   - (c) Dependencia: quebra de ciclo / pacote extraido / alias → prerequisite primeiro
   - Rename workspace-wide = barreira: sozinho, sem concorrencia
4. **Apresente inventario + plano de PRs e AGUARDE sign-off explicito** (via question tool). Nao toque codigo antes.

### Fase 2 — Stage & Land (humano no loop)

Modelo de git para refactors staged:

1. Crie branch de integracao `refactor/<topico>` a partir de `main`; semeie com `// REFACTOR(step N): ...` markers
2. Para cada linha do inventario, na ordem definida:
   - Dispare **um sub-agent por mudanca** em worktree isolado, branch `refactor/<topico>-<step>`
   - Aplique **uma** mudanca, preferindo transform tool-driven (gopls Rename/Inline, `gofmt -r`, `eg`, `gopatch`, `go/analysis` fixers) sobre edicao manual
   - Safety net: se cobertura do blast radius for LOW/MEDIUM, adicione testes antes (characterization tests) — veja `safety-net.md`
   - Verifique: `go build ./... && go vet ./... && go test ./...` (+ `-race` se concorrencia, `benchstat` se hot path)
   - Abra PR com `gh pr create --base refactor/<topico>` — **pronto para review, nao draft**; 100–500 linhas por PR
   - Sub-agent reporta resumo compacto (pass/fail, PR link); orquestrador mantem so o resumo
3. Humano revisa e mergeia cada PR em `refactor/<topico>` (S = review rapido, B = review completo; para B com logica, carregue `golang-security`+`golang-safety`)
4. Apos todas as linhas landarem e TODO-markers limpos, abra PR final `refactor/<topico>` → `main` **como draft** para revisao final deliberada

**Nunca mergeie PR intermediario direto em `main`.**

### Fase 3 — Simple-sweep (atalho)

Se for **uma unica transform mecanica tree-wide** sem dependencias (ex.: `gofmt -r`, `modernize` fixer): pode usar `ultracode`/single-pass, verifique green e commite. Nao use para multi-step com checkpoints humanos.

### Verificacao continua

- Harness do projeto: `make lint` (`gofmt`+`go vet`), `make test`/`make test-race`, `make ci`
- Eval de estilo: `python3 harness/eval/run.py run --skill go-style-combined --dry-run` (ou `make eval`)
- Para renames: gopls recusa em shadowing/interface-break — confie no tool; mas valide reflexao (struct tags, templates) com testes

## Ao retomar

Leia inventario/PRs em `refactor/<topico>` e reconcilie com `git status` e PRs abertos antes de propor proximo passo.
