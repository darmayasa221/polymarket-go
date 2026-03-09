package sqlite

import (
	"context"
	"database/sql"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// ListCursor retrieves a page of users using cursor-based pagination on the id column.
// Returns InternalServerError on query or scan failure.
func (r *Repository) ListCursor(ctx context.Context, params pagination.CursorParams) (pagination.CursorResult[*user.User], error) {
	rows, err := r.queryCursorRows(ctx, params)
	if err != nil {
		return pagination.CursorResult[*user.User]{}, errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return pagination.CursorResult[*user.User]{}, err
	}

	return buildCursorResult(users, params), nil
}

// queryCursorRows executes the appropriate cursor query based on direction.
func (r *Repository) queryCursorRows(ctx context.Context, params pagination.CursorParams) (*sql.Rows, error) {
	limit := params.PageSize + 1
	if params.Forward {
		return r.listCursorForward(ctx, params.Cursor, limit)
	}
	return r.listCursorBackward(ctx, params.Cursor, limit)
}

// buildCursorResult assembles the CursorResult from the fetched user slice and params.
func buildCursorResult(users []*user.User, params pagination.CursorParams) pagination.CursorResult[*user.User] {
	hasMore := len(users) > params.PageSize
	if hasMore {
		users = users[:params.PageSize]
	}
	if len(users) == 0 {
		return pagination.NewCursorResult(users, "", "", false, false)
	}
	nextCursor, prevCursor := computeCursors(users, params, hasMore)
	return pagination.NewCursorResult(users, nextCursor, prevCursor, nextCursor != "", prevCursor != "")
}

// computeCursors derives next/prev cursor strings from the result page.
func computeCursors(users []*user.User, params pagination.CursorParams, hasMore bool) (next, prev string) {
	if params.Forward {
		return computeForwardCursors(users, params.Cursor, hasMore)
	}
	return computeBackwardCursors(users, params.Cursor, hasMore)
}

// computeForwardCursors returns next/prev cursors for a forward page.
func computeForwardCursors(users []*user.User, cursor string, hasMore bool) (next, prev string) {
	if hasMore {
		next = users[len(users)-1].ID().String()
	}
	if cursor != "" {
		prev = users[0].ID().String()
	}
	return next, prev
}

// computeBackwardCursors returns next/prev cursors for a backward page.
func computeBackwardCursors(users []*user.User, cursor string, hasMore bool) (next, prev string) {
	if hasMore {
		prev = users[0].ID().String()
	}
	if cursor != "" {
		next = users[len(users)-1].ID().String()
	}
	return next, prev
}

// listCursorForward queries users with id > cursor ordered ascending.
func (r *Repository) listCursorForward(ctx context.Context, cursor string, limit int) (*sql.Rows, error) {
	if cursor == "" {
		return r.db.QueryContext(ctx,
			`SELECT id, username, email, hashed_password, full_name, created_at, updated_at FROM users ORDER BY id ASC LIMIT ?`,
			limit,
		)
	}
	return r.db.QueryContext(ctx,
		`SELECT id, username, email, hashed_password, full_name, created_at, updated_at FROM users WHERE id > ? ORDER BY id ASC LIMIT ?`,
		cursor, limit,
	)
}

// listCursorBackward queries users with id < cursor ordered descending.
func (r *Repository) listCursorBackward(ctx context.Context, cursor string, limit int) (*sql.Rows, error) {
	if cursor == "" {
		return r.db.QueryContext(ctx,
			`SELECT id, username, email, hashed_password, full_name, created_at, updated_at FROM users ORDER BY id DESC LIMIT ?`,
			limit,
		)
	}
	return r.db.QueryContext(ctx,
		`SELECT id, username, email, hashed_password, full_name, created_at, updated_at FROM users WHERE id < ? ORDER BY id DESC LIMIT ?`,
		cursor, limit,
	)
}
