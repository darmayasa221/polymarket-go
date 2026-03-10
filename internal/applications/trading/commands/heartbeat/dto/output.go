package dto

import "time"

// Output confirms the heartbeat was sent successfully.
type Output struct {
	SentAt time.Time
}
