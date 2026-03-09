package order

import "time"

const gtdBufferSeconds = 60

// GTDExpiration returns the expiration time for a GTD order.
// Formula: windowEnd + 60s mandatory buffer.
// The 60s buffer ensures the order stays live through the entire window despite clock skew.
func GTDExpiration(windowEnd time.Time) time.Time {
	return windowEnd.Add(gtdBufferSeconds * time.Second)
}
