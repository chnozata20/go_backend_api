package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware handles structured request logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()
		contentLength := c.Writer.Size()

		// Format path
		if raw != "" {
			path = path + "?" + raw
		}

		// Create structured log event
		event := log.Info()
		if statusCode >= 400 {
			event = log.Warn()
		}
		if statusCode >= 500 {
			event = log.Error()
		}

		// Add fields to log
		event.
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("ip", clientIP).
			Str("user_agent", userAgent).
			Int("content_length", contentLength).
			Str("timestamp", time.Now().Format(time.RFC3339))

		// Add error message if any
		if len(c.Errors) > 0 {
			event.Str("errors", c.Errors.String())
		}

		// Add request ID if available
		if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
			event.Str("request_id", requestID)
		}

		// Log the event
		event.Msg("HTTP Request")
	}
} 