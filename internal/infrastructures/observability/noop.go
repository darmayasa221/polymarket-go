// Package observability provides observability implementations.
// NoopTracker is used in tests and development without telemetry setup.
package observability

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/observability"
)

// Compile-time assertion: NoopTracker implements observability.Tracker.
var _ observability.Tracker = (*NoopTracker)(nil)

// NoopTracker is a no-operation Tracker — records nothing.
// Use in tests and when OpenTelemetry is not configured.
type NoopTracker struct {
	NoopMetrics
	NoopTracing
}

// NewNoopTracker creates a NoopTracker.
func NewNoopTracker() *NoopTracker { return &NoopTracker{} }
