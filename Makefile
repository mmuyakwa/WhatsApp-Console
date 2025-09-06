.PHONY: build run clean deps

# Variablen
BINARY_NAME=whatsapp-console
BINARY_PATH=./bin/$(BINARY_NAME)

# Build der Anwendung
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) .
	@echo "Build complete: $(BINARY_PATH)"

# Anwendung ausführen
run:
	@go run .

# Abhängigkeiten installieren
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Projekt bereinigen
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f whatsapp.db
	@echo "Clean complete"

# Development Setup
dev-setup: deps
	@echo "Setting up development environment..."
	@go install github.com/cosmtrek/air@latest
	@echo "Development setup complete"

# Hot reload für Development
dev:
	@air

# Tests ausführen
test:
	@go test ./...

# Hilfe
help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  deps      - Install dependencies"
	@echo "  clean     - Clean build artifacts"
	@echo "  dev-setup - Setup development environment"
	@echo "  dev       - Run with hot reload"
	@echo "  test      - Run tests"
	@echo "  help      - Show this help"

# Cross-Platform Build
build-all:
	@echo "Building for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/whatsapp-console-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/whatsapp-console-linux-arm64 .
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/whatsapp-console-windows-amd64.exe .
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/whatsapp-console-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/whatsapp-console-darwin-arm64 .
	@echo "All builds complete in dist/"