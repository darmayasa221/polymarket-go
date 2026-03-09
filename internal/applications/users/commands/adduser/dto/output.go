package dto

import "time"

// Output carries the result of a successful AddUser operation.
type Output struct {
	ID        string
	Username  string
	Email     string
	FullName  string
	CreatedAt time.Time
}
