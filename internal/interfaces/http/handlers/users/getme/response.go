// Package getme defines the response DTO for the get-me action.
package getme

import "github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"

// Response is the HTTP response body for the authenticated user's own profile.
type Response struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	FullName  string            `json:"full_name"`
	CreatedAt response.JSONTime `json:"created_at"`
	UpdatedAt response.JSONTime `json:"updated_at"`
}
