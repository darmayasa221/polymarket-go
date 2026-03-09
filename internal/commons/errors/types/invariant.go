package types

import "net/http"

// InvariantError represents a 422 Unprocessable Entity domain error for business rule violations.
type InvariantError struct {
	baseError
}

var _ DomainError = (*InvariantError)(nil)

// NewInvariantError creates an InvariantError with the given error code.
func NewInvariantError(code string) *InvariantError {
	return &InvariantError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusUnprocessableEntity,
		},
	}
}

// NewInvariantErrorWithMessage creates an InvariantError with a code and a human-readable message.
func NewInvariantErrorWithMessage(code, message string) *InvariantError {
	return &InvariantError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusUnprocessableEntity,
		},
	}
}
