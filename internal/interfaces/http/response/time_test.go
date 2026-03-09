package response_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

func TestJSONTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "RFC3339 without nanoseconds",
			input:    time.Date(2026, 3, 10, 12, 34, 56, 123456789, time.UTC),
			expected: `"2026-03-10T12:34:56Z"`,
		},
		{
			name:     "non-UTC time is converted to UTC",
			input:    time.Date(2026, 3, 10, 12, 34, 56, 0, time.FixedZone("WIB", 7*3600)),
			expected: `"2026-03-10T05:34:56Z"`,
		},
		{
			name:     "zero time marshals correctly",
			input:    time.Time{},
			expected: `"0001-01-01T00:00:00Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jt := response.JSONTime(tt.input)
			data, err := json.Marshal(jt)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestJSONTime_UnmarshalJSON(t *testing.T) {
	t.Run("round-trips correctly", func(t *testing.T) {
		original := time.Date(2026, 3, 10, 12, 34, 56, 0, time.UTC)
		jt := response.JSONTime(original)

		data, err := json.Marshal(jt)
		require.NoError(t, err)

		var result response.JSONTime
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		assert.Equal(t, original, time.Time(result))
	})

	t.Run("parses valid RFC3339 string", func(t *testing.T) {
		data := []byte(`"2026-03-10T12:34:56Z"`)
		var jt response.JSONTime
		err := json.Unmarshal(data, &jt)
		require.NoError(t, err)

		expected := time.Date(2026, 3, 10, 12, 34, 56, 0, time.UTC)
		assert.Equal(t, expected, time.Time(jt))
	})

	t.Run("returns error for invalid format", func(t *testing.T) {
		data := []byte(`"not-a-date"`)
		var jt response.JSONTime
		err := json.Unmarshal(data, &jt)
		assert.Error(t, err)
	})

	t.Run("null input returns zero value without error", func(t *testing.T) {
		var jt response.JSONTime
		err := json.Unmarshal([]byte("null"), &jt)
		require.NoError(t, err)
		assert.Equal(t, time.Time{}, time.Time(jt))
	})
}
