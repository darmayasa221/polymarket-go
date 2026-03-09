// Package telemetry provides observability utilities for tracing and metrics.
package telemetry

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/commons/crypto"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// NewRequestID generates a new request ID.
func NewRequestID() string { return crypto.GenerateUUID() }

// WithRequestID stores a request ID in the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestIDFrom extracts the request ID from context.
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
