.PHONY: build test lint generate clean

# Build
build:
	go build -v ./...

# Test
test:
	go test -v -race ./...

# Test with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Lint
lint:
	golangci-lint run ./...

# Generate metadata (requires mdatagen)
generate:
	go install go.opentelemetry.io/collector/cmd/mdatagen@latest
	mdatagen metadata.yaml

# Clean
clean:
	rm -f coverage.out
