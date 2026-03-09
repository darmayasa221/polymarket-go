// Package redis implements the Authentication repository using Redis (Adapter Pattern).
package redis

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// Compile-time assertion: Repository implements repository.Authentication.
var _ repository.Authentication = (*Repository)(nil)

// Repository is the Redis implementation of repository.Authentication.
type Repository struct {
	cache *redisclient.Client
}

// New creates a new Redis Authentication repository.
func New(cache *redisclient.Client) *Repository {
	return &Repository{cache: cache}
}

// tokenJSON is the JSON representation of a token stored in Redis.
type tokenJSON struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	Purpose   string    `json:"purpose"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
