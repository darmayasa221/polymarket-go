package signing

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	gocrypto "github.com/ethereum/go-ethereum/crypto"
)

// Signer signs 32-byte EIP-712 hashes with an EOA private key.
// This is the ONLY place in the codebase where the private key is used.
type Signer struct {
	key *ecdsa.PrivateKey
}

// NewSigner parses a hex-encoded private key (with or without 0x prefix).
func NewSigner(privateKeyHex string) (*Signer, error) {
	if privateKeyHex == "" {
		return nil, fmt.Errorf("signing: private key is empty")
	}
	hexKey := strings.TrimPrefix(privateKeyHex, "0x")
	key, err := gocrypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, fmt.Errorf("signing: parse private key: %w", err)
	}
	return &Signer{key: key}, nil
}

// Sign produces a 65-byte ECDSA signature [R || S || V] for the given 32-byte hash.
// V is adjusted to 27 or 28 (Ethereum EIP-712 convention).
func (s *Signer) Sign(hash []byte) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("signing: hash must be 32 bytes, got %d", len(hash))
	}
	sig, err := gocrypto.Sign(hash, s.key)
	if err != nil {
		return nil, fmt.Errorf("signing: ecdsa sign: %w", err)
	}
	// go-ethereum Sign returns V as 0 or 1.
	// Polymarket expects V = 27 or 28 (pre-EIP-155 Ethereum convention).
	sig[64] += 27
	return sig, nil
}
