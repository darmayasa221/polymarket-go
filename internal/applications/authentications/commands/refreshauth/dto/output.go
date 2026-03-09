package dto

import "time"

// Output holds the new token pair returned after a successful refresh.
type Output struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}
