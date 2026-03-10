package dto

// Output reports the result of closing the position.
type Output struct {
	PositionID  string
	RealisedPnL string // decimal string; positive = profit, negative = loss
}
