package sqlite

import (
	"context"
	"database/sql"
	"errors"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// GetPassword retrieves the hashed password for the given username.
// Returns NotFoundError if not found, InternalServerError on failure.
func (r *Repository) GetPassword(ctx context.Context, username string) (user.HashedPassword, error) {
	var hash string
	err := r.db.QueryRowContext(ctx, `SELECT hashed_password FROM users WHERE username = ? LIMIT 1`, username).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errtypes.NewNotFoundError(repository.ErrUserNotFound)
		}
		return "", errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}
	return user.HashedPassword(hash), nil
}
