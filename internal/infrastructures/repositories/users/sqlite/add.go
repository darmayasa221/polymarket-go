package sqlite

import (
	"context"
	"strings"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
)

// Add inserts a new user into the database.
// Returns ConflictError on unique constraint violation, InternalServerError on other failures.
func (r *Repository) Add(ctx context.Context, u *user.User) error {
	const query = `INSERT INTO users (id, username, email, hashed_password, full_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		u.ID().String(), u.Username(), u.Email().String(),
		u.HashedPassword().String(), u.FullName(), u.CreatedAt(), u.UpdatedAt(),
	)
	if err != nil {
		if strings.Contains(err.Error(), sqlitedb.ConstraintUniqueFailed) {
			if strings.Contains(err.Error(), constraintEmailColumn) {
				return errtypes.NewConflictError(repository.ErrEmailTaken)
			}
			return errtypes.NewConflictError(repository.ErrUsernameTaken)
		}
		return errtypes.NewInternalServerError(repository.ErrUserAddFailed)
	}
	return nil
}
