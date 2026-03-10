package ports

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
)

// WindowStateStore is an ephemeral, in-memory store for per-asset window state.
// State is not persisted across bot restarts — StartWindow re-initializes on each window.
type WindowStateStore interface {
	SaveWindowState(ctx context.Context, state windowstate.WindowState) error
	GetWindowState(ctx context.Context, asset string) (windowstate.WindowState, error)
}
