package dto

// Output reports how many markets were refreshed.
type Output struct {
	Refreshed int
	Assets    []string
}
