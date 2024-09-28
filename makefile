.SILENT:
include .env

run:
	@go run ./cmd/api

build:
	@go build -ldflags='-s' -o=./bin/api ./cmd/api

force:
	@migrate -path=./internal/pkg/db/migrations -database="$(DB_URL)" force 1

swagit:
	@swag init -g ./cmd/api/main.go -o ./cmd/api/docs

dbuild:
	@docker build --no-cache .

up: dbuild
	@docker compose up -d

down:
	@docker compose down
