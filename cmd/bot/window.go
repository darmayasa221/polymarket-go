package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shopspring/decimal"

	refreshmarketsdto "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/refreshmarkets/dto"
	getactivemarketdto "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/getactivemarket/dto"
	openpositiondto "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/openposition/dto"
	getcurrentsignaldto "github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/getcurrentsignal/dto"
	placeorderdto "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/placeorder/dto"
	startwindowdto "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container/botcontainer"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/signing"
)

// assets are the four 5-minute markets the bot trades.
var assets = []string{"btc", "eth", "sol", "xrp"}

// confidenceMin is the minimum signal confidence required to enter a trade.
const confidenceMin = "0.30"

// orderSize is the fixed number of shares per order.
const orderSize = "15"

// openWindowsForAssets tries to open a trading window for each asset.
// It refreshes markets once then iterates assets. Errors are logged, not fatal.
func openWindowsForAssets(ctx context.Context, bc *botcontainer.BotContainer, signer *signing.Signer, funderAddr string, clobOrderIDs map[string]string) {
	if _, err := bc.RefreshMarkets.Execute(ctx, refreshmarketsdto.Input{}); err != nil {
		log.Printf("window: refresh markets: %v", err)
		return
	}

	for _, asset := range assets {
		if err := openWindowForAsset(ctx, bc, signer, funderAddr, asset, clobOrderIDs); err != nil {
			log.Printf("window: asset %s: %v", asset, err)
		}
	}
}

func openWindowForAsset(ctx context.Context, bc *botcontainer.BotContainer, signer *signing.Signer, funderAddr, asset string, clobOrderIDs map[string]string) error {
	// 1. Get active market
	mkt, err := bc.GetActiveMarket.Execute(ctx, getactivemarketdto.Input{Asset: asset})
	if err != nil {
		return fmt.Errorf("get active market: %w", err)
	}

	// 2. Get current signal
	sig, err := bc.GetCurrentSignal.Execute(ctx, getcurrentsignaldto.Input{Asset: asset})
	if err != nil {
		return fmt.Errorf("get signal: %w", err)
	}

	minConf := decimal.RequireFromString(confidenceMin)
	if sig.Signal.Confidence.LessThan(minConf) {
		log.Printf("window: %s: confidence %.2f < %.2f — skipping", asset, sig.Signal.Confidence.InexactFloat64(), minConf.InexactFloat64())
		return nil
	}

	// 3. Start window state
	if _, err = bc.StartWindow.Execute(ctx, startwindowdto.Input{
		Asset:       asset,
		MarketID:    mkt.MarketID,
		ConditionID: mkt.ConditionID,
		UpTokenID:   mkt.UpTokenID,
		DownTokenID: mkt.DownTokenID,
		TickSize:    mkt.TickSize,
		OpenPrice:   sig.Signal.OpenPrice,
	}); err != nil {
		return fmt.Errorf("start window: %w", err)
	}

	// 4. Choose token based on signal direction; side is always "buy"
	tokenID, outcome := chooseOrder(sig.Signal.Predicted, mkt.UpTokenID, mkt.DownTokenID)

	// 5. Fetch live fee rate
	feeRate, err := bc.FeeRateProvider.FetchFeeRate(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("fetch fee rate: %w", err)
	}

	// 6. Place order (computes EIP-712 hash, saves order locally)
	price := defaultEntryPrice(sig.Signal.Predicted)
	placed, err := bc.PlaceOrder.Execute(ctx, placeorderdto.Input{
		Asset:         asset,
		Outcome:       outcome,
		Side:          "buy",
		Price:         price,
		Size:          decimal.RequireFromString(orderSize),
		TokenID:       tokenID,
		FeeRateBps:    feeRate,
		FunderAddress: funderAddr,
	})
	if err != nil {
		return fmt.Errorf("place order: %w", err)
	}

	// 7. Sign the EIP-712 hash
	signature, err := signer.Sign(placed.UnsignedHash)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	// 8. Fetch the saved order domain object (needed by OrderSubmitter.Submit)
	o, err := bc.OrderRepository.FindByID(ctx, polyid.OrderID(placed.OrderID))
	if err != nil {
		return fmt.Errorf("fetch order: %w", err)
	}

	// 9. Submit to CLOB
	clobOrderID, err := bc.OrderSubmitter.Submit(ctx, o, signature)
	if err != nil {
		return fmt.Errorf("submit order: %w", err)
	}

	// 10. Track CLOB order ID in memory (needed for cancellation)
	clobOrderIDs[placed.OrderID] = clobOrderID
	log.Printf("window: %s: placed %s order — local=%s clob=%s", asset, outcome, placed.OrderID, clobOrderID)

	// 11. Open position for tracking
	if _, err = bc.OpenPosition.Execute(ctx, openpositiondto.Input{
		Asset:    asset,
		TokenID:  tokenID,
		Outcome:  outcome,
		Size:     orderSize,
		AvgPrice: price.String(),
		MarketID: mkt.MarketID,
	}); err != nil {
		log.Printf("window: %s: open position: %v (non-fatal)", asset, err)
	}

	return nil
}

// chooseOrder picks the outcome token based on the predicted direction.
// For "Up": buy the Up token. For "Down": buy the Down token.
// Side is always "buy" — the bot never shorts.
func chooseOrder(predicted, upTokenID, downTokenID string) (tokenID, outcome string) {
	if predicted == "Up" {
		return upTokenID, "Up"
	}
	return downTokenID, "Down"
}

// defaultEntryPrice returns the default limit price for a new order.
// Uses a conservative entry slightly away from 0.50 toward the predicted outcome.
func defaultEntryPrice(_ string) decimal.Decimal {
	return decimal.NewFromFloat(0.52)
}
