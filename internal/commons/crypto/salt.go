package crypto

import (
	"crypto/rand"
	"math/big"
)

// maxUint256 is 2^256 - 1 — the maximum value for a Solidity uint256.
var maxUint256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

// GenerateSalt returns a cryptographically random uint256 for EIP-712 order salt.
// Panics only if the system entropy source fails (unrecoverable).
func GenerateSalt() *big.Int {
	n, err := rand.Int(rand.Reader, maxUint256)
	if err != nil {
		panic("crypto.GenerateSalt: entropy source failure: " + err.Error())
	}
	return n
}
