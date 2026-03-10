package closeposition

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/closeposition/dto"
)

// UseCase defines the ClosePosition command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
