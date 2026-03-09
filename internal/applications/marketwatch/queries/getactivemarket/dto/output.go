package dto

// Output holds the active market data needed to place orders.
type Output struct {
	MarketID    string
	Asset       string
	ConditionID string
	UpTokenID   string
	DownTokenID string
	TickSize    string
	WindowStart string // RFC3339
	WindowEnd   string // RFC3339
}
