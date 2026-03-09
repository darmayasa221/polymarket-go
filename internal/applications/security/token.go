// Package security defines interfaces for security application services.
package security

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/security/types"
)

// TokenManager defines token creation and verification operations.
// Implemented in infrastructures/security/token/jwt/.
type TokenManager interface {
	// CreateTokenPair generates an access + refresh token pair for a user.
	CreateTokenPair(ctx context.Context, userID string) (types.TokenPair, error)

	// VerifyAccessToken validates an access token and returns its claims.
	VerifyAccessToken(ctx context.Context, tokenValue string) (*types.TokenClaims, error)

	// VerifyRefreshToken validates a refresh token and returns its claims.
	VerifyRefreshToken(ctx context.Context, tokenValue string) (*types.TokenClaims, error)
}
