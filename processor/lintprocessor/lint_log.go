package lintprocessor

import (
	"context"

	"github.com/ymtdzzz/otel-lint/pkg/linter"
	"go.opentelemetry.io/collector/pdata/plog"
)

type logLintProcessor struct {
	enable bool
	l      *linter.Linter
}

func newLogLintProcessor(cfg *Config) *logLintProcessor {
	return &logLintProcessor{
		enable: cfg.Enable,
		l:      linter.NewLinter(cfg.LinterOpts()...),
	}
}

func (p *logLintProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	return p.l.RunLog(ld)
}
