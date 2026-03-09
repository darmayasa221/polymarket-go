// Package keys defines machine-readable error code constants.
// Format: DOMAIN.ERROR_CODE
package keys

const (
	ErrInternalServer     = "GENERAL.INTERNAL_SERVER_ERROR"
	ErrInvalidRequestBody = "GENERAL.INVALID_REQUEST_BODY"
	ErrNotFound           = "GENERAL.NOT_FOUND"
	ErrUnauthorized       = "GENERAL.UNAUTHORIZED"
	ErrForbidden          = "GENERAL.FORBIDDEN"
	ErrValidationFailed   = "GENERAL.VALIDATION_FAILED"
	ErrTimeout            = "GENERAL.TIMEOUT"
)
