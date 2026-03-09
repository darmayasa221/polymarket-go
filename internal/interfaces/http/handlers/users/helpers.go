package users

import (
	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// unauthorizedError returns an authentication error for missing user ID.
func unauthorizedError() error {
	return types.NewAuthenticationError(errkeys.ErrUnauthorized)
}
