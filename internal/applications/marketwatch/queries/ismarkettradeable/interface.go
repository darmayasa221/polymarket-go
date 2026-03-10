package ismarkettradeable

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/ismarkettradeable/dto"
)

// UseCase defines the IsMarketTradeable query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
