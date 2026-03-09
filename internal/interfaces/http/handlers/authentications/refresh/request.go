// Package refresh defines request/response DTOs for the refresh action.
package refresh

// Request is the HTTP request body for token refresh.
type Request struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
