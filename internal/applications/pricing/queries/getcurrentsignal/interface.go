package getcurrentsignal

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/getcurrentsignal/dto"
)

// UseCase defines the GetCurrentSignal query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
