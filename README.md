# iLO5 Receiver for OpenTelemetry Collector

[![Go Reference](https://pkg.go.dev/badge/github.com/murasame29/opentelemetry-collector-ilo5-reciever.svg)](https://pkg.go.dev/github.com/murasame29/opentelemetry-collector-ilo5-reciever)
[![CI](https://github.com/murasame29/opentelemetry-collector-ilo5-reciever/actions/workflows/ci.yaml/badge.svg)](https://github.com/murasame29/opentelemetry-collector-ilo5-reciever/actions/workflows/ci.yaml)

HPE iLO5 からメトリクスを収集する OpenTelemetry Collector Receiver です。

## Metrics

| Metric | Description | Unit |
|--------|-------------|------|
| `ilo.system.power_state` | System power state (1=On, 0=Off) | 1 |
| `ilo.system.health` | System health (1=OK, 2=Warning, 3=Critical) | 1 |
| `ilo.power.consumption` | Power consumption | W |
| `ilo.power.capacity` | Power capacity | W |
| `ilo.power.voltage` | Voltage reading | V |
| `ilo.power.psu.output` | PSU output power | W |
| `ilo.power.psu.input_voltage` | PSU input voltage | V |
| `ilo.power.psu.health` | PSU health (1=OK, 2=Warning, 3=Critical) | 1 |
| `ilo.thermal.temperature` | Temperature reading | Cel |
| `ilo.fan.speed` | Fan speed | % |
| `ilo.storage.drive.health` | Drive health (1=OK, 2=Warning, 3=Critical) | 1 |

## Resource Attributes

| Attribute | Description |
|-----------|-------------|
| `host.name` | Server hostname |
| `host.serial_number` | Server serial number |
| `ilo.endpoint` | iLO endpoint URL |
| `ilo.model` | iLO model |

## Installation

```bash
go get github.com/murasame29/opentelemetry-collector-ilo5-reciever/receiver/ilo5receiver
```

## Configuration

```yaml
receivers:
  ilo5:
    endpoint: "https://ilo.example.com"
    username: "admin"
    password: "${ILO_PASSWORD}"
    collection_interval: 60s
    insecure_skip_verify: true
```

## Development

### Requirements

- Go 1.25+

### Build

```bash
make build
```

### Test

```bash
make test
```

### Test with Coverage

```bash
make test-coverage
```

### Generate Metadata

```bash
make generate
```

## Usage with OpenTelemetry Collector Builder

Add to your `builder-config.yaml`:

```yaml
receivers:
  - gomod: github.com/murasame29/opentelemetry-collector-ilo5-reciever/receiver/ilo5receiver v0.1.0
```

### Example Collector Config

```yaml
receivers:
  ilo5:
    endpoint: "https://192.168.1.100"
    username: "admin"
    password: "${ILO_PASSWORD}"
    collection_interval: 60s
    insecure_skip_verify: true

exporters:
  otlp:
    endpoint: "otlp.nr-data.net:4317"
    headers:
      api-key: "${NEW_RELIC_LICENSE_KEY}"

service:
  pipelines:
    metrics:
      receivers: [ilo5]
      exporters: [otlp]
```

## License

MIT
