package types

import "net/http"

// ConfigError represents a 500 Internal Server Error domain error for startup configuration failures.
type ConfigError struct {
	baseError
}

var _ DomainError = (*ConfigError)(nil)

// NewConfigError creates a ConfigError with the given error code.
func NewConfigError(code string) *ConfigError {
	return &ConfigError{
		baseError: baseError{
			code:       code,
			httpStatus: http.StatusInternalServerError,
		},
	}
}

// NewConfigErrorWithMessage creates a ConfigError with a code and a human-readable message.
func NewConfigErrorWithMessage(code, message string) *ConfigError {
	return &ConfigError{
		baseError: baseError{
			code:       code,
			message:    message,
			httpStatus: http.StatusInternalServerError,
		},
	}
}
