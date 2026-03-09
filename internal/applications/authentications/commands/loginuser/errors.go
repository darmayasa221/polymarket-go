// Package loginuser implements the login user command use case.
package loginuser

// Format: LOGIN_USER.ERROR_CODE.
const (
	ErrUsernameRequired   = "LOGIN_USER.USERNAME_REQUIRED"
	ErrPasswordRequired   = "LOGIN_USER.PASSWORD_REQUIRED"
	ErrInvalidCredentials = "LOGIN_USER.INVALID_CREDENTIALS"
	ErrGetIDFailed        = "LOGIN_USER.GET_ID_FAILED"
	ErrTokenEntityFailed  = "LOGIN_USER.TOKEN_ENTITY_FAILED"
	ErrPersistTokenFailed = "LOGIN_USER.PERSIST_TOKEN_FAILED"
)
