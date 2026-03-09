package logoutuser

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser/dto"
)

// UseCase is the write contract for invalidating a user session token.
type UseCase interface {
	// Execute deletes the given token, ending the user session.
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
