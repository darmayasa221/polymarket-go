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

// Capital management constants.
const (
	reserveRatio   = "0.20"    // always keep 20% of balance untouched
	minOrderTokens = int64(5)  // Polymarket minimum order size
	maxOrderTokens = int64(50) // cap per position to avoid over-concentration
)

// candidate holds a pre-fetched signal for an asset that passed the confidence check.
type candidate struct {
	asset  string
	signal getcurrentsignaldto.Output
}

// openWindowsForAssets fetches balance, collects confident signals, calculates
// position size, then opens windows for all qualifying assets.
func openWindowsForAssets(ctx context.Context, bc *botcontainer.BotContainer, signer *signing.Signer, funderAddr string, clobOrderIDs map[string]string) {
	if _, err := bc.RefreshMarkets.Execute(ctx, refreshmarketsdto.Input{}); err != nil {
		log.Printf("window: refresh markets: %v", err)
		return
	}

	// Fetch live USDC.e balance — this drives position sizing.
	balance, err := bc.BalanceProvider.FetchBalance(ctx)
	if err != nil {
		log.Printf("window: fetch balance: %v — skipping window", err)
		return
	}
	log.Printf("window: balance=%.2f USDC", balance.InexactFloat64())

	// Pass 1 — collect assets whose signal confidence meets the threshold.
	minConf := decimal.RequireFromString(confidenceMin)
	candidates := make([]candidate, 0, len(assets))
	for _, asset := range assets {
		sig, sigErr := bc.GetCurrentSignal.Execute(ctx, getcurrentsignaldto.Input{Asset: asset})
		if sigErr != nil {
			log.Printf("window: %s: signal error: %v", asset, sigErr)
			continue
		}
		if sig.Signal.Confidence.LessThan(minConf) {
			log.Printf("window: %s: confidence %.2f < %.2f — skipping", asset,
				sig.Signal.Confidence.InexactFloat64(), minConf.InexactFloat64())
			continue
		}
		candidates = append(candidates, candidate{asset: asset, signal: sig})
	}

	if len(candidates) == 0 {
		log.Printf("window: no confident signals — sitting out this window")
		return
	}

	// Calculate how many tokens to buy per asset.
	// Deployable capital is split evenly across all confident assets.
	tokenCount := calcTokenCount(balance, len(candidates))
	if tokenCount == 0 {
		log.Printf("window: balance %.2f too low to afford minimum order (%d tokens) across %d asset(s)",
			balance.InexactFloat64(), minOrderTokens, len(candidates))
		return
	}

	log.Printf("window: %d confident signal(s) — deploying %d tokens each (balance=%.2f reserve=20%%)",
		len(candidates), tokenCount, balance.InexactFloat64())

	// Pass 2 — open a window for each qualifying asset.
	for _, c := range candidates {
		if err := openWindowForAsset(ctx, bc, signer, funderAddr, c.asset, c.signal, tokenCount, clobOrderIDs); err != nil {
			log.Printf("window: asset %s: %v", c.asset, err)
		}
	}
}

// calcTokenCount returns the number of tokens to buy per asset.
// Spreads 80% of balance evenly across numAssets; clamped to [minOrderTokens, maxOrderTokens].
// Returns 0 if the per-asset budget is below the minimum order size.
func calcTokenCount(balance decimal.Decimal, numAssets int) int64 {
	reserve := decimal.RequireFromString(reserveRatio)
	deployable := balance.Mul(decimal.NewFromInt(1).Sub(reserve))
	perAsset := deployable.Div(decimal.NewFromInt(int64(numAssets)))
	tokens := perAsset.Div(defaultEntryPrice()).Floor().IntPart()

	if tokens < minOrderTokens {
		return 0
	}
	if tokens > maxOrderTokens {
		return maxOrderTokens
	}
	return tokens
}

func openWindowForAsset(
	ctx context.Context,
	bc *botcontainer.BotContainer,
	signer *signing.Signer,
	funderAddr, asset string,
	sig getcurrentsignaldto.Output,
	tokenCount int64,
	clobOrderIDs map[string]string,
) error {
	// 1. Get active market
	mkt, err := bc.GetActiveMarket.Execute(ctx, getactivemarketdto.Input{Asset: asset})
	if err != nil {
		return fmt.Errorf("get active market: %w", err)
	}

	// 2. Start window state
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

	// 3. Choose token based on signal direction; side is always "buy"
	tokenID, outcome := chooseOrder(sig.Signal.Predicted, mkt.UpTokenID, mkt.DownTokenID)

	// 4. Fetch live fee rate
	feeRate, err := bc.FeeRateProvider.FetchFeeRate(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("fetch fee rate: %w", err)
	}

	// 5. Place order (computes EIP-712 hash, saves order locally)
	price := defaultEntryPrice()
	placed, err := bc.PlaceOrder.Execute(ctx, placeorderdto.Input{
		Asset:         asset,
		Outcome:       outcome,
		Side:          "buy",
		Price:         price,
		Size:          decimal.NewFromInt(tokenCount),
		TokenID:       tokenID,
		FeeRateBps:    feeRate,
		FunderAddress: funderAddr,
	})
	if err != nil {
		return fmt.Errorf("place order: %w", err)
	}

	// 6. Sign the EIP-712 hash
	signature, err := signer.Sign(placed.UnsignedHash)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	// 7. Fetch the saved order domain object (needed by OrderSubmitter.Submit)
	o, err := bc.OrderRepository.FindByID(ctx, polyid.OrderID(placed.OrderID))
	if err != nil {
		return fmt.Errorf("fetch order: %w", err)
	}

	// 8. Submit to CLOB
	clobOrderID, err := bc.OrderSubmitter.Submit(ctx, o, signature)
	if err != nil {
		return fmt.Errorf("submit order: %w", err)
	}

	// 9. Track CLOB order ID in memory (needed for cancellation)
	clobOrderIDs[placed.OrderID] = clobOrderID
	log.Printf("window: %s: placed %s order — tokens=%d price=%s local=%s clob=%s",
		asset, outcome, tokenCount, price.String(), placed.OrderID, clobOrderID)

	// 10. Open position for tracking
	if _, err = bc.OpenPosition.Execute(ctx, openpositiondto.Input{
		Asset:    asset,
		TokenID:  tokenID,
		Outcome:  outcome,
		Size:     fmt.Sprintf("%d", tokenCount),
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

// defaultEntryPrice returns the conservative limit price for a new order.
// Buying at $0.52 means only filling when the market agrees with the signal direction.
func defaultEntryPrice() decimal.Decimal {
	return decimal.NewFromFloat(0.52)
}
