package timeutil

import "time"

const windowSeconds = 300 // 5 minutes

// Formula: floor(unix / 300) * 300.
func WindowStart(t time.Time) time.Time {
	unix := t.UTC().Unix()
	start := (unix / windowSeconds) * windowSeconds
	return time.Unix(start, 0).UTC()
}

// WindowEnd returns the end of the 5-minute window containing t.
// Equal to WindowStart(t) + 5 minutes.
func WindowEnd(t time.Time) time.Time {
	return WindowStart(t).Add(windowSeconds * time.Second)
}

// SecondsRemaining returns how many seconds remain in the current 5-minute window.
func SecondsRemaining(t time.Time) int {
	end := WindowEnd(t)
	return int(end.Unix() - t.UTC().Unix())
}
