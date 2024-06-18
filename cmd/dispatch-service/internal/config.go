package internal

import (
	"fmt"
	"subscription-api/config"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode config.Mode `env:"MODE"`

	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT"`
	SMTPEmail    string `env:"SMTP_EMAIL"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`

	DispatchServiceAddress string `env:"DISPATCH_SERVICE_ADDRESS"`
	DispatchServicePort    string `env:"DISPATCH_SERVICE_PORT"`

	PostgreSQLConnString string `env:"POSTGRESQL_CONN_STRING"`

	CurrencyServiceAddress string `env:"CURRENCY_SERVICE_ADDRESS"`
	CurrencyServicePort    string `env:"CURRENCY_SERVICE_PORT"`
}

func Env() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
