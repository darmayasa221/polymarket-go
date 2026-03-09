package types

import "net/http"

// FieldViolation represents a single field-level validation failure.
type FieldViolation struct {
	// Field is the snake_case field name (e.g. "email", "full_name").
	Field string
	// Message is a human-readable description of the violation (e.g. "must be a valid email address").
	Message string
}

// ValidationError represents a 422 Unprocessable Entity domain error with per-field violations.
type ValidationError struct {
	baseError
	violations []FieldViolation
}

var _ DomainError = (*ValidationError)(nil)

// NewValidationError creates a ValidationError with the given error code and field violations.
// A nil violations slice is normalised to an empty slice.
func NewValidationError(code string, violations []FieldViolation) *ValidationError {
	if violations == nil {
		violations = []FieldViolation{}
	}

	return &ValidationError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusUnprocessableEntity,
		},
		violations: violations,
	}
}

// NewValidationErrorWithMessage creates a ValidationError with a code, a human-readable message,
// and field violations. A nil violations slice is normalised to an empty slice.
func NewValidationErrorWithMessage(code, message string, violations []FieldViolation) *ValidationError {
	if violations == nil {
		violations = []FieldViolation{}
	}

	return &ValidationError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusUnprocessableEntity,
		},
		violations: violations,
	}
}

// GetViolations returns the list of field-level violations associated with this error.
func (e *ValidationError) GetViolations() []FieldViolation {
	return e.violations
}
