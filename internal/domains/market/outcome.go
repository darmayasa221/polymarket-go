package market

// Outcome represents the binary resolution of a 5-minute market.
// Always "Up" or "Down" — never "Yes" or "No".
type Outcome string

const (
	// Up means close price >= open price.
	Up Outcome = "Up"
	// Down means close price < open price.
	Down Outcome = "Down"
)
