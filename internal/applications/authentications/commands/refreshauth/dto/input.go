// Package dto defines data transfer objects for the refreshauth command.
package dto

// Input holds the refresh token to exchange for a new token pair.
type Input struct {
	RefreshToken string
}
