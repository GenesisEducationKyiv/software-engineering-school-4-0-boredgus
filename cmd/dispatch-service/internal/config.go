package internal

import (
	"subscription-api/config"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode                   config.Mode `env:"MODE"`
	MailmanHost            string      `env:"MAILMAN_HOST"`
	MailmanPort            int         `env:"MAILMAN_PORT"`
	MailmanEmail           string      `env:"MAILMAN_EMAIL"`
	MailmanPassword        string      `env:"MAILMAN_PASSWORD"`
	DispatchServiceAddress string      `env:"DISPATCH_SERVICE_ADDRESS"`
	DispatchServicePort    string      `env:"DISPATCH_SERVICE_PORT"`
	PostgreSQLConnString   string      `env:"POSTGRESQL_CONN_STRING"`
}

func Env() *envVars {
	var vars *envVars
	if err := env.Parse(&vars); err != nil {
		config.Log().Errorf("failed to parse env variables: %w", err)
	}
	return vars
}
