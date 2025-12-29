.PHONY: build test lint install clean help all dev

# Variables
BINARY_NAME=rusky
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X github.com/tito-sala/rusky/internal/cli.version=$(VERSION)"
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Default target
all: lint test build

# Development build (with debug info)
dev:
	@echo "Building development version..."
	go build -ldflags "-X github.com/tito-sala/rusky/internal/cli.version=$(VERSION)" -o $(BINARY_NAME) cmd/rusky/main.go
	@echo "Built: $(BINARY_NAME) (version: $(VERSION))"

# Production build (optimized)
build:
	@echo "Building production version..."
	go build $(LDFLAGS) -o $(BINARY_NAME) cmd/rusky/main.go
	@echo "Built: $(BINARY_NAME) (version: $(VERSION))"

# Build for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 cmd/rusky/main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 cmd/rusky/main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 cmd/rusky/main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 cmd/rusky/main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe cmd/rusky/main.go
	@echo "All builds complete in dist/"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Coverage report: coverage.out"

# Run tests with coverage report
coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage HTML report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run --timeout=5m

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w $(GO_FILES)
	goimports -w $(GO_FILES)

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./cmd/rusky

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@echo "Cleaned"

# Run the application
run: build
	./$(BINARY_NAME)

# Display help
help:
	@echo "Rusky Makefile Commands:"
	@echo ""
	@echo "  make build      - Build optimized binary"
	@echo "  make dev        - Build development binary (with debug info)"
	@echo "  make build-all  - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make test       - Run tests with race detection"
	@echo "  make coverage   - Generate HTML coverage report"
	@echo "  make lint       - Run golangci-lint"
	@echo "  make fmt        - Format code with gofmt and goimports"
	@echo "  make install    - Install binary to GOPATH/bin"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make run        - Build and run the application"
	@echo "  make all        - Run lint, test, and build (default)"
	@echo "  make help       - Display this help message"
	@echo ""
	@echo "Current version: $(VERSION)"
