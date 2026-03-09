package observability

// NoopMetrics implements observability.Metrics with no-op operations.
type NoopMetrics struct{}

func (NoopMetrics) RecordRequest(method, path string, status int, durationMs float64) {}
func (NoopMetrics) RecordDBQuery(operation, table string, durationMs float64)         {}
func (NoopMetrics) RecordCacheHit(keyPrefix string)                                   {}
func (NoopMetrics) RecordCacheMiss(keyPrefix string)                                  {}
