// Package cached implements a caching Decorator for the User repository.
// Wraps any repository.User implementation with Redis caching.
// Decorator Pattern: transparently adds caching without modifying the base.
package cached

import (
	"encoding/json"

	cacheconstants "github.com/darmayasa221/polymarket-go/internal/commons/constants/cache"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	rediscache "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// Compile-time assertion: Repository implements repository.User.
var _ repository.User = (*Repository)(nil)

// defaultTTL is the time-to-live for cached user entries.
const defaultTTL = cacheconstants.DefaultUserCacheTTL

// Repository wraps a base repository.User with Redis caching.
// Implements the Decorator Pattern — all calls delegate to base with caching behavior added.
type Repository struct {
	base  repository.User
	cache *rediscache.Client
}

// New creates a new cached User repository.
// base is the underlying implementation (e.g. SQLite), cache is the Redis client.
func New(base repository.User, cache *rediscache.Client) *Repository {
	return &Repository{base: base, cache: cache}
}

// toParams converts a *user.User into ReconstitutedParams for JSON serialization.
func toParams(u *user.User) user.ReconstitutedParams {
	return user.ReconstitutedParams{
		ID:             u.ID().String(),
		Username:       u.Username(),
		Email:          u.Email().String(),
		HashedPassword: u.HashedPassword().String(),
		FullName:       u.FullName(),
		CreatedAt:      u.CreatedAt(),
		UpdatedAt:      u.UpdatedAt(),
	}
}

// marshalUser serializes a *user.User to a JSON string. Returns empty string on failure.
func marshalUser(u *user.User) (string, bool) {
	data, err := json.Marshal(toParams(u))
	if err != nil {
		return "", false
	}
	return string(data), true
}

// unmarshalUser deserializes JSON data into a *user.User. Returns nil on failure.
func unmarshalUser(data string) *user.User {
	var p user.ReconstitutedParams
	if json.Unmarshal([]byte(data), &p) != nil {
		return nil
	}
	return user.Reconstitute(p)
}
