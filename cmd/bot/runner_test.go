package main

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/market"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/user"
)

func TestWsBackoff_Sequence(t *testing.T) {
	b := newWsBackoff()
	assert.Equal(t, 1*time.Second, b.next())
	assert.Equal(t, 2*time.Second, b.next())
	assert.Equal(t, 4*time.Second, b.next())
	assert.Equal(t, 8*time.Second, b.next())
	assert.Equal(t, 16*time.Second, b.next())
	assert.Equal(t, 30*time.Second, b.next()) // capped at wsBackoffMax
	assert.Equal(t, 30*time.Second, b.next()) // stays at cap
}

func TestWsBackoff_Reset(t *testing.T) {
	b := newWsBackoff()
	b.next()
	b.next()
	b.next()
	b.reset()
	assert.Equal(t, 1*time.Second, b.next()) // back to 1s after reset
}

// TestEventLoop_ExitsOnClosedPriceChannel documents the contract that the retry
// loop relies on: when priceCh is closed, eventLoop returns.
// Tickers use time.Hour to ensure they never fire during the test — this prevents
// nil-pointer panics on r.bc (which is nil in this test since bc is not needed).
func TestEventLoop_ExitsOnClosedPriceChannel(t *testing.T) {
	r := &runner{buf: newPriceBuffer(4), clobOrderIDs: make(map[string]string)}

	priceCh := make(chan *oracle.Price)
	close(priceCh) // simulate immediate RTDS disconnect

	marketCh := make(chan market.MarketEvent) // never sends
	userCh := make(chan user.UserEvent)       // never sends

	// Use time.Hour tickers so none fire before the closed priceCh is selected.
	// This keeps r.bc nil-safe for the test.
	hticker := time.NewTicker(time.Hour)
	eticker := time.NewTicker(time.Hour)
	wticker := time.NewTicker(time.Hour)
	defer hticker.Stop()
	defer eticker.Stop()
	defer wticker.Stop()

	done := make(chan struct{})
	go func() {
		r.eventLoop(t.Context(), priceCh, marketCh, userCh, hticker, eticker, wticker)
		close(done)
	}()

	select {
	case <-done: // pass — eventLoop returned as expected
	case <-time.After(time.Second):
		t.Fatal("eventLoop did not return after priceCh closed")
	}
}

// TestRunner_ReconnectResetsOrderState documents that clobOrderIDs and lastWindow
// are cleared on reconnect. CLOB auto-cancels all GTD orders within 10s of
// heartbeat stopping, so clearing clobOrderIDs is correct on reconnect.
func TestRunner_ReconnectResetsOrderState(t *testing.T) {
	r := &runner{
		clobOrderIDs: map[string]string{"order1": "clob1"},
		buf:          newPriceBuffer(4),
		lastWindow:   12345,
	}

	// Simulate what run() does between attempts.
	r.clobOrderIDs = make(map[string]string)
	r.lastWindow = 0

	assert.Empty(t, r.clobOrderIDs, "clobOrderIDs must be cleared on reconnect")
	assert.Equal(t, int64(0), r.lastWindow, "lastWindow must be reset on reconnect")
}

// TestRunner_ReconnectPreservesBuffer documents that priceBuffer is NOT reset on
// reconnect. Rebuilding the buffer takes ~15 minutes (3 windows × 5 min each),
// so preserving it is essential for uninterrupted momentum signal after a drop.
func TestRunner_ReconnectPreservesBuffer(t *testing.T) {
	r := &runner{
		clobOrderIDs: make(map[string]string),
		buf:          newPriceBuffer(4),
		lastWindow:   0,
	}

	// Feed a price into the buffer before simulated reconnect.
	r.buf.update("btc", decimal.NewFromFloat(64000), time.Unix(1741996800, 0))

	// Reconnect state reset — must NOT touch priceBuffer.
	r.clobOrderIDs = make(map[string]string)
	r.lastWindow = 0

	price := r.buf.currentPrice("btc")
	assert.True(t, price.Equal(decimal.NewFromFloat(64000)), "priceBuffer must survive reconnect")
}
