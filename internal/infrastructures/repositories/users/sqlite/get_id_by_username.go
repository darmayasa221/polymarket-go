package sqlite

import (
	"context"
	"database/sql"
	"errors"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// GetIDByUsername retrieves only the user's ID for the given username.
// Returns NotFoundError if not found, InternalServerError on failure.
func (r *Repository) GetIDByUsername(ctx context.Context, username string) (user.UserID, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE username = ? LIMIT 1`, username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errtypes.NewNotFoundError(repository.ErrUserNotFound)
		}
		return "", errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}
	return user.UserID(id), nil
}
