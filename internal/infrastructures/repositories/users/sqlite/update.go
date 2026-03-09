package sqlite

import (
	"context"
	"strings"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
)

// Update persists changes to an existing user record.
// Returns NotFoundError if the user does not exist, InternalServerError on failure.
func (r *Repository) Update(ctx context.Context, u *user.User) error {
	const query = `UPDATE users SET username=?, email=?, hashed_password=?, full_name=?, updated_at=? WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		u.Username(), u.Email().String(), u.HashedPassword().String(), u.FullName(), u.UpdatedAt(), u.ID().String(),
	)
	if err != nil {
		if strings.Contains(err.Error(), sqlitedb.ConstraintUniqueFailed) {
			if strings.Contains(err.Error(), constraintEmailColumn) {
				return errtypes.NewConflictError(repository.ErrEmailTaken)
			}
			return errtypes.NewConflictError(repository.ErrUsernameTaken)
		}
		return errtypes.NewInternalServerError(repository.ErrUserUpdateFailed)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrUserUpdateFailed)
	}
	if rows == 0 {
		return errtypes.NewNotFoundError(repository.ErrUserNotFound)
	}
	return nil
}
