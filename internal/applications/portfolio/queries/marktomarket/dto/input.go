package dto

// Input provides current market prices per outcome token ID.
type Input struct {
	Prices map[string]string // tokenID → current price (decimal string, 0–1 range)
}
