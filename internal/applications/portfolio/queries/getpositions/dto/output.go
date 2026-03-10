package dto

import "time"

// PositionDTO is a read-only snapshot of an open position.
type PositionDTO struct {
	PositionID string
	Asset      string
	TokenID    string
	Outcome    string
	Size       string
	AvgPrice   string
	OpenedAt   time.Time
}

// Output holds all matching open positions.
type Output struct {
	Positions []PositionDTO
}
