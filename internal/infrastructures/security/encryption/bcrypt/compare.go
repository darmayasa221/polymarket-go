package bcrypt

import (
	"context"

	goBcrypt "golang.org/x/crypto/bcrypt"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compare verifies a plain text password against a bcrypt hash.
// Returns nil if they match, AuthenticationError if not.
func (e *Encryption) Compare(_ context.Context, hashedPassword, plainPassword string) error {
	err := goBcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return errtypes.NewAuthenticationError(security.ErrPasswordMismatch)
	}
	return nil
}
