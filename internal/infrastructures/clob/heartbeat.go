package clob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
)

// Compile-time assertion: HeartbeatSender implements tradingports.HeartbeatSender.
var _ tradingports.HeartbeatSender = (*HeartbeatSender)(nil)

// HeartbeatSender sends POST /v1/heartbeats to the CLOB every 5 seconds.
// Without this, all GTD orders auto-cancel after 10 seconds.
// The CLOB uses a chain mechanism: each heartbeat returns the ID for the next one.
type HeartbeatSender struct {
	client      *Client
	heartbeatID *string // nil until the first chain is established
}

// NewHeartbeatSender creates a HeartbeatSender.
func NewHeartbeatSender(client *Client) *HeartbeatSender { return &HeartbeatSender{client: client} }

type heartbeatRequest struct {
	HeartbeatID *string `json:"heartbeat_id"`
}

type heartbeatResponse struct {
	HeartbeatID string `json:"heartbeat_id"`
}

// Send posts the keepalive. Call every 5 seconds while orders are open.
// On the first call (no chain yet), it sends null and uses the ID returned by the
// 400 error to seed the chain, then immediately sends a valid heartbeat.
func (s *HeartbeatSender) Send(ctx context.Context) error {
	if s.heartbeatID == nil {
		// No chain yet — seed it by sending null, extract returned ID from error response.
		id, err := s.seedChain(ctx)
		if err != nil {
			return fmt.Errorf("heartbeat: seed chain: %w", err)
		}
		s.heartbeatID = &id
	}

	nextID, err := s.send(ctx, s.heartbeatID)
	if err != nil {
		// Chain broken (e.g. long gap) — reset to reseed on next tick.
		s.heartbeatID = nil
		return fmt.Errorf("heartbeat: send: %w", err)
	}
	s.heartbeatID = &nextID
	return nil
}

// seedChain sends heartbeat_id=null and extracts the seed ID from the 400 error response.
func (s *HeartbeatSender) seedChain(ctx context.Context) (string, error) {
	respData, statusCode, err := s.rawPost(ctx, nil)
	if err != nil {
		return "", err
	}
	// Server returns 400 with {"error": "...", "heartbeat_id": "<seed>"} on first call.
	if statusCode == http.StatusBadRequest {
		var errResp heartbeatResponse
		if jsonErr := json.Unmarshal(respData, &errResp); jsonErr == nil && errResp.HeartbeatID != "" {
			return errResp.HeartbeatID, nil
		}
	}
	// If 200 on first call (unexpected), parse normally.
	var resp heartbeatResponse
	if jsonErr := json.Unmarshal(respData, &resp); jsonErr == nil && resp.HeartbeatID != "" {
		return resp.HeartbeatID, nil
	}
	return "", fmt.Errorf("unexpected response seeding chain: status=%d body=%s", statusCode, respData)
}

// send posts a heartbeat with the given ID and returns the next ID.
func (s *HeartbeatSender) send(ctx context.Context, id *string) (string, error) {
	respData, statusCode, err := s.rawPost(ctx, id)
	if err != nil {
		return "", err
	}
	if statusCode != http.StatusOK {
		return "", fmt.Errorf("status %d: %s", statusCode, respData)
	}
	var resp heartbeatResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return resp.HeartbeatID, nil
}

// rawPost executes the POST /v1/heartbeats request and returns raw body + status code.
func (s *HeartbeatSender) rawPost(ctx context.Context, id *string) (body []byte, status int, err error) {
	var marshalErr error
	body, marshalErr = json.Marshal(heartbeatRequest{HeartbeatID: id})
	if marshalErr != nil {
		return nil, 0, fmt.Errorf("marshal: %w", marshalErr)
	}

	req, err := buildCLOBRequest(ctx, s.client.cfg, http.MethodPost, "/v1/heartbeats", body)
	if err != nil {
		return nil, 0, err
	}

	resp, err := s.client.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	return buf, resp.StatusCode, nil
}
