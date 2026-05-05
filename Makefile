# Variables
BINARY_NAME=api
GO_FILES=$(shell find . -name "*.go")
DB_URL=postgres://keeperadmin:securepass@localhost:5432/keeperdb?sslmode=disable
MIGRATE_PATH=migrations

.PHONY: build run test clean migrate-up migrate-down migrate-create

## Build: Compiles the project
build:
	go build -o bin/$(BINARY_NAME) cmd/api/main.go

## Run: Builds and executes the application
run: build
	./bin/$(BINARY_NAME)

## Test: Runs all tests with the race detector
test:
	go test -v -race ./...

## Clean: Removes build artifacts
clean:
	rm -rf bin/
	go clean

# --- Migrations ---

## migrate-create: Create a new migration file (usage: make migrate-create name=add_users_table)
migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

## migrate-up: Apply all up migrations
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" -verbose up

## migrate-down: Roll back the last migration
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" -verbose down 1

## migrate-force: Force migration to a specific version if it's in a dirty state
migrate-force:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(version)