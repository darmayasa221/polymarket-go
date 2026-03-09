package order

// OrderStatus tracks the lifecycle state of an order.
type OrderStatus string

const (
	// StatusOpen means the order is active on the CLOB.
	StatusOpen OrderStatus = "open"
	// StatusFilled means the order has been fully matched.
	StatusFilled OrderStatus = "filled"
	// StatusCancelled means the order was canceled before filling.
	StatusCancelled OrderStatus = "canceled"
	// StatusExpired means a GTD order passed its expiration time.
	StatusExpired OrderStatus = "expired"
)
