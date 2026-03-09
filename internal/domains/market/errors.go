package market

const (
	// ErrIDRequired is returned when the Gamma event ID is empty.
	ErrIDRequired = "MARKET.ID_REQUIRED"
	// ErrInvalidAsset is returned when the asset is not BTC/ETH/SOL/XRP.
	ErrInvalidAsset = "MARKET.INVALID_ASSET"
	// ErrWindowStartRequired is returned when WindowStart is zero.
	ErrWindowStartRequired = "MARKET.WINDOW_START_REQUIRED"
	// ErrConditionIDRequired is returned when ConditionID is empty.
	ErrConditionIDRequired = "MARKET.CONDITION_ID_REQUIRED"
	// ErrUpTokenRequired is returned when UpTokenID is empty.
	ErrUpTokenRequired = "MARKET.UP_TOKEN_REQUIRED"
	// ErrDownTokenRequired is returned when DownTokenID is empty.
	ErrDownTokenRequired = "MARKET.DOWN_TOKEN_REQUIRED"
	// ErrTickSizeInvalid is returned when TickSize is zero or negative.
	ErrTickSizeInvalid = "MARKET.TICK_SIZE_INVALID"
)
