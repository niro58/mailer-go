package middleware

import (
	env "mailer-go/internal/environment"
	handler "mailer-go/internal/handlers"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		auth = strings.ReplaceAll(auth, "Bearer ", "")
		if auth != env.Env.ApiAuth {
			handler.Respond(c, nil, handler.ErrUnauthorized)
			c.Abort()
		}

		c.Next()
	}
}
