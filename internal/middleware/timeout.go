package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const TimeoutTime = 10 * time.Second

func TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), TimeoutTime)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
