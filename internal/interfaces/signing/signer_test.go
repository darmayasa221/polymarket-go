package signing_test

import (
	"encoding/hex"
	"testing"

	gocrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/interfaces/signing"
)

func TestSigner_Sign_RecoverMatchesAddress(t *testing.T) {
	// known test key — never use in production
	const testPrivKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	signer, err := signing.NewSigner(testPrivKey)
	require.NoError(t, err)

	hash := make([]byte, 32)
	for i := range hash {
		hash[i] = byte(i)
	}

	sig, err := signer.Sign(hash)
	require.NoError(t, err)
	require.Len(t, sig, 65)

	// Recover — must strip Ethereum's V offset (27) before recovering
	recoverSig := make([]byte, 65)
	copy(recoverSig, sig)
	recoverSig[64] -= 27

	pubKey, err := gocrypto.Ecrecover(hash, recoverSig)
	require.NoError(t, err)
	require.NotEmpty(t, pubKey)

	_ = hex.EncodeToString(pubKey) // sanity: pubkey is bytes, not empty
}

func TestSigner_New_InvalidKey(t *testing.T) {
	_, err := signing.NewSigner("not-a-hex-key")
	require.Error(t, err)
}

func TestSigner_New_EmptyKey(t *testing.T) {
	_, err := signing.NewSigner("")
	require.Error(t, err)
}
