package controllers

import (
	"context"
	"subscription-api/config"
)

type Context interface {
	Logger() config.Logger
	Context() context.Context
	Status(status int)
	String(status int, data string)
	BindJSON(data any) error
}
