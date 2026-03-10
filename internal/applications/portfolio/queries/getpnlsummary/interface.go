package getpnlsummary

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpnlsummary/dto"
)

// UseCase defines the GetPnLSummary query contract.
type UseCase interface {
	Execute(ctx context.Context, input dto.Input) (dto.Output, error)
}
