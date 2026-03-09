// Package polyid_test tests the polyid value objects.
package polyid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
)

func TestConditionID(t *testing.T) {
	t.Parallel()

	id := polyid.ConditionID("0xabc123")
	assert.Equal(t, "0xabc123", id.String())
	assert.False(t, id.IsEmpty())
	assert.True(t, polyid.ConditionID("").IsEmpty())
}

func TestTokenID(t *testing.T) {
	t.Parallel()

	id := polyid.TokenID("12345678901234567890")
	assert.Equal(t, "12345678901234567890", id.String())
	assert.False(t, id.IsEmpty())
	assert.True(t, polyid.TokenID("").IsEmpty())
}

func TestOrderID(t *testing.T) {
	t.Parallel()

	id := polyid.OrderID("0xdeadbeef")
	assert.Equal(t, "0xdeadbeef", id.String())
	assert.False(t, id.IsEmpty())
	assert.True(t, polyid.OrderID("").IsEmpty())
}

func TestSlugID(t *testing.T) {
	t.Parallel()

	id := polyid.SlugID("btc-updown-5m-1700000100")
	assert.Equal(t, "btc-updown-5m-1700000100", id.String())
	assert.False(t, id.IsEmpty())
	assert.True(t, polyid.SlugID("").IsEmpty())
}
