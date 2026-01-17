package main

import (
	"log"

	"github.com/murasame29/ilo5-receiver/receiver/ilo5receiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/service/telemetry/otelconftelemetry"
)

func main() {
	info := component.BuildInfo{
		Command:     "ilo5-collector",
		Description: "Local collector for iLO5 receiver verification",
		Version:     "1.0.0",
	}

	receivers, err := otelcol.MakeFactoryMap(
		ilo5receiver.NewFactory(),
	)
	if err != nil {
		log.Fatalf("failed to make receiver factory map: %v", err)
	}

	debugFactory := debugexporter.NewFactory()
	exporters := map[component.Type]exporter.Factory{
		debugFactory.Type(): debugFactory,
	}

	factories := otelcol.Factories{
		Receivers: receivers,
		Exporters: exporters,
		Telemetry: otelconftelemetry.NewFactory(),
	}

	settings := otelcol.CollectorSettings{
		BuildInfo: info,
		Factories: func() (otelcol.Factories, error) {
			return factories, nil
		},
		ConfigProviderSettings: otelcol.ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				ProviderFactories: []confmap.ProviderFactory{
					fileprovider.NewFactory(),
					envprovider.NewFactory(),
				},
				DefaultScheme: "file",
			},
		},
	}

	cmd := otelcol.NewCommand(settings)
	if err := cmd.Execute(); err != nil {
		log.Fatalf("collector failed: %v", err)
	}
}
