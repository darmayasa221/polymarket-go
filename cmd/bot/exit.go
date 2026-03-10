package main

import (
	"context"
	"log"
	"time"

	"github.com/shopspring/decimal"

	closepositiondto "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/closeposition/dto"
	getpositionsdto "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpositions/dto"
	getcurrentsignaldto "github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/getcurrentsignal/dto"
	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	cancelorderdto "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/cancelorder/dto"
	closewindowdto "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/closewindow/dto"
	getwindowstatedto "github.com/darmayasa221/polymarket-go/internal/applications/trading/queries/getwindowstate/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container/botcontainer"
)

// exitThreshold is how many dollars of oracle price move triggers stop loss / take profit.
var exitThreshold = decimal.NewFromFloat(0.20)

// timeStopBuffer is how long before window end the time stop triggers.
const timeStopBuffer = 30 * time.Second

// checkExits evaluates mid-window exit conditions for all assets with open windows.
func checkExits(ctx context.Context, bc *botcontainer.BotContainer, clobOrderIDs map[string]string) {
	now := timeutil.Now()
	for _, asset := range assets {
		if err := checkExitForAsset(ctx, bc, asset, now, clobOrderIDs); err != nil {
			log.Printf("exit: asset %s: %v", asset, err)
		}
	}
}

func checkExitForAsset(ctx context.Context, bc *botcontainer.BotContainer, asset string, now time.Time, clobOrderIDs map[string]string) error {
	stateOut, stateErr := bc.GetWindowState.Execute(ctx, getwindowstatedto.Input{Asset: asset})
	if stateErr != nil {
		return nil //nolint:nilerr // no active window for this asset is expected, not an error
	}

	if stateOut.Status != windowstate.WindowOpen {
		return nil
	}

	// Get current oracle signal to compare price movement
	sigOut, sigErr := bc.GetCurrentSignal.Execute(ctx, getcurrentsignaldto.Input{Asset: asset})
	if sigErr != nil {
		return nil //nolint:nilerr // no signal yet — skip exit check, never exit blindly
	}

	openPrice := stateOut.OpenPrice
	currentPrice := sigOut.Signal.CurrentPrice
	delta := currentPrice.Sub(openPrice).Abs()

	timeStop := now.After(stateOut.WindowEnd.Add(-timeStopBuffer))
	stopLoss := delta.GreaterThanOrEqual(exitThreshold) && currentPrice.LessThan(openPrice)
	takeProfit := delta.GreaterThanOrEqual(exitThreshold) && currentPrice.GreaterThan(openPrice)

	if !timeStop && !stopLoss && !takeProfit {
		return nil
	}

	reason := exitReason(timeStop, stopLoss, takeProfit)
	log.Printf("exit: %s triggered (%s) — openPrice=%s currentPrice=%s", asset, reason, openPrice, currentPrice)

	return triggerExit(ctx, bc, stateOut.MarketID, asset, clobOrderIDs, currentPrice)
}

func triggerExit(ctx context.Context, bc *botcontainer.BotContainer, marketID, asset string, clobOrderIDs map[string]string, currentPrice decimal.Decimal) error {
	// Cancel all open CLOB orders (best-effort)
	openOrders, _ := bc.OrderRepository.ListOpenByMarket(ctx, marketID)
	for _, o := range openOrders {
		localID := o.ID().String()
		clobID, ok := clobOrderIDs[localID]
		if !ok {
			continue
		}
		if _, err := bc.CancelOrder.Execute(ctx, cancelorderdto.Input{
			OrderID:     localID,
			ClobOrderID: clobID,
		}); err != nil {
			log.Printf("exit: cancel order %s: %v (non-fatal)", localID, err)
		}
		delete(clobOrderIDs, localID)
	}

	// Close the window state in Redis
	if _, err := bc.CloseWindow.Execute(ctx, closewindowdto.Input{Asset: asset}); err != nil {
		log.Printf("exit: close window %s: %v", asset, err)
	}

	// Close any open positions for this asset at current oracle price
	posOut, posErr := bc.GetPositions.Execute(ctx, getpositionsdto.Input{})
	if posErr != nil {
		return nil //nolint:nilerr // best-effort: skip position close if query fails
	}
	for _, pos := range posOut.Positions {
		if pos.Asset != asset {
			continue
		}
		if _, err := bc.ClosePosition.Execute(ctx, closepositiondto.Input{
			PositionID: pos.PositionID,
			ExitPrice:  currentPrice.String(),
		}); err != nil {
			log.Printf("exit: close position %s: %v (non-fatal)", pos.PositionID, err)
		}
	}

	return nil
}

func exitReason(timeStop, stopLoss, takeProfit bool) string {
	switch {
	case timeStop:
		return "time-stop"
	case stopLoss:
		return "stop-loss"
	case takeProfit:
		return "take-profit"
	default:
		return "unknown"
	}
}
