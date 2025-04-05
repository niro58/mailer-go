package middleware

import (
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
			Str("service", "everest-api").
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", raw).
			Int("status_code", statusCode).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Dur("latency_ms", latency)

		// Only add error fields if there's an actual error
		if errorMessage != "" {
			event.Str("error", errorMessage)
		}
		if statusCode >= 400 {
			event.Bool("is_error", true)
		}

		// Add stack trace for server errors
		if statusCode >= 500 {
			event.Str("stack_trace", string(debug.Stack()))
		}

		event.Msg("http_request")
	}
}
