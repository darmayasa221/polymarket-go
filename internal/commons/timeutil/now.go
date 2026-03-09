// Package timeutil provides time utilities ensuring UTC consistency.
// Always use timeutil.Now() — never time.Now() directly.
package timeutil

import "time"

// Now returns the current UTC time.
// Use this everywhere instead of time.Now() for UTC consistency.
func Now() time.Time {
	return time.Now().UTC()
}

// Unix returns the current UTC Unix timestamp.
func Unix() int64 {
	return Now().Unix()
}
