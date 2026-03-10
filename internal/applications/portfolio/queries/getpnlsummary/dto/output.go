package dto

// Output summarizes realized PnL across all closed positions.
// TotalUnrealisedPnL is always "0" — use MarkToMarket for the live figure.
type Output struct {
	TotalRealisedPnL   string // sum of all closed positions' realized PnL
	TotalUnrealisedPnL string // always "0"; use MarkToMarket for live figure
	WinCount           int    // positions closed at a profit
	LossCount          int    // positions closed at a loss
	TotalCount         int    // total closed positions
}
