package types

import "net/http"

// AuthenticationError represents a 401 Unauthorized domain error.
type AuthenticationError struct {
	baseError
}

var _ DomainError = (*AuthenticationError)(nil)

// NewAuthenticationError creates an AuthenticationError with the given error code.
func NewAuthenticationError(code string) *AuthenticationError {
	return &AuthenticationError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusUnauthorized,
		},
	}
}

// NewAuthenticationErrorWithMessage creates an AuthenticationError with a code and a human-readable message.
func NewAuthenticationErrorWithMessage(code, message string) *AuthenticationError {
	return &AuthenticationError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusUnauthorized,
		},
	}
}
