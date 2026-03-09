package dto

// Input holds the data for a tick size update triggered by a WS tick_size_change event.
type Input struct {
	ConditionID string // 0x hex identifier of the market
	TokenID     string // informational — identifies which side triggered the event
	OldTickSize string // decimal for validation/logging
	NewTickSize string // decimal to store (must be one of: 0.1, 0.01, 0.001, 0.0001)
}
