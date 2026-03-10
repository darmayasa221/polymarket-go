package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/updateticksize/dto"
	getactivemarketdto "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/getactivemarket/dto"
	recordpricedto "github.com/darmayasa221/polymarket-go/internal/applications/pricing/commands/recordprice/dto"
	heartbeatdto "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/heartbeat/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container/botcontainer"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/market"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/user"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/signing"
)

const (
	heartbeatInterval   = 5 * time.Second
	exitMonitorInterval = 30 * time.Second
	windowCheckInterval = 30 * time.Second // check for new windows every 30s
)

// runner orchestrates the bot's event loop.
type runner struct {
	bc           *botcontainer.BotContainer
	signer       *signing.Signer
	clobCfg      clob.Config
	clobOrderIDs map[string]string // localOrderID → clobOrderID
	lastWindow   int64             // Unix timestamp of last opened window boundary
}

func newRunner(bc *botcontainer.BotContainer, signer *signing.Signer, cfg clob.Config) *runner {
	return &runner{
		bc:           bc,
		signer:       signer,
		clobCfg:      cfg,
		clobOrderIDs: make(map[string]string),
	}
}

// run starts WS handlers, tickers, and delegates to eventLoop.
func (r *runner) run(ctx context.Context) error {
	// Refresh markets before subscribing to ensure condition IDs are available
	if _, err := r.bc.RefreshMarkets.Execute(ctx, struct{}{}); err != nil {
		log.Printf("runner: initial refresh markets: %v (continuing)", err)
	}

	priceCh, err := r.bc.RTDSHandler.Start(ctx)
	if err != nil {
		return fmt.Errorf("runner: start RTDS WS: %w", err)
	}

	conditionIDs := r.allConditionIDs(ctx)
	marketCh, err := r.bc.MarketHandler.Start(ctx, conditionIDs)
	if err != nil {
		return fmt.Errorf("runner: start market WS: %w", err)
	}

	userCh, err := r.bc.UserHandler.Start(ctx, r.clobCfg)
	if err != nil {
		return fmt.Errorf("runner: start user WS: %w", err)
	}

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	exitTicker := time.NewTicker(exitMonitorInterval)
	windowTicker := time.NewTicker(windowCheckInterval)
	defer heartbeatTicker.Stop()
	defer exitTicker.Stop()
	defer windowTicker.Stop()

	log.Println("runner: started")
	r.eventLoop(ctx, priceCh, marketCh, userCh, heartbeatTicker, exitTicker, windowTicker)
	return nil
}

// eventLoop processes events from all channels and tickers.
// The multi-channel select is the canonical Go event dispatcher pattern; its complexity is structural.
//
//nolint:gocognit // 7-channel select dispatcher — complexity is inherent to the architecture
func (r *runner) eventLoop(ctx context.Context, priceCh <-chan *oracle.Price, marketCh <-chan market.MarketEvent, userCh <-chan user.UserEvent, heartbeatTicker, exitTicker, windowTicker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			log.Println("runner: shutting down")
			return
		case price, ok := <-priceCh:
			if !ok {
				return
			}
			r.onPrice(ctx, price)
		case ev, ok := <-marketCh:
			if !ok {
				return
			}
			r.onMarketEvent(ctx, ev)
		case ev, ok := <-userCh:
			if !ok {
				return
			}
			log.Printf("runner: user event type=%s orderID=%s", ev.Type, ev.OrderID)
		case <-heartbeatTicker.C:
			if _, err := r.bc.Heartbeat.Execute(ctx, heartbeatdto.Input{}); err != nil {
				log.Printf("runner: heartbeat: %v", err)
			}
		case <-windowTicker.C:
			boundary := timeutil.WindowStart(timeutil.Now()).Unix()
			if boundary != r.lastWindow {
				r.lastWindow = boundary
				log.Printf("runner: new window boundary %d — opening windows", boundary)
				openWindowsForAssets(ctx, r.bc, r.signer, r.clobCfg.Address, r.clobOrderIDs)
			}
		case <-exitTicker.C:
			checkExits(ctx, r.bc, r.clobOrderIDs)
		}
	}
}

// onPrice records an oracle price reading from the RTDS WebSocket.
func (r *runner) onPrice(ctx context.Context, price *oracle.Price) {
	_, _ = r.bc.RecordPrice.Execute(ctx, recordpricedto.Input{
		Asset:      price.Asset(),
		Source:     string(price.Source()),
		Value:      price.Value().String(),
		RoundedAt:  price.RoundedAt(),
		ReceivedAt: price.ReceivedAt(),
	})
}

// onMarketEvent handles market WS events, currently only tick_size_change.
func (r *runner) onMarketEvent(ctx context.Context, ev market.MarketEvent) {
	if ev.TickSizeChange == nil {
		return
	}
	_, _ = r.bc.UpdateTickSize.Execute(ctx, dto.Input{
		ConditionID: ev.TickSizeChange.ConditionID,
		NewTickSize: ev.TickSizeChange.NewTickSize,
	})
}

// allConditionIDs fetches condition IDs for all active markets to subscribe the market WS.
func (r *runner) allConditionIDs(ctx context.Context) []string {
	ids := make([]string, 0, len(assets))
	for _, asset := range assets {
		out, err := r.bc.GetActiveMarket.Execute(ctx, getactivemarketdto.Input{Asset: asset})
		if err != nil {
			continue
		}
		ids = append(ids, out.ConditionID)
	}
	return ids
}
