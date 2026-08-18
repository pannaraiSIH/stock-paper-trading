include .env
export

.PHONY: format generate migrate run test

format:
	@echo "Formatting code..."
	go fmt ./...

generate:
	@echo "Generating sql code..."
	sqlc generate 

migrate:
	@echo "Starting schema migration..."
	migrate -path migrations -database "$(DATABASE_URL)" up

run:
	@echo "Starting server with Air..."
	air

test:
	@echo "Running test..."
	go test -v ./...

