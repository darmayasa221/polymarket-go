package bcrypt

import (
	"context"

	goBcrypt "golang.org/x/crypto/bcrypt"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Hash hashes a plain text password using bcrypt.
func (e *Encryption) Hash(_ context.Context, password string) (string, error) {
	hash, err := goBcrypt.GenerateFromPassword([]byte(password), e.cfg.Cost)
	if err != nil {
		return "", errtypes.NewInternalServerError(security.ErrPasswordHashFailed)
	}
	return string(hash), nil
}
