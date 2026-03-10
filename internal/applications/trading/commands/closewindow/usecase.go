package closewindow

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	tradingcmds "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/closewindow/dto"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	store     tradingports.WindowStateStore
	orderRepo tradingports.OrderRepository
}

// New creates a CloseWindow use case.
func New(store tradingports.WindowStateStore, orderRepo tradingports.OrderRepository) UseCase {
	return &useCase{store: store, orderRepo: orderRepo}
}

// Execute closes the trading window for an asset and expires all open orders.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Asset == "" {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
	}

	state, err := uc.store.GetWindowState(ctx, input.Asset)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrStateNotFound)
	}

	if state.Status != windowstate.WindowOpen {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrWindowNotOpen)
	}

	// Expire all open orders — best-effort; window close takes priority.
	openOrders, err := uc.orderRepo.ListOpenByMarket(ctx, state.MarketID)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrSaveFailed)
	}
	for _, o := range openOrders {
		_ = uc.orderRepo.UpdateStatus(ctx, o.ID(), order.StatusExpired)
	}

	state.Status = windowstate.WindowClosed
	state.OpenOrders = []windowstate.OrderSummary{}
	if err := uc.store.SaveWindowState(ctx, state); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrStateSaveFailed)
	}

	return dto.Output{
		Asset:         state.Asset,
		OrdersExpired: len(openOrders),
	}, nil
}
