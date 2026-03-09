package types

import "net/http"

// AuthorizationError represents a 403 Forbidden domain error.
type AuthorizationError struct {
	baseError
}

var _ DomainError = (*AuthorizationError)(nil)

// NewAuthorizationError creates an AuthorizationError with the given error code.
func NewAuthorizationError(code string) *AuthorizationError {
	return &AuthorizationError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusForbidden,
		},
	}
}
