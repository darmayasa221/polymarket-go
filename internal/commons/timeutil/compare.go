package timeutil

import "time"

// IsExpired returns true if t is before now (i.e., the time has passed).
func IsExpired(t time.Time) bool {
	return t.UTC().Before(Now())
}

// IsAfter returns true if t is after the reference time.
func IsAfter(t, reference time.Time) bool {
	return t.UTC().After(reference.UTC())
}
