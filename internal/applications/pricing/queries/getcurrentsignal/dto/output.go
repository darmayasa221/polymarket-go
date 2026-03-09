package dto

import "github.com/darmayasa221/polymarket-go/internal/applications/shared/signal"

// Output wraps the signal computed by GetCurrentSignal.
type Output struct {
	Signal signal.Signal
}
