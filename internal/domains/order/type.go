package order

// OrderType defines the time-in-force behavior of an order.
type OrderType string

const (
	// GTD is Good Till Date — expires at the specified time.
	GTD OrderType = "GTD"
	// GTC is Good Till Canceled — remains open until explicitly canceled.
	GTC OrderType = "GTC"
	// FOK is Fill or Kill — must fill immediately in full or cancel.
	FOK OrderType = "FOK"
	// FAK is Fill and Kill — fills what it can immediately, cancels the rest.
	FAK OrderType = "FAK"
)
