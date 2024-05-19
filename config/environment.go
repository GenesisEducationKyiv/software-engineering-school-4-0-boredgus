package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Mode string

const (
	DevMode  Mode = "dev"
	TestMode Mode = "test"
	ProdMode Mode = "prod"
)

// type Environment struct {
// 	Mode                   Mode   `env:"mode"`
// 	ExchangeCurrencyAPIKey string `env:"EXCHANGE_CURRENCY_API_KEY"`

// 	PostgreSQLAddress  string `env:"POSTGRESQL_ADDRESS"`
// 	PostgreSQLPort     string `env:"POSTGRESQL_PORT"`
// 	PostgreSQLDatabase string `env:"POSTGRESQL_DATABASE"`
// 	PostgreSQLUsername string `env:"POSTGRESQL_USERNAME"`
// 	PostgreSQLPassword string `env:"POSTGRESQL_PASSWORD"`

// 	SMTPHost     string `env:"SMTP_HOST"`
// 	SMTPPort     int    `env:"SMTP_PORT"`
// 	SMTPUsername string `env:"SMTP_USERNAME"`
// 	SMTPPassword string `env:"SMTP_PASSWORD"`
// }

// var envVars Environment

func LoadEnvFile(envFiles []string) {
	for _, filename := range envFiles {
		_, err := os.Stat(filename)
		if errors.Is(err, os.ErrNotExist) {
			Log().Infof("%v file is not provided, skipping loading", filename)
			return
		}
	}
	if err := godotenv.Load(envFiles...); err != nil {
		Log().Fatalf("failed to load %v file: %v", strings.Join(envFiles, ","), err)
	}
}
func InitEnvVariables /* [T interface{}] */ ( /* target T, */ envFiles ...string) {
	LoadEnvFile(envFiles)
	// if err := env.Parse(&envVars); err != nil {
	// 	Logger().Error("failed to load env file", err)
	// 	envVars = Environment{}
	// }
}

// func Env() Environment {
// 	return envVars
// }
