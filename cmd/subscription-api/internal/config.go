package internal

import (
	"subscription-api/config"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode                   config.Mode `env:"MODE"`
	Port                   string      `env:"API_PORT"`
	CurrencyServiceAddress string      `env:"CURRENCY_SERVICE_ADDRESS"`
	CurrencyServicePort    string      `env:"CURRENCY_SERVICE_PORT"`
	DispatchServiceAddress string      `env:"DISPATCH_SERVICE_ADDRESS"`
	DispatchServicePort    string      `env:"DISPATCH_SERVICE_PORT"`
}

var vars envVars

func init() {
	if err := env.Parse(&vars); err != nil {
		config.Log().Errorf("failed to parse env variables: %w", err)
	}
}

func Env() envVars {
	return vars
}
