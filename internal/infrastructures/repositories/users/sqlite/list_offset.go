package sqlite

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// ListOffset retrieves a page of users using offset-based pagination.
// Returns InternalServerError on query or scan failure.
func (r *Repository) ListOffset(ctx context.Context, params pagination.OffsetParams) (pagination.OffsetResult[*user.User], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return pagination.OffsetResult[*user.User]{}, errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}

	const query = `SELECT id, username, email, hashed_password, full_name, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, params.PageSize, params.Offset())
	if err != nil {
		return pagination.OffsetResult[*user.User]{}, errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return pagination.OffsetResult[*user.User]{}, err
	}
	return pagination.NewOffsetResult(users, total, params.Page, params.PageSize), nil
}
