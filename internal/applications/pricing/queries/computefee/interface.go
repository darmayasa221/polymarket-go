package computefee

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/computefee/dto"
)

// UseCase defines the ComputeFee query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
