package sqlite

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// GetByUsername retrieves a user by their username.
// Returns NotFoundError if not found, InternalServerError on failure.
func (r *Repository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	const query = `SELECT id, username, email, hashed_password, full_name, created_at, updated_at FROM users WHERE username = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, username)
	return scanUser(row)
}
