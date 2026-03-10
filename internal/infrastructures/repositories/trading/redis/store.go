// Package redis implements the WindowStateStore using Redis.
package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion: Store implements tradingports.WindowStateStore.
var _ tradingports.WindowStateStore = (*Store)(nil)

const (
	errWindowStateNotFound = "TRADING.WINDOW_STATE_NOT_FOUND"
	keyPrefix              = "polymarket:window:"
)

// Store is a Redis implementation of tradingports.WindowStateStore.
// State survives bot restarts (unlike in-memory). StartWindow overwrites on each new window.
type Store struct{ client *goredis.Client }

// New creates a Store using the given Redis client.
func New(client *goredis.Client) *Store { return &Store{client: client} }

// SaveWindowState serializes state as JSON and writes it to Redis keyed by asset.
func (s *Store) SaveWindowState(ctx context.Context, state windowstate.WindowState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("redis window store: marshal: %w", err)
	}
	if err := s.client.Set(ctx, keyPrefix+state.Asset, data, 0).Err(); err != nil {
		return fmt.Errorf("redis window store: set: %w", err)
	}
	return nil
}

// GetWindowState retrieves and deserializes window state for asset.
// Returns NotFoundError when no state exists for the asset.
func (s *Store) GetWindowState(ctx context.Context, asset string) (windowstate.WindowState, error) {
	data, err := s.client.Get(ctx, keyPrefix+asset).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return windowstate.WindowState{}, errtypes.NewNotFoundError(errWindowStateNotFound)
		}
		return windowstate.WindowState{}, fmt.Errorf("redis window store: get: %w", err)
	}
	var state windowstate.WindowState
	if err := json.Unmarshal(data, &state); err != nil {
		return windowstate.WindowState{}, fmt.Errorf("redis window store: unmarshal: %w", err)
	}
	return state, nil
}
