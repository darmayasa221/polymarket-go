package types

import "net/http"

// InternalServerError represents a 500 Internal Server Error domain error.
type InternalServerError struct {
	baseError
}

var _ DomainError = (*InternalServerError)(nil)

// NewInternalServerError creates an InternalServerError with the given error code.
func NewInternalServerError(code string) *InternalServerError {
	return &InternalServerError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusInternalServerError,
		},
	}
}
