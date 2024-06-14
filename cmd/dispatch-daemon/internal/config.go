package internal

import (
	"fmt"
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

func Env() (*envVars, error) {
	var vars envVars
	if err := env.Parse(&vars); err != nil {
		return nil, fmt.Errorf("failed to parse env variables: %w", err)
	}
	return &vars, nil
}
