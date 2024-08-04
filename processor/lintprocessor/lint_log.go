package lintprocessor

import (
	"context"

	"github.com/ymtdzzz/otel-lint/pkg/linter"
	"go.opentelemetry.io/collector/pdata/plog"
)

type logLintProcessor struct{}

func newLogLintProcessor() *logLintProcessor {
	return &logLintProcessor{}
}

func (p *logLintProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	return linter.OtelLinter.RunLog(ld)
}
