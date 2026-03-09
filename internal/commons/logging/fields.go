package logging

import "go.uber.org/zap"

// Standard field constructors for consistent log field naming.

// FieldRequestID returns a zap field for request ID.
func FieldRequestID(id string) zap.Field { return zap.String("request_id", id) }

// FieldCorrelationID returns a zap field for correlation ID.
func FieldCorrelationID(id string) zap.Field { return zap.String("correlation_id", id) }

// FieldUserID returns a zap field for user ID.
func FieldUserID(id string) zap.Field { return zap.String("user_id", id) }

// FieldOperation returns a zap field for operation name.
func FieldOperation(op string) zap.Field { return zap.String("operation", op) }

// FieldError returns a zap field for an error.
func FieldError(err error) zap.Field { return zap.Error(err) }

// FieldDuration returns a zap field for duration in milliseconds.
func FieldDuration(ms float64) zap.Field { return zap.Float64("duration_ms", ms) }

// FieldMethod returns a zap field for HTTP method.
func FieldMethod(method string) zap.Field { return zap.String("method", method) }

// FieldPath returns a zap field for HTTP path.
func FieldPath(path string) zap.Field { return zap.String("path", path) }

// FieldStatus returns a zap field for HTTP status code.
func FieldStatus(status int) zap.Field { return zap.Int("status", status) }

// FieldLayer returns a zap field for architectural layer.
func FieldLayer(layer string) zap.Field { return zap.String("layer", layer) }
