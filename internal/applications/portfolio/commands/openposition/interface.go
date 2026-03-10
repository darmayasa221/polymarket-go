package openposition

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/openposition/dto"
)

// UseCase defines the OpenPosition command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
