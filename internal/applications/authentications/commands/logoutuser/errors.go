// Package logoutuser implements the logout user command use case.
package logoutuser

// Format: LOGOUT_USER.ERROR_CODE.
const (
	ErrTokenRequired = "LOGOUT_USER.TOKEN_REQUIRED"
	ErrDeleteFailed  = "LOGOUT_USER.DELETE_FAILED"
)
