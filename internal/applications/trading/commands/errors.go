package commands

const (
	// ErrInvalidAsset is returned when the asset ticker is not supported.
	ErrInvalidAsset = "TRADING.INVALID_ASSET"
	// ErrWindowNotOpen is returned when an operation requires an open window.
	ErrWindowNotOpen = "TRADING.WINDOW_NOT_OPEN"
	// ErrOrderNotFound is returned when an order cannot be found.
	ErrOrderNotFound = "TRADING.ORDER_NOT_FOUND"
	// ErrSaveFailed is returned when persisting an order fails.
	ErrSaveFailed = "TRADING.SAVE_FAILED"
	// ErrSubmitFailed is returned when the CLOB rejects the order.
	ErrSubmitFailed = "TRADING.SUBMIT_FAILED"
	// ErrCancelFailed is returned when the CLOB rejects a cancel request.
	ErrCancelFailed = "TRADING.CANCEL_FAILED"
	// ErrHeartbeatFailed is returned when the keepalive POST fails.
	ErrHeartbeatFailed = "TRADING.HEARTBEAT_FAILED"
	// ErrStateNotFound is returned when no window state exists for an asset.
	ErrStateNotFound = "TRADING.STATE_NOT_FOUND"
	// ErrStateSaveFailed is returned when persisting window state fails.
	ErrStateSaveFailed = "TRADING.STATE_SAVE_FAILED"
	// ErrInvalidPrice is returned when a price is out of range.
	ErrInvalidPrice = "TRADING.INVALID_PRICE"
	// ErrInvalidSize is returned when a size is below minimum.
	ErrInvalidSize = "TRADING.INVALID_SIZE"
)
