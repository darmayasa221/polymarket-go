package timeutil

import "time"

// AddDuration returns a new time by adding duration to now.
func AddDuration(d time.Duration) time.Time {
	return Now().Add(d)
}

// AddDays returns a new time by adding n days to now.
func AddDays(n int) time.Time {
	return Now().AddDate(0, 0, n)
}
