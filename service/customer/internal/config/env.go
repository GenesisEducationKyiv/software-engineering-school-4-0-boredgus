package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode Mode `env:"MODE"`

	MicroserviceName string `env:"CUSTOMER_SERVICE_NAME"`
	MetricsPort      string `env:"CUSTOMER_SERVICE_METRICS_PORT"`
	MetricsRoute     string `env:"CUSTOMER_SERVICE_METRICS_ROUTE"`

	CustomerServiceAddress string `env:"CUSTOMER_SERVICE_ADDRESS"`
	CustomerServicePort    string `env:"CUSTOMER_SERVICE_PORT"`

	DatabaseSchema string `env:"CUSTOMERS_DB"`
	DatabaseURL    string `env:"CUSTOMMERS_DB_CONN_STRING"`

	MetricsGatewayURL string `env:"METRICS_GATEWAY_URL"`
}

func GetEnv() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
