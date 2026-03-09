// Package repository defines the User repository contract.
// The interface lives in domain — implementations live in infrastructure.
package repository

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// User defines all persistence operations for the User aggregate.
// All methods accept context.Context for cancellation and tracing.
type User interface {
	// Add persists a new user. Returns ConflictError if username/email exists.
	Add(ctx context.Context, u *user.User) error

	// GetByID retrieves a user by ID. Returns NotFoundError if not found.
	GetByID(ctx context.Context, id user.UserID) (*user.User, error)

	// GetByUsername retrieves a user by username. Returns NotFoundError if not found.
	GetByUsername(ctx context.Context, username string) (*user.User, error)

	// GetIDByUsername retrieves only the user's ID by username.
	GetIDByUsername(ctx context.Context, username string) (user.UserID, error)

	// GetPassword retrieves the hashed password for a username.
	GetPassword(ctx context.Context, username string) (user.HashedPassword, error)

	// VerifyUsername returns true if the username already exists.
	VerifyUsername(ctx context.Context, username string) (bool, error)

	// ListOffset retrieves a paginated list using offset pagination.
	ListOffset(ctx context.Context, params pagination.OffsetParams) (pagination.OffsetResult[*user.User], error)

	// ListCursor retrieves a paginated list using cursor pagination.
	ListCursor(ctx context.Context, params pagination.CursorParams) (pagination.CursorResult[*user.User], error)

	// Update persists changes to an existing user.
	Update(ctx context.Context, u *user.User) error

	// Delete removes a user by ID.
	Delete(ctx context.Context, id user.UserID) error
}
