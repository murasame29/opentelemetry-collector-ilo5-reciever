.PHONY: build install-ocb generate run clean

# OCB version
OCB_VERSION := 0.143.0

# Install ocb (OpenTelemetry Collector Builder)
install-ocb:
	go install go.opentelemetry.io/collector/cmd/builder@v$(OCB_VERSION)

# Generate collector using ocb
generate: install-ocb
	builder --config=builder-config.yaml

# Build the collector
build: generate
	cd dist && go build -o ../ilo5-collector .

# Run with New Relic config
run:
	./ilo5-collector --config=config-newrelic.yaml

# Run with debug config (local testing)
run-debug:
	./ilo5-collector --config=config.yaml

# Clean build artifacts
clean:
	rm -rf dist ilo5-collector

# All-in-one: build and run
all: build run
