// Package dto defines data transfer objects for the adduser command.
package dto

// Input carries the data required to register a new user.
type Input struct {
	Username string
	Email    string
	Password string
	FullName string
}
