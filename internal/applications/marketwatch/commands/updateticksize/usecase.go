package updateticksize

import (
	"context"

	"github.com/shopspring/decimal"

	mwcmds "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/updateticksize/dto"
	mwports "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

// validTickSizes contains the only tick sizes Polymarket supports.
// Tick size changes when price > 0.96 or < 0.04.
var validTickSizes = map[string]struct{}{
	"0.1":    {},
	"0.01":   {},
	"0.001":  {},
	"0.0001": {},
}

type useCase struct {
	repo mwports.MarketRepository
}

// New creates an UpdateTickSize use case.
func New(repo mwports.MarketRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute updates the tick size for a market after a WS tick_size_change event.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.ConditionID == "" {
		return dto.Output{}, errtypes.NewClientError(mwcmds.ErrMarketNotFound)
	}

	newTick, err := decimal.NewFromString(input.NewTickSize)
	if err != nil {
		return dto.Output{}, errtypes.NewClientError(mwcmds.ErrInvalidTickSize)
	}

	if _, ok := validTickSizes[input.NewTickSize]; !ok {
		return dto.Output{}, errtypes.NewClientError(mwcmds.ErrInvalidTickSize)
	}

	if err := uc.repo.UpdateTickSize(ctx, input.ConditionID, newTick); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(mwcmds.ErrUpdateFailed)
	}

	return dto.Output{
		ConditionID: input.ConditionID,
		NewTickSize: newTick.String(),
	}, nil
}
