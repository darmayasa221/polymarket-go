package oracle

const (
	// ErrAssetRequired is returned when asset is empty.
	ErrAssetRequired = "ORACLE.ASSET_REQUIRED"
	// ErrInvalidSource is returned when the price source is unknown.
	ErrInvalidSource = "ORACLE.INVALID_SOURCE"
	// ErrPriceValueInvalid is returned when the price is zero or negative.
	ErrPriceValueInvalid = "ORACLE.PRICE_VALUE_INVALID"
	// ErrReceivedAtRequired is returned when ReceivedAt is zero.
	ErrReceivedAtRequired = "ORACLE.RECEIVED_AT_REQUIRED"
)
