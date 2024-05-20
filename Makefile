# Simple Makefile for a Go project

build:
	@echo "Building..."
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test: clean build
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean: clear-cert
	@echo "Cleaning..."
	@rm -f main

# Generate Self-Signed Certificates (local development only)
gen-cert:
	@echo "Generating self-signed certificates..."
	@if ! command -v openssl > /dev/null; then \
	    echo "Error: OpenSSL is not installed. Please install OpenSSL to generate certificates."; \
	    exit 1; \
	fi
	@openssl req -x509 -newkey rsa:4096 -keyout .cert/server.key \
		-out .cert/server.crt -nodes -days 365 \
		-subj '/CN=localhost'
	@echo "Certificates generated successfully."

clear-cert:
	@echo "Removing generated certificates..."
	@rm -f .cert/*

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test gen-cert clear-cert clean
