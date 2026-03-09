package redis

import (
	"context"
	"encoding/json"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// Add serializes the token to JSON and stores it in Redis with a TTL derived from ExpiresAt.
// Also registers the token key in the user's token set for bulk deletion.
// Returns an error if the token is already expired or if serialization/storage fails.
func (r *Repository) Add(ctx context.Context, t *token.Token) error {
	ttl := t.ExpiresAt().Sub(timeutil.Now())
	if ttl <= 0 {
		return errtypes.NewAuthenticationError(repository.ErrTokenExpired)
	}

	data, err := json.Marshal(tokenJSON{
		ID:        t.ID().String(),
		UserID:    t.UserID().String(),
		Value:     t.Value().String(),
		Type:      t.Type(),
		Purpose:   t.Purpose(),
		ExpiresAt: t.ExpiresAt(),
		CreatedAt: t.CreatedAt(),
	})
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenAddFailed)
	}

	tokenKey := redisclient.AuthTokenKey(t.Value().String())
	if err := r.cache.Set(ctx, tokenKey, string(data), ttl); err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenAddFailed)
	}

	userSetKey := redisclient.AuthUserSetKey(t.UserID().String())
	if err := r.cache.Client().SAdd(ctx, userSetKey, tokenKey).Err(); err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenAddFailed)
	}

	// Extend the user-set TTL only when the new token's lifetime is longer than the current TTL.
	// ExpireGT prevents a short-lived access token (15 m) from shrinking the set's TTL
	// below a long-lived refresh token's remaining lifetime (up to 168 h).
	_ = r.cache.Client().ExpireGT(ctx, userSetKey, ttl)

	return nil
}
