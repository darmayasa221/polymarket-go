package dto

// Output reports how many orders were expired on window close.
type Output struct {
	Asset         string
	OrdersExpired int
}
