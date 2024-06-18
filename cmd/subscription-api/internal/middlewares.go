package internal

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger zap.SugaredLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func onTimeoutRunOut(ctx *gin.Context) {
	ctx.Status(http.StatusRequestTimeout)
}

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(2*time.Second),
		timeout.WithResponse(onTimeoutRunOut),
	)
}
