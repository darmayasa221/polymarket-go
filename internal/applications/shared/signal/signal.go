package signal

import (
	"time"

	"github.com/shopspring/decimal"
)

// Signal is the output of GetCurrentSignal query.
type Signal struct {
	Asset        string
	Predicted    string          // "Up" | "Down"
	Confidence   decimal.Decimal // 0.00–1.00; capped at 5% price move = full confidence
	OpenPrice    decimal.Decimal
	CurrentPrice decimal.Decimal
	Source       string // "chainlink" | "binance"
	RecordedAt   time.Time
}
