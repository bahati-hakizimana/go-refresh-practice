# ===============================
# Go Refresh Course - Makefile
# ===============================

# Build production binary
build:
	@echo "Building application..."
	@go build -o bin/go-refresh-course cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run application (Windows-safe)
run:
	@echo "Building application..."
	@go build -o bin/app.exe cmd/main.go
	@echo "Starting application..."
	@bin/app.exe

# Create new migration
migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

# Apply migrations
migrate-up:
	@echo "Applying migrations..."
	@go run cmd/migrate/main.go up

# Rollback migrations
migrate-down:
	@echo "Rolling back migrations..."
	@go run cmd/migrate/main.go down