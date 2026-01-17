.PHONY: build test test-coverage clean generate

RECEIVER_DIR := receiver/ilo5receiver

# Generate metadata
generate:
	cd $(RECEIVER_DIR) && mdatagen metadata.yaml

# Build
build:
	cd $(RECEIVER_DIR) && go build -v ./...

# Test
test:
	cd $(RECEIVER_DIR) && go test -v -race ./...

# Test with coverage
test-coverage:
	cd $(RECEIVER_DIR) && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Clean
clean:
	rm -f $(RECEIVER_DIR)/coverage.out
