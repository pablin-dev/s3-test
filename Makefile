# Build the application
build:
	go build -o bin/app ./cmd/app

# Run all tests
test:
	go test ./...

# Run linters
lint:
	golangci-lint run --config .golangci.yml ./...

# Generate mocks
mock:
	mockery --all

# Run all verification steps
all: lint test build
