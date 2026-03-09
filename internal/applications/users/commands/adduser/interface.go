// Package adduser implements the add user command use case.
package adduser

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser/dto"
)

// UseCase is the write contract for registering a new user.
type UseCase interface {
	// Execute registers a new user. Returns ConflictError if username exists.
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
