package internal

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger zap.SugaredLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
