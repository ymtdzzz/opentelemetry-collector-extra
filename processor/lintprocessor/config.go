package lintprocessor

import (
	"github.com/ymtdzzz/otel-lint/pkg/linter"
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	Enable             bool     `mapstructure:"enable" default:"true"`
	IgnoreExperimental bool     `mapstructure:"ignore_experimental" default:"false"`
	IgnoreWarn         bool     `mapstructure:"ignore_warn" default:"false"`
	IgnoreRules        []string `mapstructure:"ignore_rules"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}

func (cfg *Config) LinterOpts() []linter.Option {
	opts := []linter.Option{}
	if cfg.IgnoreExperimental {
		opts = append(opts, linter.IgnoreExperimental())
	}
	if cfg.IgnoreWarn {
		opts = append(opts, linter.IgnoreWarn())
	}
	if len(cfg.IgnoreRules) > 0 {
		opts = append(opts, linter.IgnoreRules(cfg.IgnoreRules))
	}
	return opts
}
