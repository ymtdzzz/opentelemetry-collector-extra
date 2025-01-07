package newrelicreceiver

import (
	"errors"

	"go.opentelemetry.io/collector/component"
)

type Config struct {
	ServerCert string `mapstructure:"server_cert"`
	ServerKey  string `mapstructure:"server_key"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if cfg.ServerCert == "" || cfg.ServerKey == "" {
		return errors.New("server_cert and server_key must be set")
	}
	return nil
}
