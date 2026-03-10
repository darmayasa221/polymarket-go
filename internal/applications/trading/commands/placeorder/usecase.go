package placeorder

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	tradingcmds "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/placeorder/dto"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	"github.com/darmayasa221/polymarket-go/internal/commons/crypto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

// polygonChainID is the chain ID for Polygon mainnet.
const polygonChainID = 137

// minOrderSize is the minimum number of shares per order on Polymarket.
const minOrderSize = "5"

type useCase struct {
	orderRepo tradingports.OrderRepository
	store     tradingports.WindowStateStore
}

// New creates a PlaceOrder use case.
func New(orderRepo tradingports.OrderRepository, store tradingports.WindowStateStore) UseCase {
	return &useCase{orderRepo: orderRepo, store: store}
}

// Execute prepares a limit order and returns its EIP-712 unsigned hash.
// CRITICAL: This use case NEVER performs ECDSA signing.
// The interfaces layer must sign UnsignedHash with the EOA private key.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	state, err := uc.store.GetWindowState(ctx, input.Asset)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrStateNotFound)
	}
	if state.Status != windowstate.WindowOpen {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrWindowNotOpen)
	}

	outcome, err := parseOutcome(input.Outcome)
	if err != nil {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
	}

	side, err := parseSide(input.Side)
	if err != nil {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
	}

	minSize := decimal.RequireFromString(minOrderSize)
	if input.Size.LessThan(minSize) {
		return dto.Output{}, errtypes.NewClientError(tradingcmds.ErrInvalidSize)
	}

	expiry := order.GTDExpiration(state.WindowEnd)

	// MakerAmount and TakerAmount for a buy: maker pays USDC, taker receives shares.
	// For buy: makerAmount = price × size, takerAmount = size.
	makerAmount := input.Price.Mul(input.Size)
	takerAmount := input.Size

	saltBig := crypto.GenerateSalt()

	unsignedOrder := order.UnsignedOrder{
		Salt:          saltBig,
		Maker:         input.FunderAddress,
		Signer:        input.FunderAddress,
		Taker:         "0x0000000000000000000000000000000000000000",
		TokenID:       polyid.TokenID(input.TokenID),
		MakerAmount:   makerAmount,
		TakerAmount:   takerAmount,
		Expiration:    expiry.Unix(),
		Nonce:         0,
		FeeRateBps:    input.FeeRateBps,
		Side:          side,
		SignatureType: 0, // EOA
	}

	hash, err := order.SigningHash(unsignedOrder, polygonChainID)
	if err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrSaveFailed)
	}

	o, err := order.New(order.Params{
		MarketID:      state.MarketID,
		TokenID:       polyid.TokenID(input.TokenID),
		Side:          side,
		Outcome:       outcome,
		Price:         input.Price,
		Size:          input.Size,
		Type:          order.GTD,
		Expiration:    expiry,
		FeeRateBps:    input.FeeRateBps,
		SignatureType: 0,
	})
	if err != nil {
		return dto.Output{}, err
	}

	if err := uc.orderRepo.Save(ctx, o); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrSaveFailed)
	}

	return dto.Output{
		OrderID:      o.ID().String(),
		UnsignedHash: hash,
		GTDExpiry:    expiry.Unix(),
		FeePerShare:  computeFeePerShare(input.Price),
	}, nil
}

// parseOutcome maps "Up"/"Down" string to market.Outcome.
func parseOutcome(s string) (market.Outcome, error) {
	switch s {
	case string(market.Up):
		return market.Up, nil
	case string(market.Down):
		return market.Down, nil
	}
	return "", errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
}

// parseSide maps "buy"/"sell" string to order.Side.
func parseSide(s string) (order.Side, error) {
	switch s {
	case "buy":
		return order.Buy, nil
	case "sell":
		return order.Sell, nil
	}
	return 0, errtypes.NewClientError(tradingcmds.ErrInvalidAsset)
}

// computeFeePerShare applies the parabolic fee formula: fee = p × (1-p) × 0.0625.
// Duplicated here to avoid cross-context import from pricing context.
func computeFeePerShare(p decimal.Decimal) decimal.Decimal {
	one := decimal.NewFromInt(1)
	c := decimal.NewFromFloat(0.0625)
	return p.Mul(one.Sub(p)).Mul(c)
}
