package main

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// windowBoundary returns the start timestamp of the 5-minute window containing t.
// Formula: floor(unix/300)*300 — matches the slug pattern.
func windowBoundary(t time.Time) int64 {
	unix := t.Unix()
	return unix - (unix % 300)
}

// windowClose records a single window's close Chainlink price.
type windowClose struct {
	boundary int64           // Unix timestamp of window start
	price    decimal.Decimal // last Chainlink price received in this window
}

// priceBuffer maintains an in-memory ring buffer of per-asset Chainlink window closes.
// It is safe for concurrent use; update is called from the RTDS goroutine.
type priceBuffer struct {
	mu      sync.RWMutex
	closes  map[string][]windowClose // asset → ring buffer of closes
	current map[string]windowClose   // asset → live (in-progress) window
	maxN    int
}

// newPriceBuffer creates a priceBuffer that retains the last maxN window closes.
func newPriceBuffer(maxN int) *priceBuffer {
	return &priceBuffer{
		closes:  make(map[string][]windowClose),
		current: make(map[string]windowClose),
		maxN:    maxN,
	}
}

// update records a live Chainlink price for asset at time t.
// When t crosses a new window boundary, the previous window's price is snapshotted.
func (b *priceBuffer) update(asset string, price decimal.Decimal, t time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	boundary := windowBoundary(t)
	cur, hasCur := b.current[asset]

	if hasCur && boundary > cur.boundary {
		// Crossed into a new window — snapshot the previous close.
		b.appendClose(asset, cur)
	}

	b.current[asset] = windowClose{boundary: boundary, price: price}
}

// appendClose appends a close to the ring buffer, evicting the oldest if full.
func (b *priceBuffer) appendClose(asset string, wc windowClose) {
	buf := b.closes[asset]
	buf = append(buf, wc)
	if len(buf) > b.maxN {
		buf = buf[len(buf)-b.maxN:]
	}
	b.closes[asset] = buf
}

// momentum returns the dominant direction and confidence across the last n window closes.
// Returns ("", zero) when insufficient history is available.
// Confidence = (majority count) / (total comparisons).
// Each close is compared to its predecessor; the first close in the window is skipped if
// no prior close exists.
func (b *priceBuffer) momentum(asset string, n int) (string, decimal.Decimal) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	closes := b.closes[asset]
	if len(closes) < n {
		return "", decimal.Zero
	}

	allCloses := closes
	startIdx := len(allCloses) - n

	ups, downs := 0, 0
	for i := 0; i < n; i++ {
		idx := startIdx + i
		if idx == 0 {
			// No prior close to compare against — skip.
			continue
		}
		if allCloses[idx].price.GreaterThanOrEqual(allCloses[idx-1].price) {
			ups++
		} else {
			downs++
		}
	}

	total := ups + downs
	if total == 0 {
		return "", decimal.Zero
	}
	if ups >= downs {
		return "Up", decimal.NewFromInt(int64(ups)).Div(decimal.NewFromInt(int64(total)))
	}
	return "Down", decimal.NewFromInt(int64(downs)).Div(decimal.NewFromInt(int64(total)))
}
