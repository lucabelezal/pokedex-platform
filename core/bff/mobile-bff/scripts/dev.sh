#!/bin/bash

# Inicializador do servidor em desenvolvimento com suporte a PostgreSQL
# Uso: ./scripts/dev.sh [local|postgres]

set -e

MODE="${1:-local}"
PORT="${MOBILE_BFF_PORT:-8080}"

echo "Iniciando mobile-bff no modo $MODE na porta $PORT..."

case "$MODE" in
  local)
    echo "Usando repositorios mock (modo local)"
    MOBILE_BFF_PORT=$PORT go run ./cmd/server/main.go
    ;;
  postgres)
    echo "Conectando ao PostgreSQL..."
    # Requer docker-compose.test.yml em execucao
    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pokedex"
    MOBILE_BFF_PORT=$PORT go run ./cmd/server/main.go
    ;;
  *)
    echo "Modo desconhecido: $MODE"
    echo "Uso: ./scripts/dev.sh [local|postgres]"
    exit 1
    ;;
esac
