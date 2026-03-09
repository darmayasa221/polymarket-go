package user

import (
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/validation/rules"
)

// Validate checks all business invariants for the User.
// Returns an InvariantError if any rule is violated.
// Called automatically inside New().
func (u *User) Validate() error {
	if u.id.IsEmpty() {
		return types.NewInvariantError(ErrIDRequired)
	}
	if !rules.IsRequired(u.username) {
		return types.NewInvariantError(ErrUsernameRequired)
	}
	if !rules.IsMinLength(u.username, UsernameMinLength) {
		return types.NewInvariantError(ErrUsernameTooShort)
	}
	if !rules.IsMaxLength(u.username, UsernameMaxLength) {
		return types.NewInvariantError(ErrUsernameTooLong)
	}
	if !rules.IsRequired(u.email.String()) {
		return types.NewInvariantError(ErrEmailRequired)
	}
	if !rules.IsEmail(u.email.String()) {
		return types.NewInvariantError(ErrEmailInvalid)
	}
	if !rules.IsMaxLength(u.email.String(), EmailMaxLength) {
		return types.NewInvariantError(ErrEmailTooLong)
	}
	if !rules.IsRequired(u.hashedPassword.String()) {
		return types.NewInvariantError(ErrPasswordRequired)
	}
	if !rules.IsRequired(u.fullName) {
		return types.NewInvariantError(ErrFullNameRequired)
	}
	if !rules.IsMinLength(u.fullName, FullNameMinLength) {
		return types.NewInvariantError(ErrFullNameTooShort)
	}
	if !rules.IsMaxLength(u.fullName, FullNameMaxLength) {
		return types.NewInvariantError(ErrFullNameTooLong)
	}
	return nil
}
