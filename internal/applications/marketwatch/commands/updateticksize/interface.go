package updateticksize

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/updateticksize/dto"
)

// UseCase defines the UpdateTickSize command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
