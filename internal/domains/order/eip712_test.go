package order_test

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

func TestSigningHash_ReturnsBytes(t *testing.T) {
	t.Parallel()

	unsigned := order.UnsignedOrder{
		Salt:          big.NewInt(12345),
		Maker:         "0xMakerAddress",
		Signer:        "0xSignerAddress",
		Taker:         "0x0000000000000000000000000000000000000000",
		TokenID:       polyid.TokenID("111"),
		MakerAmount:   decimal.NewFromFloat(10),
		TakerAmount:   decimal.NewFromFloat(6.5),
		Expiration:    1700000460,
		Nonce:         0,
		FeeRateBps:    50,
		Side:          order.Buy,
		SignatureType: 0,
	}

	hash, err := order.SigningHash(unsigned, 137)
	require.NoError(t, err)
	assert.Len(t, hash, 32, "EIP-712 signing hash must be 32 bytes")
}
