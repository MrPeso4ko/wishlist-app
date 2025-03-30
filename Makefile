# Makefile for Wishlist App

.PHONY: build run migrate test docker clean

# Build the application
build:
	go build -o bin/main ./cmd/main.go

# Run the application
run: build
	./bin/main

# Run tests
test:
	go test -v -cover ./...

# Build and run with Docker
docker:
	docker-compose up --build

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download