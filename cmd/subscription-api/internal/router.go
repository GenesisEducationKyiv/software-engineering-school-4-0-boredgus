package internal

import (
	"context"
	"subscription-api/config"
	"subscription-api/internal/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type ctx struct {
	c      *gin.Context
	ctx    context.Context
	logger config.Logger
}

func NewContext(c *gin.Context, cx context.Context, logger config.Logger) *ctx {
	return &ctx{c: c, ctx: c, logger: logger}
}
func (c *ctx) Status(status int) {
	c.c.Status(status)
}
func (c *ctx) String(status int, data string) {
	c.c.String(status, data)
}
func (c *ctx) BindJSON(data any) error {
	return c.c.BindJSON(data)
}
func (c *ctx) Context() context.Context {
	return c.ctx
}
func (c *ctx) Logger() config.Logger {
	return c.logger
}

type APIParams struct {
	CurrencyService controllers.CurrencyService
	DispatchService controllers.DispatchService
	Logger          config.Logger
}

var MaxRequestDuration = 2000 * time.Millisecond

func GetRouter(params APIParams) *gin.Engine {
	r := gin.Default()

	r.GET("/rate", func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), MaxRequestDuration)
		defer cancel()
		controllers.GetExchangeRate(NewContext(ctx, c, params.Logger), params.CurrencyService)
	})

	r.POST("/subscribe", func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), MaxRequestDuration)
		defer cancel()
		controllers.SubscribeForDailyDispatch(NewContext(ctx, c, params.Logger), params.DispatchService)
	})

	return r
}
