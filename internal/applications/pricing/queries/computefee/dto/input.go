package dto

// Input holds the token price for fee computation.
type Input struct {
	// TokenPrice is the current token price (0.01–0.99).
	// E.g. "0.50" for a 50/50 market.
	TokenPrice string
}
