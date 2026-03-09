package crypto_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/commons/crypto"
)

func TestGenerateSalt(t *testing.T) {
	t.Parallel()

	s1 := crypto.GenerateSalt()
	s2 := crypto.GenerateSalt()

	assert.NotNil(t, s1)
	assert.NotNil(t, s2)
	assert.NotEqual(t, s1.String(), s2.String(), "salts must be unique")

	// uint256 max = 2^256 - 1
	uint256Max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	assert.LessOrEqual(t, s1.Cmp(big.NewInt(0)), 1, "salt must be non-negative")
	assert.LessOrEqual(t, s1.Cmp(uint256Max), 0, "salt must be <= uint256 max")
}
