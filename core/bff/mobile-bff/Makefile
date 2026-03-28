.PHONY: help test unit-test integration-test coverage build run clean dev docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make test              - Run all tests (unit + integration)"
	@echo "  make unit-test         - Run unit tests only"
	@echo "  make integration-test  - Run integration tests"
	@echo "  make coverage          - Run tests with coverage report"
	@echo "  make build             - Build the application"
	@echo "  make run               - Run the application"
	@echo "  make docker-up         - Start PostgreSQL containers for testing"
	@echo "  make docker-down       - Stop PostgreSQL containers"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make dev               - Run in development mode with hot reload"

test: unit-test
	@echo "Running all tests passed!"

unit-test:
	@echo "Running unit tests..."
	@go test -v -race ./tests/unit

integration-test: docker-up
	@echo "Running integration tests..."
	@DATABASE_URL="postgres://postgres:postgres@localhost:5433/pokedex_test" go test -v -race ./tests/integration -timeout 30s
	@$(MAKE) docker-down

coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./tests/unit ./tests/mocks ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

build:
	@echo "Building mobile-bff..."
	@go build -o bin/mobile-bff ./cmd/server/main.go

run: build
	@echo "Running mobile-bff..."
	@MOBILE_BFF_PORT=8080 ./bin/mobile-bff

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

dev:
	@echo "Running in development mode..."
	@go run ./cmd/server/main.go

docker-up:
	@echo "Starting PostgreSQL containers..."
	@docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for database to be ready..."
	@sleep 10

docker-down:
	@echo "Stopping PostgreSQL containers..."
	@docker-compose -f docker-compose.test.yml down
	@go tool cover -func=coverage.out | grep total

build:
	@echo "Building mobile-bff..."
	@go build -o bin/mobile-bff ./cmd/server

run: build
	@echo "Running mobile-bff..."
	@MOBILE_BFF_PORT=8080 ./bin/mobile-bff

dev:
	@echo "Running in development mode..."
	@MOBILE_BFF_PORT=8080 go run ./cmd/server

clean:
	@echo "Cleaning up..."
	@rm -rf bin/ coverage.out coverage.html

lint:
	@echo "Running go fmt and go vet..."
	@go fmt ./...
	@go vet ./...
