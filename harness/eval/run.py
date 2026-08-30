#!/usr/bin/env python3
"""
Harness de avaliacao das skills e AGENTS.md.

Uso:
  python3 harness/eval/run.py list                  # lista evals descobertos
  python3 harness/eval/run.py schema                # valida schema dos JSON
  python3 harness/eval/run.py run [opts]            # roda evals contra LLM

Opcoes de run:
  --skill NAME        filtra uma skill (ex.: go-style-combined)
  --case ID           filtra um caso (ex.: 1 ou 1.2)
  --provider P        anthropic | openai (default: env EVAL_PROVIDER ou anthropic)
  --model ID          model id (default: env EVAL_MODEL)
  --judge-model ID    model separado para judge
  --output DIR        diretorio de saida (default: harness/eval/results/<timestamp>)
  --dry-run           so imprime o que rodaria, sem chamar LLM
  --help              ajuda

Env:
  EVAL_API_KEY        chave da API do provider
  EVAL_PROVIDER       anthropic | openai
  EVAL_MODEL          model candidato
  EVAL_JUDGE_MODEL    model judge (opcional, default = MODEL)
  EVAL_BASE_URL       base URL para provider openai-compatible
  ANTHROPIC_API_KEY   alias para EVAL_API_KEY quando provider=anthropic
  OPENAI_API_KEY      alias para EVAL_API_KEY quando provider=openai
"""

from __future__ import annotations

import argparse
import datetime
import hashlib
import json
import os
import re
import sys
import time
import urllib.request
import urllib.error
from pathlib import Path

# ---------------------------------------------------------------------------
# Paths
# ---------------------------------------------------------------------------
REPO_ROOT = Path(__file__).resolve().parents[2]
SKILLS_GLOB = REPO_ROOT / ".github" / "skills"
HARNESS_EVALS = REPO_ROOT / "harness" / "evals"
RESULTS_BASE = Path(__file__).parent / "results"
AGENTS_MD = REPO_ROOT / "AGENTS.md"

# ---------------------------------------------------------------------------
# Discovery
# ---------------------------------------------------------------------------

def discover_eval_files() -> list[Path]:
    files: list[Path] = []
    # Skills: .github/skills/*/evals/evals.json
    if SKILLS_GLOB.exists():
        for eval_file in sorted(SKILLS_GLOB.glob("*/evals/evals.json")):
            files.append(eval_file)
    # Harness: harness/evals/*.json
    if HARNESS_EVALS.exists():
        for f in sorted(HARNESS_EVALS.glob("*.json")):
            files.append(f)
    return files


def eval_skill_name(eval_file: Path) -> str:
    # .github/skills/<skill>/evals/evals.json -> <skill>
    # harness/evals/<name>.json -> harness/<name>
    try:
        rel = eval_file.relative_to(REPO_ROOT)
    except ValueError:
        return eval_file.stem
    parts = rel.parts
    if len(parts) >= 4 and parts[0] == ".github" and parts[1] == "skills":
        return parts[2]
    if len(parts) >= 3 and parts[0] == "harness" and parts[1] == "evals":
        return f"harness/{eval_file.stem}"
    return eval_file.stem


def load_context(skill_name: str) -> str:
    """Monta contexto de sistema: AGENTS.md + SKILL.md + references."""
    chunks: list[str] = []
    if AGENTS_MD.exists():
        chunks.append(f"# AGENTS.md\n\n{AGENTS_MD.read_text(encoding='utf-8')}")
    # Skill SKILL.md + references
    if not skill_name.startswith("harness/"):
        skill_dir = SKILLS_GLOB / skill_name
        skill_md = skill_dir / "SKILL.md"
        if skill_md.exists():
            chunks.append(f"# Skill: {skill_name}\n\n{skill_md.read_text(encoding='utf-8')}")
        ref_dir = skill_dir / "references"
        if ref_dir.exists():
            for ref in sorted(ref_dir.glob("*.md")):
                chunks.append(f"# Reference: {ref.name}\n\n{ref.read_text(encoding='utf-8')}")
    else:
        # harness eval: AGENTS.md ja basta, mas inclui texto do eval como contexto adicional
        pass
    return "\n\n---\n\n".join(chunks)


# ---------------------------------------------------------------------------
# Schema validation
# ---------------------------------------------------------------------------

def _extract_cases(data):
    """Normaliza: array direto ou {evals: [...]} (schema samber legado)."""
    if isinstance(data, list):
        return data
    if isinstance(data, dict) and isinstance(data.get("evals"), list):
        return data["evals"]
    return None


def validate_eval_file(path: Path) -> list[str]:
    errors: list[str] = []
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as e:
        return [f"JSON invalido: {e}"]
    except OSError as e:
        return [f"erro de leitura: {e}"]

    cases = _extract_cases(data)
    if cases is None:
        return ["top-level deve ser array ou {evals: [...]}"]
    data = cases

    seen_ids: set[int] = set()
    for idx, case in enumerate(data):
        prefix = f"[{idx}]"
        if not isinstance(case, dict):
            errors.append(f"{prefix}: caso deve ser objeto")
            continue
        for field in ("id", "name", "description", "prompt", "trap", "assertions"):
            if field not in case:
                errors.append(f"{prefix}: campo ausente `{field}`")
        cid = case.get("id")
        if isinstance(cid, int):
            if cid in seen_ids:
                errors.append(f"{prefix}: id duplicado {cid}")
            seen_ids.add(cid)
        elif cid is not None:
            errors.append(f"{prefix}: `id` deve ser int")
        name = case.get("name")
        if isinstance(name, str) and not re.match(r"^[a-z0-9][a-z0-9\-]*$", name):
            errors.append(f"{prefix}: `name` deve ser kebab-case (letras, numeros, hifen)")

        assertions = case.get("assertions")
        if assertions is not None:
            if not isinstance(assertions, list) or len(assertions) == 0:
                errors.append(f"{prefix}: `assertions` deve ser array nao vazio")
            else:
                seen_aids: set[str] = set()
                for aidx, a in enumerate(assertions):
                    aprefix = f"{prefix}.assertions[{aidx}]"
                    if not isinstance(a, dict):
                        errors.append(f"{aprefix}: deve ser objeto")
                        continue
                    if "id" not in a or "text" not in a:
                        errors.append(f"{aprefix}: precisa `id` e `text`")
                    aid = a.get("id")
                    if isinstance(aid, str) and aid in seen_aids:
                        errors.append(f"{aprefix}: assertion id duplicado `{aid}`")
                    if isinstance(aid, str):
                        seen_aids.add(aid)
    return errors


# ---------------------------------------------------------------------------
# LLM clients (stdlib only)
# ---------------------------------------------------------------------------

def _api_key(provider: str) -> str | None:
    if provider == "anthropic":
        return os.environ.get("EVAL_API_KEY") or os.environ.get("ANTHROPIC_API_KEY")
    return os.environ.get("EVAL_API_KEY") or os.environ.get("OPENAI_API_KEY")


def call_anthropic(system: str, user: str, model: str, api_key: str, max_tokens: int = 4096) -> str:
    body = json.dumps({
        "model": model,
        "max_tokens": max_tokens,
        "system": system,
        "messages": [{"role": "user", "content": user}],
    }).encode()
    req = urllib.request.Request(
        "https://api.anthropic.com/v1/messages",
        data=body,
        headers={
            "Content-Type": "application/json",
            "x-api-key": api_key,
            "anthropic-version": "2023-06-01",
        },
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            data = json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        err_body = e.read().decode(errors="replace")
        raise RuntimeError(f"Anthropic API {e.code}: {err_body}") from e
    # Anthropic response: {content: [{type: text, text: ...}]}
    texts = [b["text"] for b in data.get("content", []) if b.get("type") == "text"]
    return "\n".join(texts)


def call_openai(system: str, user: str, model: str, api_key: str, base_url: str | None, max_tokens: int = 4096) -> str:
    base = (base_url or os.environ.get("EVAL_BASE_URL") or "https://api.openai.com/v1").rstrip("/")
    url = f"{base}/chat/completions"
    body = json.dumps({
        "model": model,
        "max_tokens": max_tokens,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
    }).encode()
    req = urllib.request.Request(
        url,
        data=body,
        headers={
            "Content-Type": "application/json",
            "Authorization": f"Bearer {api_key}",
        },
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            data = json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        err_body = e.read().decode(errors="replace")
        raise RuntimeError(f"OpenAI API {e.code}: {err_body}") from e
    choices = data.get("choices", [])
    if not choices:
        raise RuntimeError(f"OpenAI resposta sem choices: {data}")
    return choices[0]["message"]["content"] or ""


def call_llm(system: str, user: str, model: str, provider: str) -> str:
    key = _api_key(provider)
    if not key:
        raise RuntimeError(
            f"API key ausente para provider={provider}. "
            f"Defina EVAL_API_KEY (ou ANTHROPIC_API_KEY/OPENAI_API_KEY)."
        )
    if provider == "anthropic":
        return call_anthropic(system, user, model, key)
    return call_openai(system, user, model, key, os.environ.get("EVAL_BASE_URL"))


# ---------------------------------------------------------------------------
# Judge
# ---------------------------------------------------------------------------

JUDGE_SYSTEM = """Voce e um avaliador rigoroso de codigo Go. Recebe OUTPUT de um candidato e ASSERTIONS a verificar.

Responda APENAS com JSON valido no formato:
{
  \"assertions\": {\"<id>\": {\"pass\": true|false, \"reason\": \"...\"}},
  \"trap_hit\": true|false,
  \"trap_reason\": \"...\"
}

Regras:
- pass=true apenas se a assertion foi CLARAMENTE atendida no output (nao inferencia generosa).
- trap_hit=true se o output caiu na armadilha descrita (TRAP).
- Seja estrito: codigo que quase atende mas viola o detalhe especifico = fail.
"""

def judge_output(output: str, case: dict, judge_model: str, provider: str) -> dict:
    assertions_text = "\n".join(f"- {a['id']}: {a['text']}" for a in case["assertions"])
    user = (
        f"TRAP (armadilha a evitar): {case['trap']}\n\n"
        f"ASSERTIONS:\n{assertions_text}\n\n"
        f"OUTPUT DO CANDIDATO:\n```\n{output[:12000]}\n```\n\n"
        f"Avalie cada assertion e se a trap foi acionada. Responda APENAS JSON."
    )
    raw = call_llm(JUDGE_SYSTEM, user, judge_model, provider)
    # Extrai JSON do raw (pode vir com markdown fence)
    m = re.search(r"\{.*\}", raw, re.DOTALL)
    if m:
        try:
            return json.loads(m.group(0))
        except json.JSONDecodeError:
            pass
    # fallback: tenta parse direto
    try:
        return json.loads(raw)
    except json.JSONDecodeError:
        return {"raw": raw, "parse_error": True, "assertions": {}, "trap_hit": None}


# ---------------------------------------------------------------------------
# Commands
# ---------------------------------------------------------------------------

def cmd_list() -> int:
    files = discover_eval_files()
    if not files:
        print("Nenhum eval encontrado.")
        return 0
    total_cases = 0
    total_assertions = 0
    print(f"{'skill':<32} {'arquivo':<60} {'casos':>5} {'asserts':>7}")
    print("-" * 110)
    for f in files:
        try:
            raw = json.loads(f.read_text(encoding="utf-8"))
            cases = _extract_cases(raw)
            if cases is None:
                n_cases, n_asserts = -1, -1
            else:
                n_cases = len(cases)
                n_asserts = sum(len(c.get("assertions", [])) for c in cases)
        except Exception:
            n_cases = -1
            n_asserts = -1
        rel = f.relative_to(REPO_ROOT)
        skill = eval_skill_name(f)
        print(f"{skill:<32} {str(rel):<60} {n_cases:>5} {n_asserts:>7}")
        if n_cases > 0:
            total_cases += n_cases
            total_assertions += n_asserts
    print("-" * 110)
    print(f"{'TOTAL':<32} {'':<60} {total_cases:>5} {total_assertions:>7}")
    return 0


def cmd_schema() -> int:
    files = discover_eval_files()
    if not files:
        print("Nenhum eval encontrado.")
        return 0
    ok = True
    for f in files:
        errs = validate_eval_file(f)
        rel = f.relative_to(REPO_ROOT)
        if errs:
            ok = False
            print(f"FAIL  {rel}")
            for e in errs:
                print(f"  - {e}")
        else:
            raw = json.loads(f.read_text(encoding="utf-8"))
            cases = _extract_cases(raw) or []
            print(f"OK    {rel} ({len(cases)} casos)")
    return 0 if ok else 1


def cmd_run(args: argparse.Namespace) -> int:
    provider = args.provider or os.environ.get("EVAL_PROVIDER") or "anthropic"
    model = args.model or os.environ.get("EVAL_MODEL") or (
        "claude-sonnet-4-20250514" if provider == "anthropic" else "gpt-4o"
    )
    judge_model = args.judge_model or os.environ.get("EVAL_JUDGE_MODEL") or model

    files = discover_eval_files()
    if args.skill:
        files = [f for f in files if eval_skill_name(f) == args.skill]
        if not files:
            print(f"Nenhum eval para skill={args.skill}", file=sys.stderr)
            return 1

    # Filtro de caso especifico (--case: "1" ou "1.1")
    case_filter = args.case

    # Valida schema primeiro
    # Valida; arquivos com schema invalido sao ignorados em `run` (warn) — use `schema` para ver todos os FAILs
    valid_files: list[Path] = []
    for f in files:
        errs = validate_eval_file(f)
        if errs:
            print(f"WARN: ignorando {f} (schema invalido): {errs[:2]}", file=sys.stderr)
            continue
        valid_files.append(f)
    if not valid_files:
        print("Nenhum eval valido encontrado.", file=sys.stderr)
        return 1
    files = valid_files

    if args.dry_run:
        for f in files:
            raw = json.loads(f.read_text(encoding="utf-8"))
            cases = _extract_cases(raw) or []
            skill = eval_skill_name(f)
            for case in cases:
                if case_filter and str(case["id"]) != str(case_filter):
                    continue
                print(f"would run: {skill} case {case['id']} ({case['name']})")
        return 0

    # Verifica API key antes de comecar
    if not _api_key(provider):
        print(
            f"ERRO: API key ausente para provider={provider}. "
            f"Defina EVAL_API_KEY (ou ANTHROPIC_API_KEY/OPENAI_API_KEY).",
            file=sys.stderr,
        )
        return 1

    # Diretorio de saida
    ts = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H%M%SZ")
    out_dir = Path(args.output) if args.output else (RESULTS_BASE / ts)
    out_dir.mkdir(parents=True, exist_ok=True)
    print(f"Provider={provider} model={model} judge={judge_model}")
    print(f"Saida: {out_dir}")

    all_results: list[dict] = []
    for eval_file in files:
        raw = json.loads(eval_file.read_text(encoding="utf-8"))
        cases = _extract_cases(raw) or []
        skill = eval_skill_name(eval_file)
        context = load_context(skill)
        # System prompt enxuto: contexto + instrucao
        candidate_system = (
            f"Voce e um engenheiro Go seguindo as regras abaixo. "
            f"Responda APENAS com codigo Go (e breve explicacao se necessario), "
            f"seguindo estritamente o contexto.\n\n{context}"
        )

        for case in cases:
            if case_filter and str(case["id"]) != str(case_filter):
                continue

            cid = case["id"]
            print(f"\n[{skill} #{cid} {case['name']}] prompt -> candidato...")
            t0 = time.time()
            try:
                candidate_output = call_llm(candidate_system, case["prompt"], model, provider)
            except Exception as e:
                print(f"  ERRO candidato: {e}", file=sys.stderr)
                result = {
                    "skill": skill,
                    "eval_file": str(eval_file.relative_to(REPO_ROOT)),
                    "case_id": cid,
                    "case_name": case["name"],
                    "error": str(e),
                    "phase": "candidate",
                }
                all_results.append(result)
                continue

            elapsed_c = time.time() - t0
            print(f"  candidato OK ({elapsed_c:.1f}s, {len(candidate_output)} chars) -> judge...")

            t1 = time.time()
            try:
                verdict = judge_output(candidate_output, case, judge_model, provider)
            except Exception as e:
                print(f"  ERRO judge: {e}", file=sys.stderr)
                result = {
                    "skill": skill,
                    "eval_file": str(eval_file.relative_to(REPO_ROOT)),
                    "case_id": cid,
                    "case_name": case["name"],
                    "candidate_output": candidate_output,
                    "error": str(e),
                    "phase": "judge",
                }
                all_results.append(result)
                continue

            elapsed_j = time.time() - t1

            # Normaliza verdict
            v_assertions = verdict.get("assertions", {})
            passes = sum(1 for v in v_assertions.values() if isinstance(v, dict) and v.get("pass") is True)
            total = len(case["assertions"])
            trap_hit = verdict.get("trap_hit")

            result = {
                "skill": skill,
                "eval_file": str(eval_file.relative_to(REPO_ROOT)),
                "case_id": cid,
                "case_name": case["name"],
                "description": case["description"],
                "trap": case["trap"],
                "assertions": case["assertions"],
                "candidate_output": candidate_output,
                "verdict": verdict,
                "passes": passes,
                "total": total,
                "trap_hit": trap_hit,
                "elapsed_candidate_s": round(elapsed_c, 2),
                "elapsed_judge_s": round(elapsed_j, 2),
            }
            all_results.append(result)

            status = "PASS" if passes == total and trap_hit is False else "FAIL"
            print(f"  judge OK ({elapsed_j:.1f}s) -> {status} {passes}/{total} trap_hit={trap_hit}")

            # Throttle leve entre casos
            time.sleep(0.5)

    # Escreve resultados
    results_json = out_dir / "results.json"
    results_json.write_text(json.dumps(all_results, ensure_ascii=False, indent=2), encoding="utf-8")

    # Relatorio markdown
    md_lines: list[str] = [f"# Eval harness — {ts}", ""]
    md_lines.append(f"Provider: `{provider}`  Model: `{model}`  Judge: `{judge_model}`")
    md_lines.append("")
    # Sumario por skill
    from collections import defaultdict
    by_skill: dict[str, list[dict]] = defaultdict(list)
    for r in all_results:
        by_skill[r["skill"]].append(r)

    md_lines.append("| skill | casos | pass | trap_hit |")
    md_lines.append("|-------|-------|------|----------|")
    for skill, cases in sorted(by_skill.items()):
        total_c = len(cases)
        passed_c = sum(1 for c in cases if c.get("passes") == c.get("total") and c.get("trap_hit") is False)
        trap_hits = sum(1 for c in cases if c.get("trap_hit") is True)
        md_lines.append(f"| {skill} | {total_c} | {passed_c}/{total_c} | {trap_hits} |")
    # Global
    total_c = len(all_results)
    passed_c = sum(1 for r in all_results if r.get("passes") == r.get("total") and r.get("trap_hit") is False)
    md_lines.append(f"| **TOTAL** | **{total_c}** | **{passed_c}/{total_c}** |  |")
    md_lines.append("")

    for skill, cases in sorted(by_skill.items()):
        md_lines.append(f"## {skill}")
        for c in sorted(cases, key=lambda x: x["case_id"]):
            cid = c["case_id"]
            name = c["case_name"]
            passes = c.get("passes", "?")
            total_a = c.get("total", "?")
            trap_hit = c.get("trap_hit")
            err = c.get("error")
            if err:
                md_lines.append(f"- **#{cid} {name}** — ERRO ({c.get('phase')}): {err}")
                continue
            icon = "✅" if passes == total_a and trap_hit is False else "❌"
            md_lines.append(f"- {icon} **#{cid} {name}** — {passes}/{total_a} trap_hit={trap_hit}")
            verdict = c.get("verdict", {})
            for aid, v in verdict.get("assertions", {}).items():
                if isinstance(v, dict):
                    picon = "✓" if v.get("pass") else "✗"
                    md_lines.append(f"  - {picon} {aid}: {v.get('reason','')}")
            if verdict.get("trap_reason"):
                md_lines.append(f"  - trap: {verdict['trap_reason']}")
        md_lines.append("")

    report_md = out_dir / "report.md"
    report_md.write_text("\n".join(md_lines), encoding="utf-8")
    print(f"\nResultados: {results_json}")
    print(f"Relatorio:  {report_md}")
    # Sumario stdout
    print("\n" + "\n".join(md_lines[:20]))

    # Exit code: falha se algum caso nao passou
    failed = sum(1 for r in all_results if r.get("passes") != r.get("total") or r.get("trap_hit") is True or "error" in r)
    return 0 if failed == 0 else 1


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def main() -> int:
    parser = argparse.ArgumentParser(description="Harness de avaliacao de skills")
    sub = parser.add_subparsers(dest="cmd", required=True)

    p_list = sub.add_parser("list", help="lista evals descobertos")
    p_schema = sub.add_parser("schema", help="valida schema dos JSON")
    p_run = sub.add_parser("run", help="roda evals contra LLM")
    p_run.add_argument("--skill", help="filtra skill (ex.: go-style-combined)")
    p_run.add_argument("--case", help="filtra caso por id (ex.: 1)")
    p_run.add_argument("--provider", choices=["anthropic", "openai"], help="provider LLM")
    p_run.add_argument("--model", help="model id candidato")
    p_run.add_argument("--judge-model", help="model id judge")
    p_run.add_argument("--output", help="diretorio de saida")
    p_run.add_argument("--dry-run", action="store_true", help="so imprime o que rodaria")

    args = parser.parse_args()
    if args.cmd == "list":
        return cmd_list()
    if args.cmd == "schema":
        return cmd_schema()
    if args.cmd == "run":
        return cmd_run(args)
    parser.print_help()
    return 1


if __name__ == "__main__":
    sys.exit(main())
