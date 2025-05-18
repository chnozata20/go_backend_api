package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics holds performance metrics
type Metrics struct {
    mu sync.RWMutex
    requestCount    int64
    totalLatency    time.Duration
    errorCount      int64
    endpointMetrics map[string]*EndpointMetrics
}

// EndpointMetrics holds metrics for a specific endpoint
type EndpointMetrics struct {
    RequestCount int64
    TotalLatency time.Duration
    ErrorCount   int64
    MinLatency   time.Duration
    MaxLatency   time.Duration
}

var metrics = &Metrics{
    endpointMetrics: make(map[string]*EndpointMetrics),
}

// MonitorMiddleware handles performance monitoring
func MonitorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        // Process request
        c.Next()

        // Calculate metrics
        latency := time.Since(start)
        status := c.Writer.Status()

        metrics.mu.Lock()
        defer metrics.mu.Unlock()

        // Update global metrics
        metrics.requestCount++
        metrics.totalLatency += latency
        if status >= 400 {
            metrics.errorCount++
        }

        // Update endpoint metrics
        if _, exists := metrics.endpointMetrics[path]; !exists {
            metrics.endpointMetrics[path] = &EndpointMetrics{
                MinLatency: latency,
                MaxLatency: latency,
            }
        }
        endpoint := metrics.endpointMetrics[path]
        endpoint.RequestCount++
        endpoint.TotalLatency += latency
        if status >= 400 {
            endpoint.ErrorCount++
        }
        if latency < endpoint.MinLatency {
            endpoint.MinLatency = latency
        }
        if latency > endpoint.MaxLatency {
            endpoint.MaxLatency = latency
        }

        // Log performance metrics
        fmt.Printf("[PERF] %s | %v | %d | %v\n",
            path,
            latency,
            status,
            time.Now().Format("2006-01-02 15:04:05"),
        )
    }
}

// GetMetrics returns current performance metrics
func GetMetrics() map[string]interface{} {
    metrics.mu.RLock()
    defer metrics.mu.RUnlock()

    return map[string]interface{}{
        "total_requests": metrics.requestCount,
        "total_latency": metrics.totalLatency,
        "error_count": metrics.errorCount,
        "endpoints": metrics.endpointMetrics,
    }
} 