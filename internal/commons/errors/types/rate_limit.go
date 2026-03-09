package types

import "net/http"

// RateLimitError represents a 429 Too Many Requests domain error.
type RateLimitError struct {
	baseError
}

var _ DomainError = (*RateLimitError)(nil)

// NewRateLimitError creates a RateLimitError with the given error code.
func NewRateLimitError(code string) *RateLimitError {
	return &RateLimitError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusTooManyRequests,
		},
	}
}
