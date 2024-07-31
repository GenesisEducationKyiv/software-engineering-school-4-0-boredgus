package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode Mode `env:"MODE"`

	TransactionManagerAddress string `env:"TRANSACTION_MANAGER_ADDRESS"`
	TransactionManagerPort    string `env:"TRANSACTION_MANAGER_PORT"`

	DispatchServiceAddress string `env:"DISPATCH_SERVICE_ADDRESS"`
	DispatchServicePort    string `env:"DISPATCH_SERVICE_PORT"`

	CustomerServiceAddress string `env:"CUSTOMER_SERVICE_ADDRESS"`
	CustomerServicePort    string `env:"CUSTOMER_SERVICE_PORT"`

	BrokerURL string `env:"NATS_URL"`
}

func Env() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
