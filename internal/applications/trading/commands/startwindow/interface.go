package startwindow

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow/dto"
)

// UseCase defines the StartWindow command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
