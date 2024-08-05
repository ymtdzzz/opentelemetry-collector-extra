package lintprocessor

import (
	"context"

	"github.com/ymtdzzz/otel-lint/pkg/linter"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type metricLintProcessor struct {
	enable bool
	l      *linter.Linter
}

func newMetricLintProcessor(cfg *Config) *metricLintProcessor {
	return &metricLintProcessor{
		enable: cfg.Enable,
		l:      linter.NewLinter(cfg.LinterOpts()...),
	}
}

func (p *metricLintProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	return p.l.RunMetric(md)
}
