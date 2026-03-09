package types

import "net/http"

// ServiceUnavailableError represents a 503 Service Unavailable domain error.
type ServiceUnavailableError struct {
	baseError
}

var _ DomainError = (*ServiceUnavailableError)(nil)

// NewServiceUnavailableError creates a ServiceUnavailableError with the given error code.
func NewServiceUnavailableError(code string) *ServiceUnavailableError {
	return &ServiceUnavailableError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusServiceUnavailable,
		},
	}
}
