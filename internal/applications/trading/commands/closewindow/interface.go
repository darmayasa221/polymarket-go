package closewindow

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/closewindow/dto"
)

// UseCase defines the CloseWindow command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
