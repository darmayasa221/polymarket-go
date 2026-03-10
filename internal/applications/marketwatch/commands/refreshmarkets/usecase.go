package refreshmarkets

import (
	"context"

	mwcmds "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/refreshmarkets/dto"
	mwports "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	source mwports.MarketSource
	repo   mwports.MarketRepository
}

// New creates a RefreshMarkets use case.
func New(source mwports.MarketSource, repo mwports.MarketRepository) UseCase {
	return &useCase{source: source, repo: repo}
}

// Execute fetches all active 5-minute markets and persists them.
func (uc *useCase) Execute(ctx context.Context, _ dto.Input) (dto.Output, error) {
	markets, err := uc.source.FetchActive5mMarkets(ctx)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(mwcmds.ErrFetchFailed)
	}

	assets := make([]string, 0, len(markets))
	for _, m := range markets {
		if err := uc.repo.Save(ctx, m); err != nil {
			return dto.Output{}, errtypes.NewInternalServerError(mwcmds.ErrSaveFailed)
		}
		assets = append(assets, string(m.Asset()))
	}

	return dto.Output{
		Refreshed: len(markets),
		Assets:    assets,
	}, nil
}
