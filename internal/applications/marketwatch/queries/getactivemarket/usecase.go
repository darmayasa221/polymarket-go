package getactivemarket

import (
	"context"

	mwcmds "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands"
	mwports "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/getactivemarket/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	repo mwports.MarketRepository
}

// New creates a GetActiveMarket query use case.
func New(repo mwports.MarketRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute retrieves the active market for a given asset and window.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Asset == "" {
		return dto.Output{}, errtypes.NewClientError(mwcmds.ErrMarketNotFound)
	}

	asset := market.Asset(input.Asset)
	windowStart := input.WindowStart
	if windowStart.IsZero() {
		windowStart = timeutil.WindowStart(timeutil.Now())
	}

	m, err := uc.repo.FindByAssetAndWindow(ctx, asset, windowStart)
	if err != nil {
		return dto.Output{}, errtypes.NewNotFoundError(mwcmds.ErrMarketNotFound)
	}

	return dto.Output{
		MarketID:    m.ID(),
		Asset:       string(m.Asset()),
		ConditionID: string(m.ConditionID()),
		UpTokenID:   string(m.UpTokenID()),
		DownTokenID: string(m.DownTokenID()),
		TickSize:    m.TickSize().String(),
		WindowStart: m.WindowStart().Format("2006-01-02T15:04:05Z07:00"),
		WindowEnd:   m.WindowEnd().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
