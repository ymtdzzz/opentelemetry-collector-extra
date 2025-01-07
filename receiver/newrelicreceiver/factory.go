package newrelicreceiver

import (
	"context"

	"github.com/ymtdzzz/opentelemetry-collector-extra/receiver/newrelicreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, metadata.TracesStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesReceiver(_ context.Context, params receiver.Settings, cfg component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	rcfg := cfg.(*Config)
	if err := rcfg.Validate(); err != nil {
		return nil, err
	}
	r, err := newNewRelicReceiver(rcfg, params)
	if err != nil {
		return nil, err
	}

	r.(*newrelicReceiver).nextTracesConsumer = consumer
	return r, nil
}
