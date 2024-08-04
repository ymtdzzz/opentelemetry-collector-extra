package lintprocessor

import (
	"context"

	"github.com/ymtdzzz/otel-lint/pkg/linter"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type metricLintProcessor struct{}

func newMetricLintProcessor() *metricLintProcessor {
	return &metricLintProcessor{}
}

func (p *metricLintProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	return linter.OtelLinter.RunMetric(md)
}
