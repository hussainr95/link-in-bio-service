# Makefile for Link in Bio Service

.PHONY: all build test run swag-install swag docker-up docker-down clean

# Default target builds the binary
all: build

# Build the Go binary
build:
	go mod tidy
	go build -o main ./cmd/server

# Run unit tests with verbose output and coverage report
test:
	go test ./... -v -cover

# Run the application locally (without Docker)
run:
	go run ./cmd/server/main.go

# Install Swagger tools (swag)
swag-install:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation (requires swag installed)
swag:
	swag init -g cmd/server/main.go

# Build and run Docker containers using docker-compose
docker-up:
	docker-compose up --build

# Stop and remove Docker containers
docker-down:
	docker-compose down

# Clean up generated files and binary
clean:
	rm -f main
	rm -rf docs
