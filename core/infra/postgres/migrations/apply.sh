#!/usr/bin/env bash
set -euo pipefail

MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DATABASE_URL="${DATABASE_URL:-${1:-}}"

if [[ -z "${DATABASE_URL}" ]]; then
  echo "ERRO: DATABASE_URL nao configurada."
  echo "Uso: DATABASE_URL=postgres://usuario:senha@host:5432/db?sslmode=disable ./core/infra/postgres/migrations/apply.sh"
  exit 1
fi

if ! command -v psql >/dev/null 2>&1; then
  echo "ERRO: psql nao encontrado no PATH."
  exit 1
fi

echo "Aplicando migrations SQL em ${MIGRATIONS_DIR}"
for migration in "${MIGRATIONS_DIR}"/*.sql; do
  [[ -f "${migration}" ]] || continue
  echo " - $(basename "${migration}")"
  psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -f "${migration}"
done

echo "Migrations aplicadas com sucesso."
