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

const wsBackoffMax = 30 * time.Second

// wsBackoff implements capped exponential backoff for WebSocket reconnects.
// Sequence: 1s → 2s → 4s → 8s → 16s → 30s (cap).
// Call reset() when the first successful RTDS price is received to restart at 1s.
type wsBackoff struct{ current time.Duration }

func newWsBackoff() *wsBackoff { return &wsBackoff{current: time.Second} }

func (b *wsBackoff) next() time.Duration {
	d := b.current
	b.current *= 2
	if b.current > wsBackoffMax {
		b.current = wsBackoffMax
	}
	return d
}

func (b *wsBackoff) reset() { b.current = time.Second }

// runner orchestrates the bot's event loop.
type runner struct {
	bc           *botcontainer.BotContainer
	signer       *signing.Signer
	clobCfg      clob.Config
	clobOrderIDs map[string]string // localOrderID → clobOrderID
	lastWindow   int64             // Unix timestamp of last opened window boundary
	buf          *priceBuffer      // in-memory cross-window Chainlink price history
	connected    chan struct{}     // closed by onPrice on first Chainlink price; signals healthy connection to run()
}

func newRunner(bc *botcontainer.BotContainer, signer *signing.Signer, cfg clob.Config) *runner {
	return &runner{
		bc:           bc,
		signer:       signer,
		clobCfg:      cfg,
		clobOrderIDs: make(map[string]string),
		buf:          newPriceBuffer(4), // keeps 4 closes → 3 valid comparisons for momentum(n=3)
		connected:    make(chan struct{}, 1),
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
	var marketCh <-chan market.MarketEvent
	if len(conditionIDs) > 0 {
		marketCh, err = r.bc.MarketHandler.Start(ctx, conditionIDs)
		if err != nil {
			return fmt.Errorf("runner: start market WS: %w", err)
		}
	} else {
		log.Println("runner: no active markets yet — market WS skipped until next refresh")
		marketCh = make(chan market.MarketEvent) // never closes, never sends
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
				log.Println("runner: RTDS price channel closed — exiting")
				return
			}
			r.onPrice(ctx, price)
		case ev, ok := <-marketCh:
			if !ok {
				log.Println("runner: market WS channel closed — exiting")
				return
			}
			r.onMarketEvent(ctx, ev)
		case ev, ok := <-userCh:
			if !ok {
				log.Println("runner: user WS channel closed — exiting")
				return
			}
			log.Printf("runner: user event type=%s order_type=%s status=%s orderID=%s", ev.EventType, ev.OrderType, ev.Status, ev.OrderID)
		case <-heartbeatTicker.C:
			if _, err := r.bc.Heartbeat.Execute(ctx, heartbeatdto.Input{}); err != nil {
				log.Printf("runner: heartbeat: %v", err)
			}
		case <-windowTicker.C:
			boundary := timeutil.WindowStart(timeutil.Now()).Unix()
			if boundary != r.lastWindow {
				r.lastWindow = boundary
				log.Printf("runner: window ticker boundary=%d — opening windows (fallback)", boundary)
				openWindowsForAssets(ctx, r.bc, r.signer, r.clobCfg.Address, r.clobOrderIDs, r.buf)
			}
		case <-exitTicker.C:
			checkExits(ctx, r.bc, r.clobOrderIDs)
		}
	}
}

// onPrice records an oracle price reading from the RTDS WebSocket.
// Chainlink prices are also fed into the cross-window momentum buffer.
func (r *runner) onPrice(ctx context.Context, price *oracle.Price) {
	_, _ = r.bc.RecordPrice.Execute(ctx, recordpricedto.Input{
		Asset:      price.Asset(),
		Source:     string(price.Source()),
		Value:      price.Value().String(),
		RoundedAt:  price.RoundedAt(),
		ReceivedAt: price.ReceivedAt(),
	})
	if price.Source() == oracle.SourceChainlink {
		r.buf.update(price.Asset(), price.Value(), price.ReceivedAt())
	}
}

// onMarketEvent handles market WS events: tick_size_change and new_market.
func (r *runner) onMarketEvent(ctx context.Context, ev market.MarketEvent) {
	switch ev.Type {
	case market.EventTickSizeChange:
		if ev.TickSizeChange == nil {
			return
		}
		_, _ = r.bc.UpdateTickSize.Execute(ctx, dto.Input{
			ConditionID: ev.TickSizeChange.ConditionID,
			NewTickSize: ev.TickSizeChange.NewTickSize,
		})
	case market.EventNewMarket:
		boundary := timeutil.WindowStart(timeutil.Now()).Unix()
		if boundary != r.lastWindow {
			r.lastWindow = boundary
			log.Printf("runner: new_market event boundary=%d — opening windows immediately", boundary)
			openWindowsForAssets(ctx, r.bc, r.signer, r.clobCfg.Address, r.clobOrderIDs, r.buf)
		}
	case market.EventMarketResolved:
		// Market resolved — no action needed; positions are tracked separately.
	}
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
