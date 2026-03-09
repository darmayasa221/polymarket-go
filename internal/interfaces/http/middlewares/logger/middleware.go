// Package logger provides HTTP request/response logging middleware.
package logger

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/commons/telemetry"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// New creates an HTTP logging middleware that logs each request with duration.
// It generates a unique request ID for every request and propagates any
// client-provided X-Correlation-ID (for distributed tracing across services).
func New(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := telemetry.NewRequestID()
		ctx := telemetry.WithRequestID(c.Request.Context(), requestID)
		c.Header(httpconst.HeaderRequestID, requestID)

		// Propagate caller-provided correlation ID for distributed tracing.
		// If absent, no correlation ID is injected — request ID alone is sufficient.
		correlationID := c.GetHeader(httpconst.HeaderCorrelationID)
		if correlationID != "" {
			ctx = telemetry.WithCorrelationID(ctx, correlationID)
			c.Header(httpconst.HeaderCorrelationID, correlationID)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		elapsed := float64(time.Since(start).Milliseconds())
		if correlationID != "" {
			logger.Info("request",
				logging.FieldRequestID(requestID),
				logging.FieldCorrelationID(correlationID),
				logging.FieldMethod(c.Request.Method),
				logging.FieldPath(c.Request.URL.Path),
				logging.FieldStatus(c.Writer.Status()),
				logging.FieldDuration(elapsed),
			)
		} else {
			logger.Info("request",
				logging.FieldRequestID(requestID),
				logging.FieldMethod(c.Request.Method),
				logging.FieldPath(c.Request.URL.Path),
				logging.FieldStatus(c.Writer.Status()),
				logging.FieldDuration(elapsed),
			)
		}
	}
}
