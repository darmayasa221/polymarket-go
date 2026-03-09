package refreshmarkets

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/refreshmarkets/dto"
)

// UseCase defines the RefreshMarkets command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
