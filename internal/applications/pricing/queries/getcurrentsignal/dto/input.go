package dto

// Input holds the parameters for the GetCurrentSignal query.
type Input struct {
	Asset string // "btc" | "eth" | "sol" | "xrp"
}
