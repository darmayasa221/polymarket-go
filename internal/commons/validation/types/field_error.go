// Package types defines validation error types.
package types

import "fmt"

// FieldError represents a validation error on a specific field.
type FieldError struct {
	Field   string
	Code    string
	Message string
}

// Error implements the error interface.
func (e *FieldError) Error() string {
	return fmt.Sprintf("field %s: %s", e.Field, e.Code)
}
