package recordprice

import (
	"context"

	"github.com/shopspring/decimal"

	pricingcmds "github.com/darmayasa221/polymarket-go/internal/applications/pricing/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/commands/recordprice/dto"
	pricingports "github.com/darmayasa221/polymarket-go/internal/applications/pricing/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	priceRepo pricingports.PriceRepository
}

// New creates a RecordPrice use case.
func New(priceRepo pricingports.PriceRepository) UseCase {
	return &useCase{priceRepo: priceRepo}
}

// Execute records a new price observation from an external feed.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Asset == "" {
		return dto.Output{}, errtypes.NewClientError(pricingcmds.ErrInvalidAsset)
	}

	source := oracle.PriceSource(input.Source)
	if !source.IsValid() {
		return dto.Output{}, errtypes.NewClientError(pricingcmds.ErrInvalidSource)
	}

	value, err := decimal.NewFromString(input.Value)
	if err != nil {
		return dto.Output{}, errtypes.NewClientError(pricingcmds.ErrInvalidAsset)
	}

	price, err := oracle.New(oracle.Params{
		Asset:      input.Asset,
		Source:     source,
		Value:      value,
		RoundedAt:  input.RoundedAt,
		ReceivedAt: input.ReceivedAt,
	})
	if err != nil {
		return dto.Output{}, err
	}

	if err := uc.priceRepo.Save(ctx, price); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(pricingcmds.ErrSaveFailed)
	}

	return dto.Output{
		Asset:      price.Asset(),
		Source:     string(price.Source()),
		Value:      price.Value().String(),
		RecordedAt: price.ReceivedAt(),
	}, nil
}
