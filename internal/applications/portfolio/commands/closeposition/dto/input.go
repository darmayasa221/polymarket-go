package dto

// Input identifies the position to close and at what price.
type Input struct {
	PositionID string
	ExitPrice  string // decimal string e.g. "0.82"
}
