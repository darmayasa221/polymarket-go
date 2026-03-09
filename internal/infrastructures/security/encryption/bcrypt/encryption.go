// Package bcrypt provides the bcrypt encryption adapter for the infrastructures layer.
package bcrypt

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
)

// Compile-time assertion: Encryption implements security.Encryption.
var _ security.Encryption = (*Encryption)(nil)

// Encryption implements bcrypt password hashing.
type Encryption struct {
	cfg Config
}

// New creates a new Encryption adapter.
// If cfg.Cost is zero, DefaultCost is used.
func New(cfg Config) *Encryption {
	if cfg.Cost == 0 {
		cfg.Cost = DefaultCost
	}
	return &Encryption{cfg: cfg}
}
