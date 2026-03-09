package providers

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	bcryptenc "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/encryption/bcrypt"
	jwttoken "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/token/jwt"
)

// ProvideEncryption creates the bcrypt encryption adapter.
func ProvideEncryption(cfg *config.Config) security.Encryption {
	return bcryptenc.New(cfg.Bcrypt)
}

// ProvideTokenManager creates the JWT token manager.
func ProvideTokenManager(cfg *config.Config) security.TokenManager {
	return jwttoken.New(cfg.JWT)
}
