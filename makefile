.SILENT:
include .env

run:
	@go run ./cmd/api

build:
	@go build -ldflags='-s' -o=./bin/api ./cmd/api

up:
	@migrate -path=./internal/pkg/db/migrations -database="$(DB_URL)" up

down:
	@migrate -path=./internal/pkg/db/migrations -database="$(DB_URL)" down

force:
	@migrate -path=./internal/pkg/db/migrations -database="$(DB_URL)" force 1

swagit:
	@swag init -g ./cmd/api/main.go -o ./cmd/api/docs

start: dbuild
	@docker compose up -d

dbuild:
	@docker build --no-cache .
