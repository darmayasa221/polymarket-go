package ismarkettradeable

import (
	"context"

	mwcmds "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands"
	mwports "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/ismarkettradeable/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	repo mwports.MarketRepository
}

// New creates an IsMarketTradeable query use case.
func New(repo mwports.MarketRepository) UseCase {
	return &useCase{repo: repo}
}

// Execute checks whether a market is currently eligible for order placement.
// A market is tradeable when: it is in the active list AND fees are enabled.
// The `enable_order_book: true` constraint is enforced by the infrastructure layer
// when fetching from the Gamma API; here we check the persisted active state.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.ConditionID == "" {
		return dto.Output{}, errtypes.NewClientError(mwcmds.ErrMarketNotFound)
	}

	markets, err := uc.repo.ListActive(ctx)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(mwcmds.ErrFetchFailed)
	}

	for _, m := range markets {
		if string(m.ConditionID()) != input.ConditionID {
			continue
		}
		// Found the market — check tradeability conditions.
		if !m.FeeEnabled() {
			return dto.Output{
				Tradeable: false,
				Reason:    "fees not enabled",
				TickSize:  m.TickSize().String(),
				Active:    m.Active(),
			}, nil
		}
		return dto.Output{
			Tradeable: true,
			TickSize:  m.TickSize().String(),
			Active:    m.Active(),
		}, nil
	}

	return dto.Output{
		Tradeable: false,
		Reason:    "market not active",
		Active:    false,
	}, nil
}
