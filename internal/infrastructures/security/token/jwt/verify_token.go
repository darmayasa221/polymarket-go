package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	sectypes "github.com/darmayasa221/polymarket-go/internal/applications/security/types"
	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// parseAndVerify parses and validates a JWT token string, returning its claims.
func (m *Manager) parseAndVerify(tokenValue string) (*sectypes.TokenClaims, error) {
	token, err := gojwt.Parse(tokenValue, func(t *gojwt.Token) (any, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errtypes.NewAuthenticationError(security.ErrTokenInvalid)
		}
		return []byte(m.cfg.SecretKey), nil
	},
		gojwt.WithIssuer(m.cfg.Issuer),
		gojwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, errtypes.NewAuthenticationError(security.ErrTokenInvalid)
	}
	if !token.Valid {
		return nil, errtypes.NewAuthenticationError(security.ErrTokenInvalid)
	}

	mapClaims, ok := token.Claims.(gojwt.MapClaims)
	if !ok {
		return nil, errtypes.NewAuthenticationError(security.ErrTokenInvalid)
	}

	userID, _ := mapClaims["sub"].(string)
	tokenType, _ := mapClaims["type"].(string)

	return &sectypes.TokenClaims{
		UserID:  userID,
		Type:    tokenType,
		Purpose: tokenconst.PurposeAuthentication,
	}, nil
}
