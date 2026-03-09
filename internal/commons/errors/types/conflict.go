package types

import "net/http"

// ConflictError represents a 409 Conflict domain error.
type ConflictError struct {
	baseError
}

var _ DomainError = (*ConflictError)(nil)

// NewConflictError creates a ConflictError with the given error code.
func NewConflictError(code string) *ConflictError {
	return &ConflictError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusConflict,
		},
	}
}
