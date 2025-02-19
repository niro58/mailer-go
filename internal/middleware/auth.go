package middleware

import (
	handler "mailer-go/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiAuth := os.Getenv("API_AUTH")
		auth := c.GetHeader("Authorization")

		if auth != apiAuth {
			c.JSON(401, handler.CreateReply(nil, handler.ErrUnauthorized))
			c.Abort()
		}

		c.Next()
	}
}
