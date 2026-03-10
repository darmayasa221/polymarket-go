package dto

// PositionMark is the mark-to-market snapshot for a single open position.
type PositionMark struct {
	PositionID    string
	Asset         string
	Outcome       string
	Size          string
	AvgPrice      string
	CurrentPrice  string
	UnrealisedPnL string // positive = profit, negative = loss
}

// Output aggregates all position marks and the total unrealized PnL.
type Output struct {
	Marks              []PositionMark
	TotalUnrealisedPnL string
}
