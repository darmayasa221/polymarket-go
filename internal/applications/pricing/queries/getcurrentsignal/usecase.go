package getcurrentsignal

import (
	"context"

	"github.com/shopspring/decimal"

	pricingports "github.com/darmayasa221/polymarket-go/internal/applications/pricing/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/getcurrentsignal/dto"
	"github.com/darmayasa221/polymarket-go/internal/applications/shared/signal"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

// confidenceThreshold is the price-move fraction that yields full confidence (1.0).
// A 5% move signals a strong directional bias in a 5-minute window.
const confidenceThreshold = "0.05"

// errNoOpenPrice is the error key when the window open price is unavailable.
const errNoOpenPrice = "PRICING.NO_OPEN_PRICE"

// errNoCurrentPrice is the error key when no current price is available.
const errNoCurrentPrice = "PRICING.NO_CURRENT_PRICE"

// errAssetRequired is the error key when asset is empty.
const errAssetRequired = "PRICING.ASSET_REQUIRED"

type useCase struct {
	priceRepo pricingports.PriceRepository
}

// New creates a GetCurrentSignal query use case.
func New(priceRepo pricingports.PriceRepository) UseCase {
	return &useCase{priceRepo: priceRepo}
}

// Execute computes the current directional signal for an asset.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Asset == "" {
		return dto.Output{}, errtypes.NewClientError(errAssetRequired)
	}

	windowStart := timeutil.WindowStart(timeutil.Now())

	openPriceObs, err := uc.priceRepo.WindowOpenPrice(ctx, input.Asset, windowStart)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(errNoOpenPrice)
	}

	// Prefer Chainlink as the current price source; fall back to any available source.
	currentObs, err := uc.priceRepo.LatestChainlinkByAsset(ctx, input.Asset)
	if err != nil {
		currentObs, err = uc.priceRepo.LatestByAsset(ctx, input.Asset)
		if err != nil {
			return dto.Output{}, errtypes.NewInternalServerError(errNoCurrentPrice)
		}
	}

	openPrice := openPriceObs.Value()
	currentPrice := currentObs.Value()

	predicted := oracle.PredictOutcome(openPrice, currentPrice)

	confidence := computeConfidence(openPrice, currentPrice)

	return dto.Output{
		Signal: signal.Signal{
			Asset:        input.Asset,
			Predicted:    string(predicted),
			Confidence:   confidence,
			OpenPrice:    openPrice,
			CurrentPrice: currentPrice,
			Source:       string(currentObs.Source()),
			RecordedAt:   currentObs.ReceivedAt(),
		},
	}, nil
}

// computeConfidence returns a 0.0–1.0 confidence score.
// A 5% or larger price move yields 1.0; flat price yields 0.0.
func computeConfidence(openPrice, currentPrice decimal.Decimal) decimal.Decimal {
	if openPrice.IsZero() {
		return decimal.Zero
	}

	threshold := decimal.RequireFromString(confidenceThreshold)
	delta := currentPrice.Sub(openPrice).Abs()
	// confidence = |delta| / openPrice / 0.05
	confidence := delta.Div(openPrice).Div(threshold)

	one := decimal.NewFromInt(1)
	if confidence.GreaterThan(one) {
		return one
	}
	return confidence
}
