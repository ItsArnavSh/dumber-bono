set shell := ["zsh", "-c"]
set export

# Common variables
bin := "dumber-bono"
ldflags := ""

# ──────────────────────────────────────────────
# Development
# ──────────────────────────────────────────────

# Run the application
run: install-deps
    go run .

# Run with race detector
run-race:
    go run -race .

# Build the binary
build:
    go build -ldflags "{{ldflags}}" -o bin/{{bin}} .

# Install dependencies
install-deps:
    go mod download

# Tidy go.mod / go.sum
tidy:
    go mod tidy

# Format all Go source files
fmt:
    gofmt -w .
    goimports -w .

# Format-check all Go source files (fails if not formatted)
fmt-check:
    gofmt -l .
    @if [ -n "$(gofmt -l .)" ]; then echo "error: files not formatted" && exit 1; fi

# Run go vet
vet:
    go vet ./...

# ──────────────────────────────────────────────
# Testing
# ──────────────────────────────────────────────

# Run all tests
test:
    go test ./...

# Run all tests verbosely
test-verbose:
    go test -v ./...

# Run all tests with race detector
test-race:
    go test -race ./...

# Run tests and show coverage report
cover:
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

# Run tests and open HTML coverage report
cover-html: cover
    go tool cover -html=coverage.out

# ──────────────────────────────────────────────
# Linting
# ──────────────────────────────────────────────

# Run the linter
lint:
    golangci-lint run ./...

# Run the linter and auto-fix what it can
lint-fix:
    golangci-lint run --fix ./...

# ──────────────────────────────────────────────
# CI (matches the GitHub Actions workflows)
# ──────────────────────────────────────────────

# Full local CI run: fmt-check -> vet -> lint -> test
ci: fmt-check vet lint test

# ──────────────────────────────────────────────
# Cleanup
# ──────────────────────────────────────────────

# Remove build artifacts and test outputs
clean:
    rm -rf bin coverage.out

# Show all available recipes
list:
    @just --list
