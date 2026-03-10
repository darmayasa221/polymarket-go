package getpositions

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpositions/dto"
)

// UseCase defines the GetPositions query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
