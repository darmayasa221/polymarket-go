package repository

// Format: AUTH_REPO.ERROR_CODE.
const (
	ErrTokenNotFound     = "AUTH_REPO.TOKEN_NOT_FOUND"
	ErrTokenExpired      = "AUTH_REPO.TOKEN_EXPIRED"
	ErrTokenAddFailed    = "AUTH_REPO.ADD_FAILED"
	ErrTokenValueTaken   = "AUTH_REPO.TOKEN_VALUE_TAKEN"
	ErrTokenCheckFailed  = "AUTH_REPO.CHECK_FAILED"
	ErrTokenDeleteFailed = "AUTH_REPO.DELETE_FAILED"
)
