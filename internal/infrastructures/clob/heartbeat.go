package clob

import (
	"context"
	"fmt"
	"net/http"

	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
)

// Compile-time assertion: HeartbeatSender implements tradingports.HeartbeatSender.
var _ tradingports.HeartbeatSender = (*HeartbeatSender)(nil)

// HeartbeatSender sends POST /keep-alive to the CLOB every 5 seconds.
// Without this, all GTD orders auto-cancel after 10 seconds.
type HeartbeatSender struct{ client *Client }

// NewHeartbeatSender creates a HeartbeatSender.
func NewHeartbeatSender(client *Client) *HeartbeatSender { return &HeartbeatSender{client: client} }

// Send posts the keepalive. Call every 5 seconds while orders are open.
func (s *HeartbeatSender) Send(ctx context.Context) error {
	if err := s.client.do(ctx, http.MethodPost, "/keep-alive", nil, nil); err != nil {
		return fmt.Errorf("heartbeat: send: %w", err)
	}
	return nil
}
