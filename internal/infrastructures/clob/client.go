package clob

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is the authenticated Polymarket CLOB HTTP client.
type Client struct {
	cfg  Config
	http *http.Client
}

// NewClient creates a new CLOB Client.
func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg, http: &http.Client{}}
}

// do executes an authenticated CLOB request and decodes the JSON response into dst.
// Pass nil dst to discard the response body.
func (c *Client) do(ctx context.Context, method, path string, body, dst any) error {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("clob: marshal body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.cfg.BaseURL+path, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("clob: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := setL2Headers(req, c.cfg, string(bodyBytes)); err != nil {
		return fmt.Errorf("clob: set auth headers: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("clob: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("clob: status %d: %s", resp.StatusCode, b)
	}

	if dst != nil {
		if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
			return fmt.Errorf("clob: decode response: %w", err)
		}
	}
	return nil
}
