// Package jwt provides the JWT token manager for the infrastructures layer.
package jwt

import "time"

// Config holds the configuration for the JWT token manager.
type Config struct {
	// SecretKey is the HMAC secret used to sign and verify tokens.
	SecretKey string
	// AccessTokenDuration is the lifetime of an access token.
	AccessTokenDuration time.Duration
	// RefreshTokenDuration is the lifetime of a refresh token.
	RefreshTokenDuration time.Duration
	// Issuer is the token issuer claim.
	Issuer string
}
