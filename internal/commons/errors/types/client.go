package types

import "net/http"

// ClientError represents a 400 Bad Request domain error.
type ClientError struct {
	baseError
}

var _ DomainError = (*ClientError)(nil)

// NewClientError creates a ClientError with the given error code.
func NewClientError(code string) *ClientError {
	return &ClientError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusBadRequest,
		},
	}
}

// NewClientErrorWithMessage creates a ClientError with a code and a human-readable message.
func NewClientErrorWithMessage(code, message string) *ClientError {
	return &ClientError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusBadRequest,
		},
	}
}
