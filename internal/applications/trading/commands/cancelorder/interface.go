package cancelorder

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/cancelorder/dto"
)

// UseCase defines the CancelOrder command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
