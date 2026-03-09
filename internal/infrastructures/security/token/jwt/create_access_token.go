package jwt

import (
	"context"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	sectypes "github.com/darmayasa221/polymarket-go/internal/applications/security/types"
	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
)

// CreateTokenPair generates an access and refresh token pair for the given userID.
func (m *Manager) CreateTokenPair(_ context.Context, userID string) (sectypes.TokenPair, error) {
	accessToken, accessExpiry, err := m.createToken(userID, tokenconst.TypeAccess, m.cfg.AccessTokenDuration)
	if err != nil {
		return sectypes.TokenPair{}, errtypes.NewInternalServerError(security.ErrTokenCreationFailed)
	}

	refreshToken, refreshExpiry, err := m.createToken(userID, tokenconst.TypeRefresh, m.cfg.RefreshTokenDuration)
	if err != nil {
		return sectypes.TokenPair{}, errtypes.NewInternalServerError(security.ErrTokenCreationFailed)
	}

	return sectypes.TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiry,
		RefreshTokenExpiresAt: refreshExpiry,
	}, nil
}

// createToken signs a JWT with the given userID, tokenType, and duration.
// It returns the signed token string, the expiry time, and any signing error.
func (m *Manager) createToken(userID, tokenType string, duration time.Duration) (string, time.Time, error) {
	now := timeutil.Now()
	expiresAt := now.Add(duration)

	claims := gojwt.MapClaims{
		"sub":  userID,
		"type": tokenType,
		"iss":  m.cfg.Issuer,
		"exp":  expiresAt.Unix(),
		"iat":  now.Unix(),
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(m.cfg.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}
