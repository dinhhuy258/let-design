.PHONY: lint
lint:
	@echo "==> Running lint check..."
	@golangci-lint run

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	migrate create -ext sql -dir migration $${name// /_}

.PHONY: up
up:
	@go run main.go service
