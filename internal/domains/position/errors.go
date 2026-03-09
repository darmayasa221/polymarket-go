package position

const (
	// ErrMarketIDRequired is returned when MarketID is empty.
	ErrMarketIDRequired = "POSITION.MARKET_ID_REQUIRED"
	// ErrSizeInvalid is returned when Size is zero or negative.
	ErrSizeInvalid = "POSITION.SIZE_INVALID"
	// ErrAvgPriceInvalid is returned when AvgPrice is zero or negative.
	ErrAvgPriceInvalid = "POSITION.AVG_PRICE_INVALID"
	// ErrTokenIDRequired is returned when TokenID is empty.
	ErrTokenIDRequired = "POSITION.TOKEN_ID_REQUIRED"
)
