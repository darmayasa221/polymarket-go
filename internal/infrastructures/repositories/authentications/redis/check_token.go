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

// CheckToken retrieves a token from Redis by its value, unmarshals it, and verifies it has not expired.
// Returns NotFoundError if the key is absent, AuthenticationError if expired, InternalServerError on failure.
func (r *Repository) CheckToken(ctx context.Context, value token.TokenValue) (*token.Token, error) {
	key := redisclient.AuthTokenKey(value.String())

	data, err := r.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, errtypes.NewNotFoundError(repository.ErrTokenNotFound)
		}
		return nil, errtypes.NewInternalServerError(repository.ErrTokenCheckFailed)
	}

	var tj tokenJSON
	if err := json.Unmarshal([]byte(data), &tj); err != nil {
		return nil, errtypes.NewInternalServerError(repository.ErrTokenCheckFailed)
	}

	t := token.Reconstitute(token.ReconstitutedParams{
		ID:        tj.ID,
		UserID:    tj.UserID,
		Value:     tj.Value,
		Type:      tj.Type,
		Purpose:   tj.Purpose,
		ExpiresAt: tj.ExpiresAt,
		CreatedAt: tj.CreatedAt,
	})

	if t.IsExpired() {
		return nil, errtypes.NewAuthenticationError(repository.ErrTokenExpired)
	}

	return t, nil
}
