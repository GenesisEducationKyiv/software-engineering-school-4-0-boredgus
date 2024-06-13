package internal

import (
	"subscription-api/config"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode                   config.Mode `env:"MODE"`
	Port                   string      `env:"CURRENCY_SERVICE_PORT"`
	CurrencyServiceServer  string      `env:"CURRENCY_SERVICE_SERVER"`
	ExchangeCurrencyAPIKey string      `env:"EXCHANGE_CURRENCY_API_KEY"`
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
