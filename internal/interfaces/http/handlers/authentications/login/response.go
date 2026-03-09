package login

import "github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"

// Response is the HTTP response body after successful login.
type Response struct {
	AccessToken           string            `json:"access_token"`
	RefreshToken          string            `json:"refresh_token"`
	AccessTokenExpiresAt  response.JSONTime `json:"access_token_expires_at"`
	RefreshTokenExpiresAt response.JSONTime `json:"refresh_token_expires_at"`
	ID                    string            `json:"id"`
}
