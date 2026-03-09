package sqlite

import (
	"context"
	"strings"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
)

// Add inserts a new token into the database.
// Returns ConflictError on unique constraint violation, InternalServerError on other failures.
func (r *Repository) Add(ctx context.Context, t *token.Token) error {
	const query = `INSERT INTO tokens (id, user_id, value, type, purpose, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		t.ID().String(), t.UserID().String(), t.Value().String(),
		t.Type(), t.Purpose(), t.ExpiresAt(), t.CreatedAt(),
	)
	if err != nil {
		if strings.Contains(err.Error(), sqlitedb.ConstraintUniqueFailed) {
			return errtypes.NewConflictError(repository.ErrTokenValueTaken)
		}
		return errtypes.NewInternalServerError(repository.ErrTokenAddFailed)
	}
	return nil
}
