.PHONY: up down build test lint migrate-up migrate-down sqlc-generate seed logs

up:
	docker compose up -d --build

down:
	docker compose down

build:
	go build -o api ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run ./...

migrate-up:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DATABASE_URL)" down

sqlc-generate:
	cd db && sqlc generate

seed:
	@echo "Set API_KEY_SEED env var and restart the API service"
	@echo "docker compose up -d api"

logs:
	docker compose logs -f
