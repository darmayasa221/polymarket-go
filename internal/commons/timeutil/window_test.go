package timeutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
)

func TestWindowStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input time.Time
		want  int64 // unix seconds of expected window start
	}{
		{
			name:  "floors to 5-minute boundary at exact boundary",
			input: time.Unix(1_700_000_000, 0).UTC(), // 1700000000 % 300 == 200 → floor
			want:  1_699_999_800,                     // floor(1700000000 / 300) * 300
		},
		{
			name:  "floors to 5-minute boundary mid-window",
			input: time.Unix(1_700_000_150, 0).UTC(), // 150s into window
			want:  1_700_000_100,                     // floor(1700000150 / 300) * 300 = 5666667*300 = 1700000100? Let's verify
			// 1700000150 / 300 = 5666667.16... → floor = 5666667 → 5666667*300 = 1700000100
		},
		{
			name:  "at window start returns same boundary",
			input: time.Unix(1_700_000_100, 0).UTC(),
			want:  1_700_000_100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := timeutil.WindowStart(tt.input)
			assert.Equal(t, tt.want, got.Unix())
			assert.Equal(t, time.UTC, got.Location())
		})
	}
}

func TestWindowEnd(t *testing.T) {
	t.Parallel()

	input := time.Unix(1_700_000_150, 0).UTC()
	start := timeutil.WindowStart(input)
	end := timeutil.WindowEnd(input)

	assert.Equal(t, start.Add(5*time.Minute), end)
}

func TestSecondsRemaining(t *testing.T) {
	t.Parallel()

	// 150 seconds into a window → 150 seconds remaining
	input := time.Unix(1_700_000_100+150, 0).UTC() // 150s past window start
	got := timeutil.SecondsRemaining(input)
	assert.Equal(t, 150, got)
}
