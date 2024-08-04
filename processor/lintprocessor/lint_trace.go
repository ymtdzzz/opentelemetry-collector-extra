package lintprocessor

import (
	"context"

	"github.com/ymtdzzz/otel-lint/pkg/linter"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type spanLintProcessor struct{}

func newSpanLintProcessor() *spanLintProcessor {
	return &spanLintProcessor{}
}

func (p *spanLintProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	return linter.OtelLinter.RunTrace(td)
}
