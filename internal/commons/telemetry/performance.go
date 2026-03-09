package telemetry

import "time"

// Elapsed returns the duration since start in milliseconds.
func Elapsed(start time.Time) float64 {
	return float64(time.Since(start).Milliseconds())
}
