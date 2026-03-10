package ports

import "context"

// HeartbeatSender sends the POST keepalive to the CLOB every 5 seconds.
// Without this, all open orders auto-cancel after 10 seconds.
type HeartbeatSender interface {
	Send(ctx context.Context) error
}
