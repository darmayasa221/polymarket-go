package redis

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// DeleteByUserID removes all tokens belonging to the given user from Redis.
// It reads the user's token set, then deletes all token keys and the set itself
// atomically in a single pipeline. This operation is idempotent — an empty set
// results in no error.
func (r *Repository) DeleteByUserID(ctx context.Context, userID valueobjects.ID) error {
	userSetKey := redisclient.AuthUserSetKey(userID.String())

	// Get all token keys belonging to this user.
	keys, err := r.cache.Client().SMembers(ctx, userSetKey).Result()
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenDeleteFailed)
	}

	// Delete all token keys + the set key atomically in one pipeline.
	// Del on a non-existent key is a no-op, so an empty set is handled correctly.
	pipe := r.cache.Client().Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}
	pipe.Del(ctx, userSetKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenDeleteFailed)
	}

	return nil
}
