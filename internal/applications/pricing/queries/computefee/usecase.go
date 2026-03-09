package computefee

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/computefee/dto"
	"github.com/darmayasa221/polymarket-go/internal/applications/shared/feecalc"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

// FeeFormulaCurveConstant is C in: feePerShare = p × (1-p) × C
// Verified against live API research table (docs/decisions/5m-market-mechanics.md):
//
//	p=0.50 → 0.015625 (156 bps)  ✓
//	p=0.20 → 0.010000 (100 bps)  ✓
//	p=0.95 → 0.002969 (~30 bps)  ✓
//
// C = 1/16 = 0.0625 is the exact constant (derivation: 0.25 × C = 0.015625).
const FeeFormulaCurveConstant = "0.0625"

const (
	errInvalidPrice  = "PRICING.INVALID_PRICE"
	errPriceOutRange = "PRICING.PRICE_OUT_OF_RANGE"
)

type useCase struct{}

// New creates a ComputeFee query use case.
// No dependencies needed — the formula is purely mathematical.
func New() UseCase {
	return &useCase{}
}

// Execute computes the parabolic fee for a given token price.
func (uc *useCase) Execute(_ context.Context, input dto.Input) (dto.Output, error) {
	if input.TokenPrice == "" {
		return dto.Output{}, errtypes.NewClientError(errInvalidPrice)
	}

	p, err := decimal.NewFromString(input.TokenPrice)
	if err != nil {
		return dto.Output{}, errtypes.NewClientError(errInvalidPrice)
	}

	zero := decimal.Zero
	one := decimal.NewFromInt(1)
	if p.LessThanOrEqual(zero) || p.GreaterThan(one) {
		return dto.Output{}, errtypes.NewClientError(errPriceOutRange)
	}

	curveConstant := decimal.RequireFromString(FeeFormulaCurveConstant)
	oneMinusP := one.Sub(p)

	// feePerShare = p × (1-p) × C
	feePerShare := p.Mul(oneMinusP).Mul(curveConstant)

	// EffectiveBps = round(feePerShare × 10000)
	bps := feePerShare.Mul(decimal.NewFromInt(10000)).Round(0).IntPart()

	return dto.Output{
		Fee: feecalc.FeeResult{
			TokenPrice:   p,
			FeePerShare:  feePerShare,
			EffectiveBps: bps,
		},
	}, nil
}
