package logging

import (
	"context"

	"go.uber.org/zap"

	"github.com/darmayasa221/polymarket-go/internal/commons/telemetry"
)

// FromContext enriches a logger with context values (request ID, user ID, correlation ID).
func FromContext(ctx context.Context, l *Logger) *Logger {
	var fields []zap.Field
	if id := telemetry.RequestIDFrom(ctx); id != "" {
		fields = append(fields, FieldRequestID(id))
	}
	if id := telemetry.CorrelationIDFrom(ctx); id != "" {
		fields = append(fields, FieldCorrelationID(id))
	}
	if id := telemetry.UserIDFrom(ctx); id != "" {
		fields = append(fields, FieldUserID(id))
	}
	if len(fields) == 0 {
		return l
	}
	return l.With(fields...)
}
