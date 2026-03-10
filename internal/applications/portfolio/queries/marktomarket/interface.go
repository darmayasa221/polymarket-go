package marktomarket

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/marktomarket/dto"
)

// UseCase defines the MarkToMarket query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
