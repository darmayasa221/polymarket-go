package marktomarket

import (
	"context"

	"github.com/shopspring/decimal"

	portfoliocmds "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands"
	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/marktomarket/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	repo portfolioports.PositionRepository
}

// New creates a MarkToMarket query use case.
func New(repo portfolioports.PositionRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute computes unrealized PnL for all open positions at current prices.
// Results are aggregate-only — no per-position PnL is persisted.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	positions, err := uc.repo.ListOpen(ctx)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(portfoliocmds.ErrSaveFailed)
	}

	marks := make([]dto.PositionMark, 0, len(positions))
	total := decimal.Zero

	for _, pos := range positions {
		priceStr, ok := input.Prices[string(pos.TokenID())]
		if !ok {
			priceStr = pos.AvgPrice().String()
		}
		currentPrice, err := decimal.NewFromString(priceStr)
		if err != nil {
			currentPrice = pos.AvgPrice()
		}

		unrealisedPnL := pos.UnrealisedPnL(currentPrice)
		total = total.Add(unrealisedPnL)

		marks = append(marks, dto.PositionMark{
			PositionID:    pos.ID(),
			Asset:         string(pos.Asset()),
			Outcome:       string(pos.Outcome()),
			Size:          pos.Size().String(),
			AvgPrice:      pos.AvgPrice().String(),
			CurrentPrice:  currentPrice.String(),
			UnrealisedPnL: unrealisedPnL.String(),
		})
	}

	return dto.Output{
		Marks:              marks,
		TotalUnrealisedPnL: total.String(),
	}, nil
}
