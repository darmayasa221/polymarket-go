package startwindow

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	tradingcmds "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow/dto"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	store tradingports.WindowStateStore
}

// New creates a StartWindow use case.
func New(store tradingports.WindowStateStore) UseCase {
	return &useCase{store: store}
}

// Execute initializes the ephemeral window state for a new 5-minute trading window.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Asset == "" {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
	}
	if !market.Asset(input.Asset).IsValid() {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
	}
	if input.MarketID == "" {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
	}

	now := timeutil.Now()
	windowStart := timeutil.WindowStart(now)
	windowEnd := timeutil.WindowEnd(now)

	state := windowstate.WindowState{
		MarketID:    input.MarketID,
		Asset:       input.Asset,
		WindowStart: windowStart,
		WindowEnd:   windowEnd,
		ConditionID: input.ConditionID,
		UpTokenID:   input.UpTokenID,
		DownTokenID: input.DownTokenID,
		TickSize:    input.TickSize,
		OpenPrice:   input.OpenPrice,
		Status:      windowstate.WindowOpen,
		OpenOrders:  []windowstate.OrderSummary{},
	}

	if err := uc.store.SaveWindowState(ctx, state); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrStateSaveFailed)
	}

	return dto.Output{
		Asset:       state.Asset,
		MarketID:    state.MarketID,
		ConditionID: state.ConditionID,
		WindowStart: state.WindowStart,
		WindowEnd:   state.WindowEnd,
	}, nil
}
