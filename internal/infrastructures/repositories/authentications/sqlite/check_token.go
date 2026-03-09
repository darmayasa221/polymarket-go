package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
)

// CheckToken retrieves a token by its value and verifies it has not expired.
// Returns NotFoundError if absent, AuthenticationError if expired, InternalServerError on scan failure.
func (r *Repository) CheckToken(ctx context.Context, value token.TokenValue) (*token.Token, error) {
	const query = `SELECT id, user_id, value, type, purpose, expires_at, created_at FROM tokens WHERE value = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, value.String())

	var (
		id        string
		userID    string
		val       string
		tokenType string
		purpose   string
		expiresAt time.Time
		createdAt time.Time
	)
	err := row.Scan(&id, &userID, &val, &tokenType, &purpose, &expiresAt, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtypes.NewNotFoundError(repository.ErrTokenNotFound)
		}
		return nil, errtypes.NewInternalServerError(repository.ErrTokenCheckFailed)
	}

	t := token.Reconstitute(token.ReconstitutedParams{
		ID:        id,
		UserID:    userID,
		Value:     val,
		Type:      tokenType,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	})

	if t.IsExpired() {
		return nil, errtypes.NewAuthenticationError(repository.ErrTokenExpired)
	}

	return t, nil
}
