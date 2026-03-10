package dto

// Input holds the fields needed to open a new position.
type Input struct {
	Asset    string // "btc" | "eth" | "sol" | "xrp"
	TokenID  string
	Outcome  string // "Up" | "Down"
	Size     string // decimal string e.g. "15"
	AvgPrice string // decimal string e.g. "0.62"
	MarketID string
}
