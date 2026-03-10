package heartbeat

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/heartbeat/dto"
)

// UseCase defines the Heartbeat command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
