package config

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config/logger"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/controllers"
	"github.com/gin-gonic/gin"
)

type ctx struct {
	c      *gin.Context
	ctx    context.Context
	logger logger.Logger
}

func NewContext(c *gin.Context, cx context.Context, logger logger.Logger) *ctx {
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
func (c *ctx) Logger() logger.Logger {
	return c.logger
}

type APIParams struct {
	CurrencyService controllers.CurrencyService
	DispatchService controllers.DispatchService
	Logger          logger.Logger
}

func GetRouter(params *APIParams) *gin.Engine {
	r := gin.Default()

	r.Use(TimeoutMiddleware())

	newContext := func(ctx *gin.Context) controllers.Context {
		return NewContext(ctx, context.Background(), params.Logger)
	}

	subscriptionController := controllers.NewSubscriptionController(params.DispatchService)

	r.GET("/rate", func(ctx *gin.Context) {
		controllers.GetExchangeRate(newContext(ctx), params.CurrencyService)
	})

	r.POST("/subscribe", func(ctx *gin.Context) {
		subscriptionController.SubscribeForDailyDispatch(newContext(ctx))
	})

	r.POST("/unsubscribe", func(ctx *gin.Context) {
		subscriptionController.UnsubscribeFromDailyDispatch(newContext(ctx))
	})

	return r
}
