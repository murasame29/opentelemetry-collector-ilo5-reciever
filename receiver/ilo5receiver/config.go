package ilo5receiver

import (
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/scraper/scraperhelper"

	"github.com/murasame29/opentelemetry-collector-ilo5-reciever/receiver/ilo5receiver/internal/metadata"
)

// Config defines the configuration for the receiver.
type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	metadata.MetricsBuilderConfig  `mapstructure:",squash"`

	Endpoint           string              `mapstructure:"endpoint"`
	Username           string              `mapstructure:"username"`
	Password           configopaque.String `mapstructure:"password"`
	InsecureSkipVerify bool                `mapstructure:"insecure_skip_verify"`
}

func (c *Config) Validate() error {
	return nil
}
