.PHONY: build run clean install help

# Binary name
BINARY_NAME=peekfetch

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/peekfetch
	@echo "✓ Build complete: $(BINARY_NAME)"

# Build optimized binary
build-release:
	@echo "Building optimized $(BINARY_NAME)..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/peekfetch
	@echo "✓ Release build complete: $(BINARY_NAME)"

# Run the application
run: build
	@./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "✓ Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "✓ Dependencies installed"

# Install binary system-wide
install: build-release
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed successfully"

# Uninstall binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled successfully"

# Run tests (placeholder for future)
test:
	@echo "Running tests..."
	@go test -v ./...

# Display help
help:
	@echo "PeekFetch - Makefile commands:"
	@echo ""
	@echo "  make build          - Build the binary"
	@echo "  make build-release  - Build optimized binary"
	@echo "  make run            - Build and run the application"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Install dependencies"
	@echo "  make install        - Install binary system-wide"
	@echo "  make uninstall      - Remove binary from system"
	@echo "  make test           - Run tests"
	@echo "  make help           - Show this help message"
	@echo ""
