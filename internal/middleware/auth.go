package middleware

import (
	handler "mailer-go/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

var apiAuth = os.Getenv("API_AUTH")

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, handler.CreateReply(nil, handler.ErrUnauthorized))
			c.Abort()
			return
		}

		if auth != apiAuth {
			c.JSON(401, handler.CreateReply(nil, handler.ErrUnauthorized))
			c.Abort()
		}

		c.Next()
	}
}
