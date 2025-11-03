build:
	@go build -o bin/go-refresh-course cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/go-refresh-course
