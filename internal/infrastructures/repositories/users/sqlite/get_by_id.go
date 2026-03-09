package sqlite

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// GetByID retrieves a user by their unique ID.
// Returns NotFoundError if the user does not exist, InternalServerError on scan failure.
func (r *Repository) GetByID(ctx context.Context, id user.UserID) (*user.User, error) {
	const query = `SELECT id, username, email, hashed_password, full_name, created_at, updated_at FROM users WHERE id = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, id.String())
	return scanUser(row)
}
