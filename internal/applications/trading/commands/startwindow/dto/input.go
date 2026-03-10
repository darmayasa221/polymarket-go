package dto

import "github.com/shopspring/decimal"

// Input holds the market data needed to initialize a new trading window.
type Input struct {
	Asset       string          // "btc" | "eth" | "sol" | "xrp"
	MarketID    string          // Gamma event ID
	ConditionID string          // 0x hex
	UpTokenID   string          // ERC1155 token ID for "Up" outcome
	DownTokenID string          // ERC1155 token ID for "Down" outcome
	TickSize    string          // e.g. "0.01"
	OpenPrice   decimal.Decimal // first Chainlink price at window boundary
}
