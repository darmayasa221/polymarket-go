package dto

import "time"

// Output holds the user data returned by the query.
type Output struct {
	ID        string
	Username  string
	Email     string
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
