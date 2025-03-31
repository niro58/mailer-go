package middleware

import (
	"mailer-go/internal/environment"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate logging fields
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Build the log event
		event := log.Info() // Default level
		if statusCode >= 400 {
			event = log.Warn() // 4xx are warnings
		}
		if statusCode >= 500 {
			event = log.Error() // 5xx are real errors
		}

		// Standard fields
		event.
			Str("service", "everest-api").                 // Critical for Loki filtering
			Str("environment", environment.Environment.Mode). // dev/stage/prod
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", raw).
			Int("status_code", statusCode).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Dur("latency_ms", latency). // Duration in milliseconds
			Str("error", errorMessage).
			Bool("is_error", statusCode >= 400) // Boolean for easy filtering

		// Add stack trace for server errors
		if statusCode >= 500 {
			event.Str("stack_trace", string(debug.Stack()))
		}

		event.Msg("http_request")
	}
}
