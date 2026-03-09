package getuser

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser/dto"
)

// UseCase is the read contract for retrieving a user by ID.
type UseCase interface {
	// Execute retrieves a user by ID. Returns NotFoundError if not found.
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
