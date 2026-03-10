package getpositions

import (
	"context"

	portfoliocmds "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands"
	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpositions/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	repo portfolioports.PositionRepository
}

// New creates a GetPositions query use case.
func New(repo portfolioports.PositionRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute retrieves open positions, optionally filtered by asset.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	positions, err := uc.repo.ListOpen(ctx)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(portfoliocmds.ErrSaveFailed)
	}

	result := make([]dto.PositionDTO, 0, len(positions))
	for _, pos := range positions {
		if input.Asset != "" && market.Asset(input.Asset) != pos.Asset() {
			continue
		}
		result = append(result, dto.PositionDTO{
			PositionID: pos.ID(),
			Asset:      string(pos.Asset()),
			TokenID:    string(pos.TokenID()),
			Outcome:    string(pos.Outcome()),
			Size:       pos.Size().String(),
			AvgPrice:   pos.AvgPrice().String(),
			OpenedAt:   pos.OpenedAt(),
		})
	}

	return dto.Output{Positions: result}, nil
}
