.PHONY: help dev build run migrate migrate-new seed test lint fmt clean docker-up docker-down

help: ## Mostrar este help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Rodar em modo desenvolvimento
	air

build: ## Build da aplicação
	go build -o bin/main ./cmd/server

run: ## Executar binário
	./bin/main

migrate: ## Aplicar migrations
	migrate -path migrations -database "$(DATABASE)" up

migrate-new: ## Criar nova migration
	@read -p "Nome da migration: " name; \
	migrate create -ext sql -dir migrations $$name

migrate-down: ## Reverter última migration
	migrate -path migrations -database "$(DATABASE)" down 1

seed: ## Popular banco com dados
	go run seed/seed.go

test: ## Rodar testes
	go test -v ./...

test-coverage: ## Rodar testes com coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Rodar linter
	golangci-lint run

fmt: ## Formatar código
	go fmt ./...
	goimports -w .

clean: ## Limpar arquivos temporários
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-up: ## Subir containers (banco + API)
	docker compose up -d

docker-down: ## Parar containers
	docker compose down

docker-logs: ## Ver logs dos containers
	docker compose logs -f

.DEFAULT_GOAL := help
