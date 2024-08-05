package lintprocessor

import (
	"context"

	"github.com/ymtdzzz/otel-lint/pkg/linter"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type spanLintProcessor struct {
	enable bool
	l      *linter.Linter
}

func newSpanLintProcessor(cfg *Config) *spanLintProcessor {
	return &spanLintProcessor{
		enable: cfg.Enable,
		l:      linter.NewLinter(cfg.LinterOpts()...),
	}
}

func (p *spanLintProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	if p.enable {
		return p.l.RunTrace(td)
	}
	return td, nil
}
