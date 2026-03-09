package loginuser

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser/dto"
)

// UseCase is the write contract for authenticating a user.
type UseCase interface {
	// Execute authenticates a user and returns a token pair.
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
