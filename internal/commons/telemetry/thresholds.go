package telemetry

import "time"

const (
	// SlowQueryThreshold is the threshold for flagging slow DB queries.
	SlowQueryThreshold = 100 * time.Millisecond
	// SlowRequestThreshold is the threshold for flagging slow HTTP requests.
	SlowRequestThreshold = 500 * time.Millisecond
)

// IsSlowQuery returns true if duration exceeds the slow query threshold.
func IsSlowQuery(d time.Duration) bool { return d > SlowQueryThreshold }

// IsSlowRequest returns true if duration exceeds the slow request threshold.
func IsSlowRequest(d time.Duration) bool { return d > SlowRequestThreshold }
