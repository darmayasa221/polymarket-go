package sqlite

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// Delete removes a user by their unique ID.
// Returns NotFoundError if the user does not exist, InternalServerError on failure.
func (r *Repository) Delete(ctx context.Context, id user.UserID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id.String())
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrUserDeleteFailed)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrUserDeleteFailed)
	}
	if rows == 0 {
		return errtypes.NewNotFoundError(repository.ErrUserNotFound)
	}
	return nil
}
