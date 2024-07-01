package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode                   Mode   `env:"MODE"`
	DispatchServiceAddress string `env:"DISPATCH_SERVICE_ADDRESS"`
	DispatchServicePort    string `env:"DISPATCH_SERVICE_PORT"`
}

// Parses environment variables.
func Env() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}