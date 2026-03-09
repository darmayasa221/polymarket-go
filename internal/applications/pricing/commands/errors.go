package commands

const (
	// ErrInvalidAsset is returned when the asset ticker is not supported.
	ErrInvalidAsset = "PRICING.INVALID_ASSET"
	// ErrInvalidSource is returned when the price source is not supported.
	ErrInvalidSource = "PRICING.INVALID_SOURCE"
	// ErrSaveFailed is returned when persisting a price observation fails.
	ErrSaveFailed = "PRICING.SAVE_FAILED"
)
