// Package observability defines interfaces for telemetry tracking.
package observability

// Tracker is the unified observability interface combining metrics and tracing.
// Implemented in infrastructures/observability/otel/.
type Tracker interface {
	Metrics
	Tracer
}
