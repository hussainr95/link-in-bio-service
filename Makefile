.PHONY: all test swag-install swag docker-up docker-down

# Default target runs tests
all: test

# Run unit tests with verbose output and coverage report
test:
	go test ./... -v -cover

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
