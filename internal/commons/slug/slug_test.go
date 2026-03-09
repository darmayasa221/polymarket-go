// Package slug_test tests the slug builder.
package slug_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/commons/slug"
)

func TestForAsset(t *testing.T) {
	t.Parallel()

	windowStart := time.Unix(1_700_000_100, 0).UTC()

	tests := []struct {
		asset string
		want  string
	}{
		{asset: "btc", want: "btc-updown-5m-1700000100"},
		{asset: "eth", want: "eth-updown-5m-1700000100"},
		{asset: "sol", want: "sol-updown-5m-1700000100"},
		{asset: "xrp", want: "xrp-updown-5m-1700000100"},
	}

	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			t.Parallel()
			got := slug.ForAsset(tt.asset, windowStart)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestCurrentWindow(t *testing.T) {
	t.Parallel()

	got := slug.CurrentWindow("btc")
	assert.True(t, got.String() != "")
	assert.Contains(t, got.String(), "btc-updown-5m-")
}

func TestNextWindow(t *testing.T) {
	t.Parallel()

	cur := slug.CurrentWindow("btc")
	nxt := slug.NextWindow("btc")
	assert.NotEqual(t, cur.String(), nxt.String())
}
