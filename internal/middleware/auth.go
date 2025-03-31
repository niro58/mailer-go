package middleware

import (
	"mailer-go/internal/environment"
	handler "mailer-go/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth != environment.Environment.ApiAuth {
			handler.Respond(c, nil, handler.ErrUnauthorized)
			c.Abort()
		}

		c.Next()
	}
}
