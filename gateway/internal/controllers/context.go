package controllers

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config/logger"
)

type Context interface {
	Logger() logger.Logger
	Context() context.Context
	Status(status int)
	String(status int, data string)
	BindJSON(data any) error
}
