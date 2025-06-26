.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make sqlc         - Generate sqlc code"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make clean        - Clean build artifacts"

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o bin/ai-matching main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: migrate-up
migrate-up:
	migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

.PHONY: migrate-down
migrate-down:
	migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: docker-build
docker-build:
	docker build -t ai-matching:latest .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 --env-file .env ai-matching:latest

.PHONY: clean
clean:
	rm -rf bin/ tmp/

.PHONY: mod
mod:
	go mod download
	go mod tidy

.PHONY: air
air:
	air

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	go fmt ./...