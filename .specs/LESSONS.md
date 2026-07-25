# LESSONS — auto-maintained by scripts/lessons.py

> Machine-owned. Do NOT hand-edit. Changes are overwritten on the next `lessons.py` write.
> Canonical state lives in `.specs/lessons.json`. Edit lessons only via the script.
> promote_threshold=2 distinct features · window_days=45 · quarantine_threshold=2

## Confirmed (load these at Specify/Design)

Corroborated across multiple features. Safe to apply as guidance.

_none_

## Candidates (under observation — do NOT load as guidance yet)

Seen once or not yet corroborated. Tracked, not trusted.

### L-001 — pokemon-catalog-service precisa de testes unitarios com mock repository
- signal: `ac_gap` · recurrence: 1 feature(s) · scope: `testing` · harmful: 0
- features: production-readiness
- evidence: PR-25 (testing)
- last seen: 2026-07-25T20:48:22Z

### L-002 — auth-service precisa de testes table-driven para Signup e Login
- signal: `ac_gap` · recurrence: 1 feature(s) · scope: `testing` · harmful: 0
- features: production-readiness
- evidence: PR-26 (testing)
- last seen: 2026-07-25T20:48:22Z

### L-003 — Escritas de favoritos ainda via PostgresFavoriteRepository. Migração requer endpoint REST no catalog-service.
- signal: `spec_deviation` · recurrence: 1 feature(s) · scope: `architecture` · harmful: 0
- features: production-readiness
- evidence: PR-11 (architecture)
- last seen: 2026-07-25T20:48:22Z

## Quarantined (failed when applied — ignore)

A confirmed lesson that recurred alongside failure. Kept for the maintainer to review.

_none_
