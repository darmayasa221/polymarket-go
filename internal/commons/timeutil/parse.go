package timeutil

import "time"

// ParseRFC3339 parses an RFC3339 time string into UTC time.
func ParseRFC3339(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

// FormatRFC3339 formats a time as RFC3339 UTC string.
func FormatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
