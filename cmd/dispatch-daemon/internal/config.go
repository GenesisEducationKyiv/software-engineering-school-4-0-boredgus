package internal

import (
	"subscription-api/config"

	"github.com/caarlos0/env/v9"
)

type envVars struct {
	Mode            config.Mode `env:"MODE"`
	MailmanHost     string      `env:"MAILMAN_HOST"`
	MailmanPort     int         `env:"MAILMAN_PORT"`
	MailmanEmail    string      `env:"MAILMAN_EMAIL"`
	MailmanPassword string      `env:"MAILMAN_PASSWORD"`
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
