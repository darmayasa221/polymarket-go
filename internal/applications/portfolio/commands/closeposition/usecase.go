package closeposition

import (
	"context"

	"github.com/shopspring/decimal"

	portfoliocmds "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/closeposition/dto"
	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	repo portfolioports.PositionRepository
}

// New creates a ClosePosition use case.
func New(repo portfolioports.PositionRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute closes a position at the given exit price and returns realized PnL.
// This is the mid-window exit: stop-loss at −$0.20, take-profit at +$0.20.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.PositionID == "" {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrPositionNotFound)
	}

	exitPrice, err := decimal.NewFromString(input.ExitPrice)
	if err != nil || exitPrice.IsZero() || exitPrice.IsNegative() {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidPrice)
	}

	pos, err := uc.repo.FindByID(ctx, input.PositionID)
	if err != nil {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrPositionNotFound)
	}

	realisedPnL := pos.RealisedPnL(exitPrice)

	closedAt := timeutil.Now()
	if err := uc.repo.Close(ctx, input.PositionID, exitPrice, closedAt); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(portfoliocmds.ErrCloseFailed)
	}

	return dto.Output{
		PositionID:  input.PositionID,
		RealisedPnL: realisedPnL.String(),
	}, nil
}
