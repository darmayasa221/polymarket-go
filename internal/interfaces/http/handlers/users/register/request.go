// Package register defines request/response DTOs for the register action.
package register

// Request is the HTTP request body for user registration.
type Request struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}
