package dto

// Input identifies the order to cancel.
type Input struct {
	OrderID     string // local UUID from order domain
	ClobOrderID string // CLOB-assigned order ID (returned by Submit)
}
