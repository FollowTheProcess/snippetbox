# Will load a .env file in, where our DB credentials are kept

set dotenv-load := true

port := ":8000"

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
    go run ./cmd/web -dsn="$DB_USER:$DB_PASSWORD@/$DB_NAME?parseTime=true" -port={{ port }}
