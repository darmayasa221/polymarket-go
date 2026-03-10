package cancelorder

import (
	"context"

	tradingcmds "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/cancelorder/dto"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	orderRepo tradingports.OrderRepository
	submitter tradingports.OrderSubmitter
}

// New creates a CancelOrder use case.
func New(orderRepo tradingports.OrderRepository, submitter tradingports.OrderSubmitter) UseCase {
	return &useCase{orderRepo: orderRepo, submitter: submitter}
}

// Execute cancels an open order on the CLOB and updates local status.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.OrderID == "" {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrOrderNotFound)
	}
	if input.ClobOrderID == "" {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrOrderNotFound)
	}

	if err := uc.submitter.Cancel(ctx, input.ClobOrderID); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrCancelFailed)
	}

	// Best-effort status update — CLOB cancel already succeeded.
	_ = uc.orderRepo.UpdateStatus(ctx, polyid.OrderID(input.OrderID), order.StatusCancelled)

	return dto.Output{OrderID: input.OrderID}, nil
}
