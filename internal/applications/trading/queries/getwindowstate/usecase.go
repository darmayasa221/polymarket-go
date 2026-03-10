package getwindowstate

import (
	"context"

	tradingcmds "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/queries/getwindowstate/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	store tradingports.WindowStateStore
}

// New creates a GetWindowState query use case.
func New(store tradingports.WindowStateStore) UseCase {
	return &useCase{store: store}
}

// Execute retrieves the current ephemeral window state for an asset.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Asset == "" {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
	}

	state, err := uc.store.GetWindowState(ctx, input.Asset)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrStateNotFound)
	}

	return dto.Output{WindowState: state}, nil
}
