package getwindowstate

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/queries/getwindowstate/dto"
)

// UseCase defines the GetWindowState query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
