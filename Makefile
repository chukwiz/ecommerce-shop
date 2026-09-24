.PHONY: build run dev lint migrate-up migrate-down help docs-generate

help:
	@echo "Makefile commands:"
	@echo "  build          Build the application"
	@echo "  run            Run the application"
	@echo "  dev            Run the application in development mode"
	@echo "  lint           Run golangci-lint on the project"
	@echo "  format         Format the codebase and rearrange imports"
	@echo "  migrate-up     Apply database migrations"
	@echo "  migrate-down   Rollback database migrations"

build:
	go build -o bin/app cmd/api/main.go

run:
	go run ./cmd/api

dev:
	go run ./cmd/api/main.go

lint: format
	golangci-lint run ./...

format:
	@gofmt -s -w .
	@goimports -w .

docs-generate:
	mkdir -p docs
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db

migrate-up:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" down

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down