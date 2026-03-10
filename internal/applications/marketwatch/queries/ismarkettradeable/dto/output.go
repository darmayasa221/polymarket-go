package dto

// Output reports whether a market is safe to trade.
type Output struct {
	Tradeable bool
	Reason    string // empty when Tradeable is true
	TickSize  string
	Active    bool
}
