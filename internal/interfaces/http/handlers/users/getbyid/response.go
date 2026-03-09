package getbyid

import "github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"

// Response is the HTTP response body for get user by ID.
type Response struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	FullName  string            `json:"full_name"`
	CreatedAt response.JSONTime `json:"created_at"`
	UpdatedAt response.JSONTime `json:"updated_at"`
}
