package token

import (
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/validation/rules"
)

// Validate checks all business invariants for the Token.
// Returns an InvariantError if any rule is violated.
// Called automatically inside New().
func (t *Token) Validate() error {
	if t.id.IsEmpty() {
		return types.NewInvariantError(ErrIDRequired)
	}
	if t.userID.IsEmpty() {
		return types.NewInvariantError(ErrUserIDRequired)
	}
	if t.value.IsEmpty() {
		return types.NewInvariantError(ErrValueRequired)
	}
	if !rules.IsMinLength(t.value.String(), TokenValueMinLength) {
		return types.NewInvariantError(ErrValueTooShort)
	}
	if !rules.IsRequired(t.tokenType) {
		return types.NewInvariantError(ErrTypeRequired)
	}
	if !rules.IsRequired(t.purpose) {
		return types.NewInvariantError(ErrPurposeRequired)
	}
	if t.expiresAt.IsZero() {
		return types.NewInvariantError(ErrExpiresAtRequired)
	}
	return nil
}
