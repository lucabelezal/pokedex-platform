#!/bin/bash

# Script de preparo para testes de integracao
# Uso: ./scripts/test-integration.sh [up|down|test|clean]

set -e

COMPOSE_FILE="docker-compose.test.yml"
TEST_DB_URL="postgres://postgres:postgres@localhost:5433/pokedex_test"

case "${1:-test}" in
  up)
    echo "Iniciando infraestrutura de testes..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "Aguardando banco ficar pronto..."
    sleep 10
    docker-compose -f "$COMPOSE_FILE" exec -T postgres-test psql -U postgres -d pokedex_test -f /docker-entrypoint-initdb.d/001_create_tables.sql
    docker-compose -f "$COMPOSE_FILE" exec -T postgres-test psql -U postgres -d pokedex_test -f /docker-entrypoint-initdb.d/002_seed_data.sql
    echo "Infraestrutura de testes pronta"
    ;;
  down)
    echo "Parando infraestrutura de testes..."
    docker-compose -f "$COMPOSE_FILE" down
    ;;
  test)
    echo "Executando testes de integracao..."
    DATABASE_URL="$TEST_DB_URL" go test -v ./tests/integration -timeout 30s
    ;;
  clean)
    echo "Limpando ambiente..."
    docker-compose -f "$COMPOSE_FILE" down -v
    ;;
  *)
    echo "Uso: ./scripts/test-integration.sh [up|down|test|clean]"
    exit 1
    ;;
esac
