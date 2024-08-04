package lintprocessor

import "go.opentelemetry.io/collector/component"

type Config struct {
	// TODO: configs such as ignore rules, versions, enable etc.
}

var _ component.Config = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
