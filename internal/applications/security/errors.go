package security

// Format: DOMAIN.ERROR_CODE.
const (
	ErrTokenCreationFailed = "SECURITY.TOKEN_CREATION_FAILED"
	ErrTokenInvalid        = "SECURITY.TOKEN_INVALID"
	ErrTokenExpired        = "SECURITY.TOKEN_EXPIRED"
	ErrPasswordHashFailed  = "SECURITY.PASSWORD_HASH_FAILED"
	ErrPasswordMismatch    = "SECURITY.PASSWORD_MISMATCH"
)
