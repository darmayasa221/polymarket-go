package openposition

import (
	"context"

	"github.com/shopspring/decimal"

	portfoliocmds "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/openposition/dto"
	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	repo portfolioports.PositionRepository
}

// New creates an OpenPosition use case.
func New(repo portfolioports.PositionRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute creates and persists a new position for a filled order.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	asset := market.Asset(input.Asset)
	if !asset.IsValid() {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidAsset)
	}
	if input.TokenID == "" {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidAsset)
	}
	outcome := market.Outcome(input.Outcome)
	if outcome != market.Up && outcome != market.Down {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidOutcome)
	}
	if input.MarketID == "" {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidAsset)
	}

	size, err := decimal.NewFromString(input.Size)
	if err != nil || size.IsZero() || size.IsNegative() {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidSize)
	}

	avgPrice, err := decimal.NewFromString(input.AvgPrice)
	if err != nil || avgPrice.IsZero() || avgPrice.IsNegative() {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidPrice)
	}

	pos, err := position.New(position.Params{
		Asset:    asset,
		TokenID:  polyid.TokenID(input.TokenID),
		Outcome:  outcome,
		Size:     size,
		AvgPrice: avgPrice,
		MarketID: input.MarketID,
	})
	if err != nil {
		return dto.Output{}, errtypes.NewClientError(portfoliocmds.ErrInvalidAsset)
	}

	if err := uc.repo.Save(ctx, pos); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(portfoliocmds.ErrSaveFailed)
	}

	return dto.Output{PositionID: pos.ID()}, nil
}
