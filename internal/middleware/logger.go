package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware handles request logging
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Stop timer
        end := time.Now()
        latency := end.Sub(start)

        // Get request details
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

        // Format log message
        if raw != "" {
            path = path + "?" + raw
        }

        // Log request details
        fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %#v\n%s",
            end.Format("2006/01/02 - 15:04:05"),
            statusCode,
            latency,
            clientIP,
            method,
            path,
            errorMessage,
        )
    }
} 