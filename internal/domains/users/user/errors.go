package user

// Format: DOMAIN.ERROR_CODE.
const (
	ErrIDRequired       = "USER.ID_REQUIRED"
	ErrUsernameRequired = "USER.USERNAME_REQUIRED"
	ErrUsernameTooShort = "USER.USERNAME_TOO_SHORT"
	ErrUsernameTooLong  = "USER.USERNAME_TOO_LONG"
	ErrEmailRequired    = "USER.EMAIL_REQUIRED"
	ErrEmailInvalid     = "USER.EMAIL_INVALID"
	ErrEmailTooLong     = "USER.EMAIL_TOO_LONG"
	ErrPasswordRequired = "USER.PASSWORD_REQUIRED"
	ErrFullNameRequired = "USER.FULL_NAME_REQUIRED"
	ErrFullNameTooShort = "USER.FULL_NAME_TOO_SHORT"
	ErrFullNameTooLong  = "USER.FULL_NAME_TOO_LONG"
)
