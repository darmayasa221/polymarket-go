// Package login defines request/response DTOs for the login action.
package login

// Request is the HTTP request body for user login.
type Request struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
