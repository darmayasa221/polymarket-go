package recordprice

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/commands/recordprice/dto"
)

// UseCase defines the RecordPrice command contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
