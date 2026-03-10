package placeorder

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/placeorder/dto"
)

// UseCase defines the PlaceOrder command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
