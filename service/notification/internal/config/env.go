package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode Mode `env:"MODE"`

	CurrencyServiceAddress string `env:"CURRENCY_SERVICE_ADDRESS"`
	CurrencyServicePort    string `env:"CURRENCY_SERVICE_PORT"`

	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT"`
	SMTPEmail    string `env:"SMTP_EMAIL"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`

	BrokerURL string `env:"NATS_URL"`
}

func GetEnv() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}

	return &vars, nil
}
