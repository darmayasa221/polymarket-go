package main

import (
	"context"
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

// run starts the bot and keeps it running by reconnecting on WebSocket drops.
// It only returns when ctx is canceled (clean shutdown).
func (r *runner) run(ctx context.Context) {
	b := newWsBackoff()
	for {
		// Fresh connected signal for this attempt.
		// attempt() calls onPrice() which sends to r.connected on first Chainlink price.
		r.connected = make(chan struct{}, 1)
		r.attempt(ctx)

		if ctx.Err() != nil {
			log.Println("runner: shutting down")
			return
		}

		// If at least one Chainlink price was received, the connection was healthy.
		// Reset backoff so the next drop starts at 1s, not wherever we left off.
		select {
		case <-r.connected:
			b.reset()
		default:
			// No price received — connection failed before delivering data. Keep backoff.
		}

		wait := b.next()
		log.Printf("runner: disconnected — reconnecting in %s", wait)

		// Reset per-connection order state.
		// priceBuffer is intentionally preserved — rebuilding takes 15+ minutes.
		r.clobOrderIDs = make(map[string]string)
		r.lastWindow = 0

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			log.Println("runner: shutting down during reconnect wait")
			return
		}
	}
}

// attempt dials all three WebSocket connections and runs the event loop until
// any connection drops or ctx is canceled. Dial failures are logged and treated
// as disconnects so the outer retry loop handles them.
//
// WS handlers are stateless beyond their URL — calling Start() multiple times is
// safe; each call dials a fresh independent connection. Goroutines from the
// previous attempt have already exited before attempt() is called again: a WS
// error causes the handler goroutine to close its output channel, eventLoop sees
// the closed channel and returns, and only then does attempt() return and run()
// schedule the next call. The same ctx propagates through, so a ctx cancellation
// also cleanly exits all handler goroutines via conn.Close().
func (r *runner) attempt(ctx context.Context) {
	attemptCtx, cancelAttempt := context.WithCancel(ctx)
	defer cancelAttempt() // ensure all WS goroutines exit when attempt returns

	if _, err := r.bc.RefreshMarkets.Execute(attemptCtx, struct{}{}); err != nil {
		log.Printf("runner: refresh markets: %v (continuing)", err)
	}

	priceCh, err := r.bc.RTDSHandler.Start(attemptCtx)
	if err != nil {
		log.Printf("runner: start RTDS WS: %v — will retry", err)
		return
	}

	conditionIDs := r.allConditionIDs(attemptCtx)
	var marketCh <-chan market.MarketEvent
	if len(conditionIDs) > 0 {
		marketCh, err = r.bc.MarketHandler.Start(attemptCtx, conditionIDs)
		if err != nil {
			log.Printf("runner: start market WS: %v — will retry", err)
			return
		}
	} else {
		log.Println("runner: no active markets yet — market WS skipped until next refresh")
		marketCh = make(chan market.MarketEvent)
	}

	userCh, err := r.bc.UserHandler.Start(attemptCtx, r.clobCfg)
	if err != nil {
		log.Printf("runner: start user WS: %v — will retry", err)
		return
	}

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	exitTicker := time.NewTicker(exitMonitorInterval)
	windowTicker := time.NewTicker(windowCheckInterval)
	defer heartbeatTicker.Stop()
	defer exitTicker.Stop()
	defer windowTicker.Stop()

	log.Println("runner: connected — starting event loop")
	r.eventLoop(attemptCtx, priceCh, marketCh, userCh, heartbeatTicker, exitTicker, windowTicker)
	log.Println("runner: event loop exited")
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
// The first Chainlink price signals a healthy connection to the reconnect loop.
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
		// Signal first healthy price to run() (non-blocking; buffered channel holds at most 1).
		select {
		case r.connected <- struct{}{}:
		default:
		}
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
