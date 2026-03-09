package sqlite

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
)

// VerifyUsername returns true if a user with the given username already exists.
// Returns InternalServerError on query failure.
func (r *Repository) VerifyUsername(ctx context.Context, username string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)
	if err != nil {
		return false, errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}
	return count > 0, nil
}
