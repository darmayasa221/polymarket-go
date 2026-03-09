package jwt

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
)

// Compile-time assertion: Manager implements security.TokenManager.
var _ security.TokenManager = (*Manager)(nil)

// Manager implements JWT token operations.
type Manager struct {
	cfg Config
}

// New creates a new JWT Manager.
func New(cfg Config) *Manager {
	return &Manager{cfg: cfg}
}
