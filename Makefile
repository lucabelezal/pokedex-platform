.PHONY: help verify-bff-env check-bff-env doctor up down restart logs ps health home detail bff-run-local

COMPOSE_FILE=core/docker-compose.yml
PROJECT_NAME=pokedex

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
