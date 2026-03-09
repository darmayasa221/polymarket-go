package telemetry

import "context"

const correlationIDKey contextKey = "correlation_id"

// WithCorrelationID stores a correlation ID in the context.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

// CorrelationIDFrom extracts the correlation ID from context.
func CorrelationIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(correlationIDKey).(string)
	return id
}
