.PHONY: help test unit-test integration-test coverage build run clean dev docker-up docker-down lint

help:
	@echo "Available commands:"
	@echo "  make test              - Run all tests (unit + integration)"
	@echo "  make unit-test         - Run unit tests only"
	@echo "  make integration-test  - Run integration tests (sobe Postgres de teste)"
	@echo "  make coverage          - Run tests with coverage report"
	@echo "  make build             - Build the application"
	@echo "  make run               - Run the application"
	@echo "  make lint              - Run go fmt and go vet"
	@echo "  make docker-up         - Start PostgreSQL containers for testing"
	@echo "  make docker-down       - Stop PostgreSQL containers"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make dev               - Run in development mode"

test: unit-test

unit-test:
	@echo "Running unit tests..."
	@go test -v -race ./internal/...

integration-test: docker-up
	@echo "Running integration tests..."
	@TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5433/pokedex_test?sslmode=disable" \
		go test -v -race -count=1 ./tests/integration -timeout 60s
	@$(MAKE) docker-down

coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | tail -1
	@echo "Coverage report generated: coverage.html"

build:
	@echo "Building mobile-bff..."
	@go build -o bin/mobile-bff ./cmd/server

run: build
	@echo "Running mobile-bff..."
	@MOBILE_BFF_PORT=8080 ./bin/mobile-bff

dev:
	@echo "Running in development mode..."
	@MOBILE_BFF_PORT=8080 go run ./cmd/server

lint:
	@echo "Running go fmt and go vet..."
	@out="$$(gofmt -l .)" && [ -z "$$out" ] || { echo "Arquivos fora do gofmt:"; echo "$$out"; exit 1; }
	@go vet ./...

docker-up:
	@echo "Starting PostgreSQL containers..."
	@docker compose -f docker-compose.test.yml up -d
	@echo "Waiting for database to be ready..."
	@sleep 10

docker-down:
	@echo "Stopping PostgreSQL containers..."
	@docker compose -f docker-compose.test.yml down

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
