package timeutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
)

func TestNow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "returns UTC location"},
		{name: "is not zero"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — no setup needed.

			// Act
			got := timeutil.Now()

			// Assert
			assert.False(t, got.IsZero(), "Now() must not be zero")
			assert.Equal(t, time.UTC, got.Location(), "Now() must return UTC location")
		})
	}
}

func TestUnix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "returns positive value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — no setup needed.

			// Act
			got := timeutil.Unix()

			// Assert
			assert.Greater(t, got, int64(0), "Unix() must return a positive timestamp")
		})
	}
}

func TestParseRFC3339(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantUTC bool
	}{
		{
			name:    "valid RFC3339 string parses to UTC",
			input:   "2026-03-09T12:00:00Z",
			wantErr: false,
			wantUTC: true,
		},
		{
			name:    "valid RFC3339 with offset parses to UTC",
			input:   "2026-03-09T12:00:00+07:00",
			wantErr: false,
			wantUTC: true,
		},
		{
			name:    "invalid string returns error",
			input:   "not-a-date",
			wantErr: true,
			wantUTC: false,
		},
		{
			name:    "empty string returns error",
			input:   "",
			wantErr: true,
			wantUTC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			got, err := timeutil.ParseRFC3339(tt.input)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, got.IsZero(), "on error the returned time must be zero")
			} else {
				require.NoError(t, err)
				assert.Equal(t, time.UTC, got.Location(), "parsed time must be UTC")
			}
		})
	}
}

func TestFormatRFC3339(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "round-trip UTC time",
			input: "2026-03-09T12:00:00Z",
		},
		{
			name:  "round-trip midnight UTC",
			input: "2000-01-01T00:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			parsed, err := timeutil.ParseRFC3339(tt.input)
			require.NoError(t, err)

			// Act
			got := timeutil.FormatRFC3339(parsed)

			// Assert
			assert.Equal(t, tt.input, got, "FormatRFC3339(ParseRFC3339(s)) must equal original string")
		})
	}
}

func TestIsExpired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    time.Time
		want bool
	}{
		{
			name: "past time returns true",
			t:    timeutil.Now().Add(-time.Hour),
			want: true,
		},
		{
			name: "future time returns false",
			t:    timeutil.Now().Add(time.Hour),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			got := timeutil.IsExpired(tt.t)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsAfter(t *testing.T) {
	t.Parallel()

	base := time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		t         time.Time
		reference time.Time
		want      bool
	}{
		{
			name:      "t after reference returns true",
			t:         base.Add(time.Hour),
			reference: base,
			want:      true,
		},
		{
			name:      "t before reference returns false",
			t:         base.Add(-time.Hour),
			reference: base,
			want:      false,
		},
		{
			name:      "t equal to reference returns false",
			t:         base,
			reference: base,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			got := timeutil.IsAfter(tt.t, tt.reference)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{name: "positive duration is after now", duration: time.Hour},
		{name: "large positive duration is after now", duration: 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			before := timeutil.Now()

			// Act
			got := timeutil.AddDuration(tt.duration)

			// Assert
			assert.True(t, got.After(before), "AddDuration with positive duration must return time after now")
		})
	}
}

func TestAddDays(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		n    int
	}{
		{name: "1 day is after now", n: 1},
		{name: "30 days is after now", n: 30},
		{name: "365 days is after now", n: 365},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			before := timeutil.Now()

			// Act
			got := timeutil.AddDays(tt.n)

			// Assert
			assert.True(t, got.After(before), "AddDays with positive n must return time after now")
		})
	}
}
