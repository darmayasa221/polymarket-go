package dto

import "time"

// Input holds the data needed to record a new price observation.
type Input struct {
	Asset      string    // "btc" | "eth" | "sol" | "xrp"
	Source     string    // "chainlink" | "binance"
	Value      string    // decimal string e.g. "67234.50"
	RoundedAt  time.Time // Chainlink: when price was recorded; zero for Binance
	ReceivedAt time.Time // wall clock when received by the bot
}
