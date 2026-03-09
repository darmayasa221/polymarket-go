package observability

// Metrics defines application metric recording operations.
type Metrics interface {
	// RecordRequest records an HTTP request with method, path, status, and duration.
	RecordRequest(method, path string, status int, durationMs float64)

	// RecordDBQuery records a database query with operation, table, and duration.
	RecordDBQuery(operation, table string, durationMs float64)

	// RecordCacheHit records a cache hit for the given key prefix.
	RecordCacheHit(keyPrefix string)

	// RecordCacheMiss records a cache miss for the given key prefix.
	RecordCacheMiss(keyPrefix string)
}
