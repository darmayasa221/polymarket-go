// Package list defines the response DTO for the list users action.
package list

import "github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"

// UserItem is a single user in the list response.
type UserItem struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	FullName  string            `json:"full_name"`
	CreatedAt response.JSONTime `json:"created_at"`
}
