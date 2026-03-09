package types

import "net/http"

// GatewayTimeoutError represents a 504 Gateway Timeout domain error.
type GatewayTimeoutError struct {
	baseError
}

var _ DomainError = (*GatewayTimeoutError)(nil)

// NewGatewayTimeoutError creates a GatewayTimeoutError with the given error code.
func NewGatewayTimeoutError(code string) *GatewayTimeoutError {
	return &GatewayTimeoutError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusGatewayTimeout,
		},
	}
}
