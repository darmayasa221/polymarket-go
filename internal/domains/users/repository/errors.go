package repository

// Format: USER_REPO.ERROR_CODE.
const (
	ErrUserNotFound     = "USER_REPO.NOT_FOUND"
	ErrUsernameTaken    = "USER_REPO.USERNAME_TAKEN"
	ErrEmailTaken       = "USER_REPO.EMAIL_TAKEN"
	ErrUserAddFailed    = "USER_REPO.ADD_FAILED"
	ErrUserGetFailed    = "USER_REPO.GET_FAILED"
	ErrUserUpdateFailed = "USER_REPO.UPDATE_FAILED"
	ErrUserDeleteFailed = "USER_REPO.DELETE_FAILED"
)
