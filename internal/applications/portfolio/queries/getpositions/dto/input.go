package dto

// Input optionally filters positions by asset. Empty string returns all open positions.
type Input struct {
	Asset string // "btc" | "eth" | "sol" | "xrp" | "" for all
}
