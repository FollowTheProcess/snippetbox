set dotenv-load := true

_default:
    @just --list

# Tidy up dependencies
tidy:
    go mod tidy

# Run go fmt on all project files
fmt:
    gofumpt -extra -s -w .

# Start the development server
start:
    go run ./cmd/web
