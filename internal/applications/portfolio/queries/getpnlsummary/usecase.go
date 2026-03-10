package getpnlsummary

import (
	"context"

	"github.com/shopspring/decimal"

	portfoliocmds "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands"
	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpnlsummary/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	repo portfolioports.PositionRepository
}

// New creates a GetPnLSummary query use case.
func New(repo portfolioports.PositionRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute returns realized PnL summary from all closed positions.
// TotalUnrealisedPnL is always "0" — callers must use MarkToMarket for live figures.
func (uc *useCase) Execute(ctx context.Context, _ dto.Input) (dto.Output, error) {
	records, err := uc.repo.ListClosedWithExitPrice(ctx)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(portfoliocmds.ErrSaveFailed)
	}

	total := decimal.Zero
	winCount := 0
	lossCount := 0

	for _, r := range records {
		pnl := r.Pos.RealisedPnL(r.ExitPrice)
		total = total.Add(pnl)
		if pnl.IsPositive() {
			winCount++
		} else {
			lossCount++
		}
	}

	return dto.Output{
		TotalRealisedPnL:   total.String(),
		TotalUnrealisedPnL: "0",
		WinCount:           winCount,
		LossCount:          lossCount,
		TotalCount:         len(records),
	}, nil
}
