package order

const (
	// ErrMarketIDRequired is returned when MarketID is empty.
	ErrMarketIDRequired = "ORDER.MARKET_ID_REQUIRED"
	// ErrTokenIDRequired is returned when TokenID is empty.
	ErrTokenIDRequired = "ORDER.TOKEN_ID_REQUIRED"
	// ErrPriceInvalid is returned when Price is zero or negative.
	ErrPriceInvalid = "ORDER.PRICE_INVALID"
	// ErrSizeInvalid is returned when Size is zero or negative.
	ErrSizeInvalid = "ORDER.SIZE_INVALID"
	// ErrExpirationRequired is returned when GTD order has no expiration.
	ErrExpirationRequired = "ORDER.EXPIRATION_REQUIRED"
)
