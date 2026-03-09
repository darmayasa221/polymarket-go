// Package dto defines data transfer objects for the logoutuser command.
package dto

// Input holds the token to invalidate on logout.
type Input struct {
	TokenValue string
}
