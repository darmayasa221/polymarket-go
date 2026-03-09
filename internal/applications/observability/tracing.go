package observability

import "context"

// Span represents an active trace span.
type Span interface {
	// End completes the span.
	End()
	// SetError marks the span as failed with the given error.
	SetError(err error)
	// SetAttribute sets a key-value attribute on the span.
	SetAttribute(key string, value any)
}

// Tracer defines distributed tracing operations.
type Tracer interface {
	// StartSpan begins a new span for the given operation name.
	StartSpan(ctx context.Context, operationName string) (context.Context, Span)
}
