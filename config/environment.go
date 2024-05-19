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
func InitEnvVariables(envFiles ...string) {
	LoadEnvFile(envFiles)
}
