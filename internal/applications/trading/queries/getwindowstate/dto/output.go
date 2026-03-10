package dto

import "github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"

// Output wraps the ephemeral window state snapshot.
type Output struct {
	windowstate.WindowState
}
