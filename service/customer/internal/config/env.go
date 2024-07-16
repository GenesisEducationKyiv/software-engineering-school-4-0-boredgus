package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode Mode `env:"MODE"`

	CustomerServiceAddress string `env:"CUSTOMER_SERVICE_ADDRESS"`
	CustomerServicePort    string `env:"CUSTOMER_SERVICE_PORT"`

	DatabaseURL string `env:"SUBS_DB_CONN_STRING"`
}

func GetEnv() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
