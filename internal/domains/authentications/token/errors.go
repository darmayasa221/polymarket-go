package token

const (
	// ErrIDRequired is returned when the token ID is missing.
	ErrIDRequired = "TOKEN.ID_REQUIRED"
	// ErrUserIDRequired is returned when the user ID is missing.
	ErrUserIDRequired = "TOKEN.USER_ID_REQUIRED"
	// ErrValueRequired is returned when the token value is empty.
	ErrValueRequired = "TOKEN.VALUE_REQUIRED"
	// ErrValueTooShort is returned when the token value is below minimum length.
	ErrValueTooShort = "TOKEN.VALUE_TOO_SHORT"
	// ErrTypeRequired is returned when the token type is missing.
	ErrTypeRequired = "TOKEN.TYPE_REQUIRED"
	// ErrPurposeRequired is returned when the token purpose is missing.
	ErrPurposeRequired = "TOKEN.PURPOSE_REQUIRED"
	// ErrExpiresAtRequired is returned when the expiration time is zero.
	ErrExpiresAtRequired = "TOKEN.EXPIRES_AT_REQUIRED"
	// ErrTokenExpired is returned when the token has already expired.
	ErrTokenExpired = "TOKEN.EXPIRED"
)
