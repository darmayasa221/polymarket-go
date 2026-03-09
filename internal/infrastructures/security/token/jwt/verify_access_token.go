package jwt

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	sectypes "github.com/darmayasa221/polymarket-go/internal/applications/security/types"
	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// VerifyAccessToken validates an access token and returns its claims.
func (m *Manager) VerifyAccessToken(_ context.Context, tokenValue string) (*sectypes.TokenClaims, error) {
	claims, err := m.parseAndVerify(tokenValue)
	if err != nil {
		return nil, err
	}
	if claims.Type != tokenconst.TypeAccess {
		return nil, errtypes.NewAuthenticationError(security.ErrTokenInvalid)
	}
	return claims, nil
}
