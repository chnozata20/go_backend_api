package middleware

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)
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

// MonitorMiddleware handles performance monitoring with Prometheus metrics
func MonitorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        // Increment in-flight requests for Prometheus
        httpRequestsInFlight.Inc()
        defer httpRequestsInFlight.Dec()

        // Get request size for Prometheus
        requestSize := c.Request.ContentLength
        if requestSize > 0 {
            httpRequestSize.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(requestSize))
        }

        // Process request
        c.Next()

        // Calculate metrics
        latency := time.Since(start)
        status := c.Writer.Status()

        // Update custom metrics
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

        // Update Prometheus metrics
        statusStr := strconv.Itoa(status)
        endpointPath := c.FullPath()
        if endpointPath == "" {
            endpointPath = "unknown"
        }

        httpRequestsTotal.WithLabelValues(c.Request.Method, endpointPath, statusStr).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, endpointPath).Observe(latency.Seconds())

        // Get response size for Prometheus
        responseSize := c.Writer.Size()
        if responseSize > 0 {
            httpResponseSize.WithLabelValues(c.Request.Method, endpointPath).Observe(float64(responseSize))
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