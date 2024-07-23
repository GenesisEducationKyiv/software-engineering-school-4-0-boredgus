package config

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	timeout "github.com/s-wijaya/gin-timeout"
)

const HandlerTimeout time.Duration = 2 * time.Second

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.TimeoutHandler(HandlerTimeout, http.StatusRequestTimeout, nil)
}
