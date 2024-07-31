package config

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config/logger"
	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode                      logger.Mode `env:"MODE"`
	MicroserviceName          string      `env:"GATEWAY_NAME"`
	Port                      string      `env:"API_PORT"`
	CurrencyServiceAddress    string      `env:"CURRENCY_SERVICE_ADDRESS"`
	CurrencyServicePort       string      `env:"CURRENCY_SERVICE_PORT"`
	TransactionManagerAddress string      `env:"TRANSACTION_MANAGER_ADDRESS"`
	TransactionManagerPort    string      `env:"TRANSACTION_MANAGER_PORT"`
	MetricsGatewayURL         string      `env:"METRICS_GATEWAY_URL"`
}

func Env() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
