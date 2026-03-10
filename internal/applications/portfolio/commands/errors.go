package commands

const (
	// ErrInvalidAsset is returned when the asset ticker is not supported.
	ErrInvalidAsset = "PORTFOLIO.INVALID_ASSET"
	// ErrInvalidOutcome is returned when the outcome is not "Up" or "Down".
	ErrInvalidOutcome = "PORTFOLIO.INVALID_OUTCOME"
	// ErrInvalidSize is returned when the position size is zero or negative.
	ErrInvalidSize = "PORTFOLIO.INVALID_SIZE"
	// ErrInvalidPrice is returned when a price is zero or negative.
	ErrInvalidPrice = "PORTFOLIO.INVALID_PRICE"
	// ErrPositionNotFound is returned when the position cannot be located.
	ErrPositionNotFound = "PORTFOLIO.POSITION_NOT_FOUND"
	// ErrSaveFailed is returned when persisting a position fails.
	ErrSaveFailed = "PORTFOLIO.SAVE_FAILED"
	// ErrCloseFailed is returned when closing a position fails.
	ErrCloseFailed = "PORTFOLIO.CLOSE_FAILED"
)
