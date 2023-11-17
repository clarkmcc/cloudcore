package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors.AllowAll().HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
