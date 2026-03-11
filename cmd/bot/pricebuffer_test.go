package main

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPriceBuffer_NoData(t *testing.T) {
	buf := newPriceBuffer(3)
	dir, conf := buf.momentum("btc", 3)
	assert.Equal(t, "", dir)
	assert.True(t, conf.IsZero())
}

func TestPriceBuffer_InsufficientHistory(t *testing.T) {
	buf := newPriceBuffer(3)
	t0 := time.Unix(1741996800, 0) // window boundary
	buf.update("btc", decimal.NewFromFloat(64000), t0)
	// Only 1 price in first window — no closes yet, momentum should be empty.
	dir, conf := buf.momentum("btc", 2)
	assert.Equal(t, "", dir)
	assert.True(t, conf.IsZero())
}

func TestPriceBuffer_ThreeWindowsUp(t *testing.T) {
	// newPriceBuffer(3) + 4 time points → 3 closes.
	// With n=3, startIdx=0, idx=0 is skipped → 2 comparisons, both up → conf=2/2=1.0.
	buf := newPriceBuffer(3)
	t0 := time.Unix(1741996800, 0)
	t1 := time.Unix(1741996800+300, 0)
	t2 := time.Unix(1741996800+600, 0)
	t3 := time.Unix(1741996800+900, 0)

	buf.update("btc", decimal.NewFromFloat(64000), t0)
	buf.update("btc", decimal.NewFromFloat(64100), t1) // closes w0 at 64000
	buf.update("btc", decimal.NewFromFloat(64250), t2) // closes w1 at 64100
	buf.update("btc", decimal.NewFromFloat(64400), t3) // closes w2 at 64250

	dir, conf := buf.momentum("btc", 3)
	assert.Equal(t, "Up", dir)
	// 2 comparisons (idx 1 vs 0, idx 2 vs 1), both Up → conf = 2/2 = 1.0
	assert.Equal(t, "1", conf.String())
}

func TestPriceBuffer_TwoOfThreeUp(t *testing.T) {
	// newPriceBuffer(4) + 5 time points → 4 closes.
	// With n=3, startIdx=1, all 3 comparisons valid → conf=2/3.
	buf := newPriceBuffer(4)
	t0 := time.Unix(1741996800, 0)
	t1 := time.Unix(1741996800+300, 0)
	t2 := time.Unix(1741996800+600, 0)
	t3 := time.Unix(1741996800+900, 0)
	t4 := time.Unix(1741996800+1200, 0)

	buf.update("btc", decimal.NewFromFloat(64000), t0) // set current w0
	buf.update("btc", decimal.NewFromFloat(63900), t1) // closes w0=64000; w1 current=63900 (down)
	buf.update("btc", decimal.NewFromFloat(64100), t2) // closes w1=63900; w2 current=64100 (up)
	buf.update("btc", decimal.NewFromFloat(64300), t3) // closes w2=64100; w3 current=64300 (up)
	buf.update("btc", decimal.NewFromFloat(64500), t4) // closes w3=64300; w4 current=64500

	// closes = [64000, 63900, 64100, 64300]; n=3, startIdx=1
	// idx1: 63900<64000 → down; idx2: 64100>63900 → up; idx3: 64300>64100 → up
	// 2 up, 1 down → dir=Up, conf=2/3
	dir, conf := buf.momentum("btc", 3)
	assert.Equal(t, "Up", dir)
	expected := decimal.NewFromInt(2).Div(decimal.NewFromInt(3))
	assert.True(t, conf.Equal(expected), "expected %s got %s", expected, conf)
}

func TestPriceBuffer_RingBufferCapacity(t *testing.T) {
	buf := newPriceBuffer(3) // max 3 closes
	base := int64(1741996800)
	price := 64000.0
	// Feed 6 windows — only last 3 closes should be retained.
	for i := 0; i < 6; i++ {
		buf.update("btc", decimal.NewFromFloat(price+float64(i*100)), time.Unix(base+int64(i)*300, 0))
	}
	// 3 closes retained, all up, 2 comparisons → conf=1.0
	dir, conf := buf.momentum("btc", 3)
	assert.Equal(t, "Up", dir)
	assert.Equal(t, "1", conf.String())
}

func TestPriceBuffer_MultipleAssets(t *testing.T) {
	// Use newPriceBuffer(4) with 5 time points for proper 3-comparison momentum.
	buf := newPriceBuffer(4)
	t0 := time.Unix(1741996800, 0)
	t1 := time.Unix(1741996800+300, 0)
	t2 := time.Unix(1741996800+600, 0)
	t3 := time.Unix(1741996800+900, 0)
	t4 := time.Unix(1741996800+1200, 0)

	// BTC goes up across all windows.
	buf.update("btc", decimal.NewFromFloat(64000), t0)
	buf.update("btc", decimal.NewFromFloat(64100), t1)
	buf.update("btc", decimal.NewFromFloat(64200), t2)
	buf.update("btc", decimal.NewFromFloat(64300), t3)
	buf.update("btc", decimal.NewFromFloat(64400), t4)

	// ETH goes down across all windows.
	buf.update("eth", decimal.NewFromFloat(3000), t0)
	buf.update("eth", decimal.NewFromFloat(2950), t1)
	buf.update("eth", decimal.NewFromFloat(2900), t2)
	buf.update("eth", decimal.NewFromFloat(2850), t3)
	buf.update("eth", decimal.NewFromFloat(2800), t4)

	btcDir, _ := buf.momentum("btc", 3)
	ethDir, _ := buf.momentum("eth", 3)
	assert.Equal(t, "Up", btcDir)
	assert.Equal(t, "Down", ethDir)
}
