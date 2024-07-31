package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode                   Mode   `env:"MODE"`
	Port                   string `env:"CURRENCY_SERVICE_PORT"`
	CurrencyServiceAddress string `env:"CURRENCY_SERVICE_ADDRESS"`
	CurrencyServicePort    string `env:"CURRENCY_SERVICE_PORT"`
	ExchangeCurrencyAPIKey string `env:"EXCHANGE_CURRENCY_API_KEY"`
	CurrencyBeaconAPIKey   string `env:"CURRENCY_BEACON_API_KEY"`

	MicroserviceName string `env:"CURRENCY_SERVICE_NAME"`
	MetricsPort      string `env:"CURRENCY_SERVICE_METRICS_PORT"`
	MetricsRoute     string `env:"CURRENCY_SERVICE_METRICS_ROUTE"`

	MetricsGatewayURL string `env:"METRICS_GATEWAY_URL"`
}

func GetEnv() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
