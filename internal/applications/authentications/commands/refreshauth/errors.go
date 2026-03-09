// Package refreshauth implements the refresh authentication command use case.
package refreshauth

// Format: REFRESH_AUTH.ERROR_CODE.
const (
	ErrTokenRequired        = "REFRESH_AUTH.TOKEN_REQUIRED"
	ErrTokenInvalid         = "REFRESH_AUTH.TOKEN_INVALID"
	ErrDeleteOldTokenFailed = "REFRESH_AUTH.DELETE_OLD_TOKEN_FAILED"
	ErrTokenEntityFailed    = "REFRESH_AUTH.TOKEN_ENTITY_FAILED"
	ErrPersistTokenFailed   = "REFRESH_AUTH.PERSIST_TOKEN_FAILED"
)
