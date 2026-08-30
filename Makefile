.PHONY: help verify-bff-env check-bff-env doctor up down restart logs ps health home detail bff-run-local db-migrate-auth-hardening test test-unit test-integration test-race coverage lint ci eval eval-run db-test-up db-test-down

COMPOSE_FILE=core/docker-compose.yml
TEST_COMPOSE_FILE=core/docker-compose.test.yml
PROJECT_NAME=pokedex

# Modulos Go da plataforma (cada um com seu proprio go.mod)
MODULES=core/app/pokemon-catalog-service core/app/auth-service core/bff/mobile-bff core/infra/postgres/json2sql

# Banco de teste usado pelos testes de integracao (espelha core/docker-compose.test.yml)
TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5433/pokedex_test?sslmode=disable

help:
	@echo "Comandos disponiveis:"
	@echo "  make up              - Sobe toda a plataforma em background"
	@echo "  make down            - Derruba toda a plataforma"
	@echo "  make restart         - Reinicia a plataforma"
	@echo "  make logs            - Mostra logs da plataforma"
	@echo "  make ps              - Lista containers da plataforma"
	@echo "  make health          - Verifica health do BFF via gateway"
	@echo "  make home            - Consulta endpoint Home via gateway"
	@echo "  make detail          - Consulta detalhe do Pokemon #1 via gateway"
	@echo "  make verify-bff-env  - Mostra status da variavel POKEMON_CATALOG_SERVICE_URL"
	@echo "  make check-bff-env   - Valida variavel e falha se ausente"
	@echo "  make doctor          - Diagnostico rapido do ambiente local"
	@echo "  make bff-run-local   - Roda mobile-bff local exigindo POKEMON_CATALOG_SERVICE_URL"
	@echo "  make db-migrate-auth-hardening - Aplica migration incremental de tokens no Postgres"
	@echo ""
	@echo "  Testes (todos os modulos):"
	@echo "  make test            - Unit + integracao (integração auto-skip sem DB)"
	@echo "  make test-unit       - Apenas unit (forca skip da integracao)"
	@echo "  make test-race       - Unit com detector de corrida (-race)"
	@echo "  make test-integration- Sobe postgres-test e roda testes de integracao"
	@echo "  make coverage        - Cobertura por modulo (go tool cover)"
	@echo "  make lint            - gofmt + go vet em todos os modulos"
	@echo "  make ci              - build + lint + test (espelha pipeline CI)"
	@echo ""
	@echo "  Harness de avaliacao:"
	@echo "  make eval            - Roda harness de eval das skills (harness/eval/run.py)"

verify-bff-env:
	@if [ -z "$$POKEMON_CATALOG_SERVICE_URL" ]; then \
		echo "AVISO: POKEMON_CATALOG_SERVICE_URL nao esta configurada no shell atual."; \
		echo "Resolucao recomendada (execucao local fora do compose):"; \
		echo "  export POKEMON_CATALOG_SERVICE_URL=http://localhost:8081"; \
		echo "Obs: No docker compose da plataforma existe valor padrao interno."; \
	else \
		echo "OK: POKEMON_CATALOG_SERVICE_URL=$$POKEMON_CATALOG_SERVICE_URL"; \
	fi

check-bff-env:
	@if [ -z "$$POKEMON_CATALOG_SERVICE_URL" ]; then \
		echo "ERRO: POKEMON_CATALOG_SERVICE_URL nao configurada."; \
		echo "Como resolver:"; \
		echo "  export POKEMON_CATALOG_SERVICE_URL=http://localhost:8081"; \
		exit 1; \
	fi

doctor:
	@echo "Diagnostico do ambiente local"
	@echo ""
	@echo "[1/4] Ferramentas essenciais"
	@if command -v docker >/dev/null 2>&1; then \
		echo "  OK docker"; \
	else \
		echo "  ERRO docker nao encontrado"; \
	fi
	@if docker compose version >/dev/null 2>&1; then \
		echo "  OK docker compose"; \
	else \
		echo "  ERRO docker compose nao disponivel"; \
	fi
	@if command -v go >/dev/null 2>&1; then \
		echo "  OK go"; \
	else \
		echo "  AVISO go nao encontrado (necessario para bff-run-local e testes)"; \
	fi
	@echo ""
	@echo "[2/4] Variavel obrigatoria do BFF fora do compose"
	@if [ -z "$$POKEMON_CATALOG_SERVICE_URL" ]; then \
		echo "  AVISO POKEMON_CATALOG_SERVICE_URL ausente no shell atual"; \
		echo "  Resolucao: export POKEMON_CATALOG_SERVICE_URL=http://localhost:8081"; \
	else \
		echo "  OK POKEMON_CATALOG_SERVICE_URL=$$POKEMON_CATALOG_SERVICE_URL"; \
	fi
	@echo ""
	@echo "[3/4] Arquivo de compose"
	@if [ -f "$(COMPOSE_FILE)" ]; then \
		echo "  OK $(COMPOSE_FILE)"; \
	else \
		echo "  ERRO arquivo $(COMPOSE_FILE) nao encontrado"; \
	fi
	@echo ""
	@echo "[4/4] Portas de runtime"
	@if command -v lsof >/dev/null 2>&1; then \
		if lsof -i :8000 >/dev/null 2>&1; then echo "  Porta 8000 em uso (Kong)"; else echo "  Porta 8000 livre"; fi; \
		if lsof -i :8001 >/dev/null 2>&1; then echo "  Porta 8001 em uso (Kong Admin)"; else echo "  Porta 8001 livre"; fi; \
		if lsof -i :8080 >/dev/null 2>&1; then echo "  Porta 8080 em uso (mobile-bff)"; else echo "  Porta 8080 livre"; fi; \
		if lsof -i :8081 >/dev/null 2>&1; then echo "  Porta 8081 em uso (pokemon-catalog-service)"; else echo "  Porta 8081 livre"; fi; \
		if lsof -i :8082 >/dev/null 2>&1; then echo "  Porta 8082 em uso (auth-service)"; else echo "  Porta 8082 livre"; fi; \
	else \
		echo "  AVISO lsof nao encontrado; nao foi possivel verificar portas"; \
	fi

up: verify-bff-env
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) up --build -d

down:
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) down

restart: down up

logs:
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) logs -f --tail=200

ps:
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) ps

health:
	@curl -fsS http://localhost:8000/bff/health | cat

home:
	@curl -fsS "http://localhost:8000/v1/home" | cat

detail:
	@curl -fsS "http://localhost:8000/v1/pokemons/1/details" | cat

bff-run-local: check-bff-env
	@cd core/bff/mobile-bff && MOBILE_BFF_PORT=8080 go run ./cmd/server/main.go

db-migrate-auth-hardening:
	@./core/infra/postgres/migrations/apply.sh

# ---------------------------------------------------------------------------
# Harness de teste — roda o mesmo comando em todos os modulos Go
# ---------------------------------------------------------------------------

test: test-unit

test-unit:
	@for m in $(MODULES); do \
		echo "=== $$m (unit) ==="; \
		( cd $$m && unset DATABASE_URL TEST_DATABASE_URL && go test -race ./... ) || exit 1; \
	done

test-race:
	@for m in $(MODULES); do \
		echo "=== $$m (race) ==="; \
		( cd $$m && unset DATABASE_URL TEST_DATABASE_URL && go test -race -count=1 ./... ) || exit 1; \
	done

test-integration: db-test-up
	@echo "Rodando testes de integracao (TEST_DATABASE_URL=$(TEST_DATABASE_URL))"
	@for m in $(MODULES); do \
		echo "=== $$m (integracao) ==="; \
		( cd $$m && TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -race -count=1 -v ./... -timeout 90s ) || exit 1; \
	done
	@$(MAKE) db-test-down

coverage:
	@for m in $(MODULES); do \
		echo "=== $$m (coverage) ==="; \
		( cd $$m && unset DATABASE_URL TEST_DATABASE_URL && go test -race -covermode=atomic -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1 ) || exit 1; \
	done

lint:
	@for m in $(MODULES); do \
		echo "=== $$m (lint) ==="; \
		( cd $$m && out="$$(gofmt -l .)" && [ -z "$$out" ] || { echo "Arquivos fora do gofmt:"; echo "$$out"; exit 1; } && go vet ./... ) || exit 1; \
	done

ci:
	@for m in $(MODULES); do \
		echo "=== build $$m ==="; \
		( cd $$m && go build ./... ) || exit 1; \
	done
	@$(MAKE) lint
	@$(MAKE) test

db-test-up:
	@docker compose -f $(TEST_COMPOSE_FILE) up -d
	@echo "Aguardando postgres-test ficar pronto..."
	@i=0; until docker compose -f $(TEST_COMPOSE_FILE) exec -T postgres-test pg_isready -U postgres -d pokedex_test >/dev/null 2>&1; do \
		i=$$((i+1)); \
		if [ $$i -ge 30 ]; then echo "ERRO: postgres-test nao ficou pronto em 30s"; exit 1; fi; \
		sleep 1; \
	done; \
	echo "OK: postgres-test pronto"

db-test-down:
	@docker compose -f $(TEST_COMPOSE_FILE) down

# ---------------------------------------------------------------------------
# Harness de avaliacao das skills
# ---------------------------------------------------------------------------

eval:
	@python3 harness/eval/run.py run --dry-run 2>&1 | grep -v "^WARN:"
	@echo ""
	@echo "Validacao detalhada: python3 harness/eval/run.py schema"
	@echo "Rodar contra LLM:    EVAL_API_KEY=... python3 harness/eval/run.py run"

eval-run:
	@python3 harness/eval/run.py run
