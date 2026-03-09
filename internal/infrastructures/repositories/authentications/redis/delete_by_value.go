package redis

import (
	"context"
	"encoding/json"
	"errors"

	goredis "github.com/redis/go-redis/v9"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// DeleteByValue removes a specific token from Redis by its value.
// It also attempts to look up the token's user ID to remove it from the user set (best-effort).
// Returns NotFoundError if the token key does not exist, InternalServerError on failure.
func (r *Repository) DeleteByValue(ctx context.Context, value token.TokenValue) error {
	key := redisclient.AuthTokenKey(value.String())

	// Read the token data first so we can remove it from the user set.
	data, err := r.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return errtypes.NewNotFoundError(repository.ErrTokenNotFound)
		}
		return errtypes.NewInternalServerError(repository.ErrTokenDeleteFailed)
	}

	// Parse to get the userID for set cleanup (best-effort).
	var tj tokenJSON
	if jsonErr := json.Unmarshal([]byte(data), &tj); jsonErr == nil {
		userSetKey := redisclient.AuthUserSetKey(tj.UserID)
		_ = r.cache.Client().SRem(ctx, userSetKey, key).Err()
	}

	if err := r.cache.Delete(ctx, key); err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenDeleteFailed)
	}

	return nil
}
