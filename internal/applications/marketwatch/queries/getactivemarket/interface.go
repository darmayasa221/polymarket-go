package getactivemarket

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/getactivemarket/dto"
)

// UseCase defines the GetActiveMarket query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
