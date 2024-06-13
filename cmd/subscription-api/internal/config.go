package internal

import (
	"subscription-api/config"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode                  config.Mode `env:"MODE"`
	Port                  string      `env:"API_PORT"`
	CurrencyServiceServer string      `env:"CURRENCY_SERVICE_SERVER"`
	DispatchServiceServer string      `env:"DISPATCH_SERVICE_SERVER"`
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
