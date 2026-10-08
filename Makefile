# --- Variables ---
DB_URL=postgres://admin:password@localhost:5432/golpher_social?sslmode=disable
MIGRATIONS_PATH=file://internal/db/migrations
DB_USER=admin
DB_NAME=costaBackend

# --- Commands ---
.PHONY: build dev run lint migrate-up migrate-down db-shell docker-up docker-down help

help:
	echo "Usage: make [target]"
	echo "Targets:"
	echo "  build         Build the binary"
	echo "  dev           Run the stack with live reload (Air, foreground)"
	echo "  run           Run the built binary"
	echo "  lint          Run golangci-lint"
	echo "  migrate-up    Run all up migrations"
	echo "  migrate-down  Rollback the last migration"
	echo "  db-shell      Open psql in the postgres container"
	echo "  docker-up     Start docker stack (detached)"
	echo "  docker-down   Stop docker stack"

build:
	go build -o bin/app ./cmd/api

dev:
	docker compose up --build

run:
	./bin/app

lint:
	golangci-lint run ./...

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
