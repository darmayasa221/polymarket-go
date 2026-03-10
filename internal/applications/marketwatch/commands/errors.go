package commands

const (
	// ErrFetchFailed is returned when the Gamma API call fails.
	ErrFetchFailed = "MARKETWATCH.FETCH_FAILED"
	// ErrSaveFailed is returned when persisting a market fails.
	ErrSaveFailed = "MARKETWATCH.SAVE_FAILED"
	// ErrMarketNotFound is returned when no market matches the query.
	ErrMarketNotFound = "MARKETWATCH.MARKET_NOT_FOUND"
	// ErrInvalidTickSize is returned when the tick size is not a valid Polymarket value.
	ErrInvalidTickSize = "MARKETWATCH.INVALID_TICK_SIZE"
	// ErrUpdateFailed is returned when persisting a tick size update fails.
	ErrUpdateFailed = "MARKETWATCH.UPDATE_FAILED"
)
