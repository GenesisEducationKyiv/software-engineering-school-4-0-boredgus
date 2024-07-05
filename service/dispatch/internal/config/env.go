package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode Mode `env:"MODE"`

	DispatchServiceAddress string `env:"DISPATCH_SERVICE_ADDRESS"`
	DispatchServicePort    string `env:"DISPATCH_SERVICE_PORT"`

	PostgreSQLConnString string `env:"POSTGRESQL_CONN_STRING"`

	BrokerURL                 string `env:"NATS_URL"`
	CreateSubscriptionSubject string `env:"CREATE_SUBSCRIPTION_SUBJECT"`
}

func Env() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
