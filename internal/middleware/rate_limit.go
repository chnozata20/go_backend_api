package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds rate limiting configuration
type RateLimiter struct {
    limiter *rate.Limiter
    mu      sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(rps), burst),
    }
}

// RateLimitMiddleware handles rate limiting
func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
    limiter := NewRateLimiter(rps, burst)
    
    return func(c *gin.Context) {
        limiter.mu.Lock()
        defer limiter.mu.Unlock()

        if !limiter.limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": time.Now().Add(time.Second).Format(time.RFC3339),
            })
            c.Abort()
            return
        }

        c.Next()
    }
} 