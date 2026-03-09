package dto

import "time"

// Output holds the token pair returned after successful login.
type Output struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	UserID                string
}
