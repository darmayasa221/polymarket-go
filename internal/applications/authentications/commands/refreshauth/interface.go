package refreshauth

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth/dto"
)

// UseCase is the write contract for refreshing an authentication token pair.
type UseCase interface {
	// Execute validates the given refresh token and returns a new token pair.
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
