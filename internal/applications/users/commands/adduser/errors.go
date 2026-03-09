// Package adduser implements the add user command use case.
package adduser

// Format: ADD_USER.ERROR_CODE.
const (
	ErrUsernameRequired = "ADD_USER.USERNAME_REQUIRED"
	ErrEmailRequired    = "ADD_USER.EMAIL_REQUIRED"
	ErrPasswordRequired = "ADD_USER.PASSWORD_REQUIRED"
	ErrFullNameRequired = "ADD_USER.FULL_NAME_REQUIRED"
	ErrAddFailed        = "ADD_USER.ADD_FAILED"
)
