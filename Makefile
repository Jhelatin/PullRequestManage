COMPOSE_PROJECT_NAME=$(shell head -n 1 go.mod | sed 's/module //')

.PHONY: up down build logs

up:
	docker-compose up -d --build

down:
	docker-compose down -v

build:
	docker-compose build

logs:
	docker-compose logs -f app

.PHONY: migrate-up migrate-down

migrate-up:
	docker-compose run --rm migrate -path /migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable" up

migrate-down:
	docker-compose run --rm migrate -path /migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable" down

.PHONY: test generate-client run-e2e-tests

generate-client:
	@mkdir -p internal/client
	oapi-codegen -generate="types,client" -package=client -o ./internal/client/api_client.go openapi.yaml

run-e2e-tests:
	go test -v -tags=e2e ./...

test:
	@echo "--- Ensuring clean state ---"
	@make down
	@echo "--- Generating API client from openapi.yaml ---"
	@make generate-client
	@echo "--- Starting services in background ---"
	@make up
	@echo "--- Waiting for services to be ready (5 seconds)... ---"
	@sleep 5
	@echo "--- Running E2E tests ---"
	@make run-e2e-tests
	@echo "--- Tearing down services ---"
	@make down
	@echo "--- Test run complete! ---"

.PHONY: help

help:
	@echo "Available commands:"
	@echo "  up                - Build and start all services in the background"
	@echo "  down              - Stop and remove all services, networks, and volumes"
	@echo "  build             - Rebuild service images"
	@echo "  logs              - Follow the logs of the application service"
	@echo "  migrate-up        - Apply all available database migrations"
	@echo "  migrate-down      - Roll back the last database migration"
	@echo "  generate-client   - Generate Go client from openapi.yaml"
	@echo "  test              - Run the full E2E test suite (up -> test -> down)"
	@echo "  run-e2e-tests     - Run only the E2E tests (assumes services are running)"