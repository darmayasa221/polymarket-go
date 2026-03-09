// Package order defines the Order aggregate for Polymarket CLOB trading.
package order

// Side represents whether an order is a buy or sell.
type Side uint8

const (
	// Buy means purchasing outcome tokens.
	Buy Side = 0
	// Sell means selling outcome tokens.
	Sell Side = 1
)
