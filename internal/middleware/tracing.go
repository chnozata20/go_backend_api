package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
)

var tracer opentracing.Tracer

// InitTracer initializes Jaeger tracer
func InitTracer(serviceName string, jaegerHost string) error {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHost,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jaegermetrics.NullFactory

	t, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		return err
	}

	tracer = t
	opentracing.SetGlobalTracer(t)

	// Keep closer for cleanup
	_ = closer

	return nil
}

// TracingMiddleware adds distributed tracing to requests
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tracer == nil {
			c.Next()
			return
		}

		// Extract span context from headers
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			// Log error but continue
		}

		// Create span
		span := tracer.StartSpan(
			c.Request.URL.Path,
			ext.RPCServerOption(spanCtx),
			ext.SpanKindRPCServer,
		)
		defer span.Finish()

		// Add span to context
		c.Set("span", span)

		// Add span context to response headers
		tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Writer.Header()))

		// Add tags
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.Component.Set(span, "gin")

		// Process request
		c.Next()

		// Add response tags
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		if len(c.Errors) > 0 {
			ext.Error.Set(span, true)
			span.SetTag("error.message", c.Errors.String())
		}
	}
}

// GetSpanFromContext returns span from gin context
func GetSpanFromContext(c *gin.Context) opentracing.Span {
	if span, exists := c.Get("span"); exists {
		if s, ok := span.(opentracing.Span); ok {
			return s
		}
	}
	return nil
}

// StartChildSpan starts a child span
func StartChildSpan(c *gin.Context, operationName string) opentracing.Span {
	parentSpan := GetSpanFromContext(c)
	if parentSpan == nil {
		return nil
	}

	return tracer.StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
} 