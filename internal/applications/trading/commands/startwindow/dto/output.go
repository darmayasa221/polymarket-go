package dto

import "time"

// Output confirms the window was initialized with key identifiers.
type Output struct {
	Asset       string
	MarketID    string
	ConditionID string
	WindowStart time.Time
	WindowEnd   time.Time
}
