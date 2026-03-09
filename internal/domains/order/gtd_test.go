package order_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

func TestGTDExpiration(t *testing.T) {
	t.Parallel()

	windowEnd := time.Unix(1_700_000_400, 0).UTC()
	exp := order.GTDExpiration(windowEnd)
	assert.Equal(t, windowEnd.Add(60*time.Second), exp)
}
