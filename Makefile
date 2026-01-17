.PHONY: build test test-coverage lint clean

RECEIVER_DIR := receiver/ilo5receiver

# Build
build:
	cd $(RECEIVER_DIR) && go build -v ./...

# Test
test:
	cd $(RECEIVER_DIR) && go test -v -race ./...

# Test with coverage
test-coverage:
	cd $(RECEIVER_DIR) && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Lint
lint:
	cd $(RECEIVER_DIR) && go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

# Clean
clean:
	rm -f $(RECEIVER_DIR)/coverage.out
