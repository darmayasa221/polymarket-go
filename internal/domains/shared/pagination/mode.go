// Package pagination defines pure domain pagination types.
// No HTTP-specific logic here — HTTP parsing is in interfaces layer.
package pagination

// Mode represents the pagination strategy.
type Mode string

const (
	ModeOffset Mode = "offset"
	ModeCursor Mode = "cursor"
)
