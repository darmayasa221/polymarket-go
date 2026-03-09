package observability

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/observability"
)

// NoopSpan is a no-operation Span — all methods are no-ops.
type NoopSpan struct{}

func (NoopSpan) End()                             {}
func (NoopSpan) SetError(err error)               {}
func (NoopSpan) SetAttribute(key string, val any) {}

// NoopTracing implements observability.Tracer with no-op operations.
type NoopTracing struct{}

// StartSpan returns the same context and a NoopSpan.
func (NoopTracing) StartSpan(ctx context.Context, operationName string) (context.Context, observability.Span) {
	return ctx, NoopSpan{}
}
