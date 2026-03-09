package types

import "net/http"

// NotFoundError represents a 404 Not Found domain error.
type NotFoundError struct {
	baseError
}

var _ DomainError = (*NotFoundError)(nil)

// NewNotFoundError creates a NotFoundError with the given error code.
func NewNotFoundError(code string) *NotFoundError {
	return &NotFoundError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusNotFound,
		},
	}
}

// NewNotFoundErrorWithMessage creates a NotFoundError with a code and a human-readable message.
func NewNotFoundErrorWithMessage(code, message string) *NotFoundError {
	return &NotFoundError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusNotFound,
		},
	}
}
