package dto

import "time"

// Input holds the search parameters for GetActiveMarket.
type Input struct {
	Asset       string    // "btc" | "eth" | "sol" | "xrp"
	WindowStart time.Time // the 5-minute window boundary
}
