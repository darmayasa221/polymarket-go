// Package dto defines data transfer objects for the loginuser command.
package dto

// Input holds login credentials.
type Input struct {
	Username string
	Password string
}
