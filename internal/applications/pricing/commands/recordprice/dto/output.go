package dto

import "time"

// Output holds the result of a successful RecordPrice command.
type Output struct {
	Asset      string
	Source     string
	Value      string
	RecordedAt time.Time
}
