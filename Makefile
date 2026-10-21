# --- Variables ---
-include .env

MIGRATIONS_DIR=internal/db/migrate/migrations
MIGRATIONS_PATH=file://$(MIGRATIONS_DIR)
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable

# --- Commands ---
.PHONY: build dev run lint seed migrate-create migrate-up migrate-down db-shell docker-up docker-down reset-docker help

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build                 Build the binary"
	@echo "  dev                   Run the stack with live reload (Air, foreground)"
	@echo "  run                   Run the built binary"
	@echo "  lint                  Run golangci-lint"
	@echo "  seed                  Go run the seed script"
	@echo "  migrate-create name=x Create a new up/down migration pair"
	@echo "  migrate-up            Run all up migrations"
	@echo "  migrate-down          Rollback the last migration"
	@echo "  db-shell              Open psql in the postgres container"
	@echo "  docker-up             Start docker stack (detached)"
	@echo "  docker-down           Stop docker stack"
	@echo "  reset-docker          Wipe volumes and rebuild from zero"

build:
	go build -o bin/app ./cmd/api

dev:
	docker compose up --build

run:
	./bin/app

lint:
	golangci-lint run ./...

seed:
	DB_HOST=localhost go run ./cmd/api/seed

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=<name>"; exit 1; fi
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $(name)

migrate-up:
	migrate -source $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -source $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

db-shell:
	docker compose exec postgres psql -U $(DB_USER) -d $(DB_NAME)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

reset-docker:
	docker compose down -v
	docker compose up --build -d
