package ilo5receiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"

	"github.com/murasame29/opentelemetry-collector-ilo5-reciever/internal/ilo"
)

type iloReceiver struct {
	config *Config
	client *ilo.Client
	params receiver.Settings
}

func newReceiver(cfg *Config, params receiver.Settings) *iloReceiver {
	return &iloReceiver{
		config: cfg,
		params: params,
	}
}

func (r *iloReceiver) Start(ctx context.Context, host component.Host) error {
	// Client initialization is handled in scraper.start
	// But we can pre-initialize here if needed for shared state
	return nil
}

func (r *iloReceiver) Shutdown(ctx context.Context) error {
	return nil
}
