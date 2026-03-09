// Package types defines types for the security application services.
package types

import "time"

// TokenClaims holds the data embedded in a JWT token.
type TokenClaims struct {
	UserID  string
	Type    string
	Purpose string
}

// TokenPair holds an access/refresh token pair generated on login.
type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}
