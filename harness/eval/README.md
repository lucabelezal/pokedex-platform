# Harness de avaliacao — skills e AGENTS.md

Valida se um agente/LLM segue as skills Go do projeto e as convencoes do `AGENTS.md`.

## Evals

| Local | O que cobre |
|-------|-------------|
| `.github/skills/<skill>/evals/evals.json` | 7 skills proprias do projeto (go-style-combined, go-api-design, go-error-handling, go-security-audit, go-test-quality, go-test-table-driven, go-architecture-review) |
| `harness/evals/agents-md.json` | Convencoes do `AGENTS.md` (pt-BR, hexagonal, commits, etc.) |

Formato de cada caso:

```json
{
  "id": 1,
  "name": "kebab-case",
  "description": "...",
  "prompt": "instrucao enviada ao candidato",
  "trap": "armadilha plausivel que o candidato pode cair",
  "assertions": [
    { "id": "1.1", "text": "criterio verificavel no output" }
  ]
}
```

O `trap` descreve o erro que um modelo desatento comete (ex.: `GetCount` em vez de `Count`, `Id` em vez de `ID`, `var s []T` que serializa `null`). O judge verifica `trap_hit`.

## Runner

`harness/eval/run.py` — Python 3 stdlib (sem dependencias). Suporta Anthropic e OpenAI-compatible.

### Comandos (sem API key)

```bash
# lista evals descobertos
python3 harness/eval/run.py list

# valida schema dos JSON
python3 harness/eval/run.py schema

# dry-run (mostra o que rodaria)
python3 harness/eval/run.py run --dry-run
python3 harness/eval/run.py run --skill go-style-combined --dry-run
```

### Rodar contra LLM (precisa API key)

```bash
# Anthropic (default)
export ANTHROPIC_API_KEY=sk-ant-...
python3 harness/eval/run.py run

# OpenAI
export OPENAI_API_KEY=sk-...
python3 harness/eval/run.py run --provider openai --model gpt-4o

# Provider openai-compatible (ex.: local)
export EVAL_API_KEY=...
export EVAL_BASE_URL=http://localhost:11434/v1
python3 harness/eval/run.py run --provider openai --model llama3

# Filtrar
python3 harness/eval/run.py run --skill go-style-combined --case 1
python3 harness/eval/run.py run --skill harness/agents-md
```

Variaveis de ambiente:

| Var | Descricao | Default |
|-----|-----------|---------|
| `EVAL_PROVIDER` | `anthropic` ou `openai` | `anthropic` |
| `EVAL_API_KEY` | chave do provider | `ANTHROPIC_API_KEY` / `OPENAI_API_KEY` como alias |
| `EVAL_MODEL` | model candidato | `claude-sonnet-4-20250514` / `gpt-4o` |
| `EVAL_JUDGE_MODEL` | model judge (opcional) | mesmo do candidato |
| `EVAL_BASE_URL` | base URL para openai-compatible | `https://api.openai.com/v1` |

### Saida

```
harness/eval/results/<timestamp>/
  results.json   # dados brutos por caso (candidate_output + verdict)
  report.md      # sumario por skill + detalhe por assertion
```

Exit code: `0` se todos os casos passaram e nenhuma trap foi acionada; `1` caso contrario.

### Makefile

```bash
make eval              # roda harness completo (precisa API key)
```

## Adicionar novos casos

1. Edite o `evals.json` da skill relevante (ou `harness/evals/agents-md.json` para convencoes gerais).
2. Siga o schema: `id` int unico, `name` kebab-case, `trap` descreve o erro a evitar, `assertions[].id` no formato `<case>.<n>`.
3. Valide: `python3 harness/eval/run.py schema`.
4. Teste dry-run: `python3 harness/eval/run.py run --skill <skill> --dry-run`.
