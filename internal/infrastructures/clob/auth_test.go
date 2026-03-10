package clob_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
)

// validTestSecret returns a base64-encoded 32-byte key for tests.
func validTestSecret() string {
	return base64.StdEncoding.EncodeToString(make([]byte, 32))
}

func TestBuildL2Signature(t *testing.T) {
	t.Parallel()
	sig, err := clob.BuildL2Signature(validTestSecret(), "1741612800", "GET", "/fee-rate", "")
	require.NoError(t, err)
	// Result must be valid base64 encoding a 32-byte HMAC-SHA256.
	decoded, decErr := base64.StdEncoding.DecodeString(sig)
	require.NoError(t, decErr)
	assert.Len(t, decoded, 32)
}

func TestBuildL2Signature_BadSecret(t *testing.T) {
	t.Parallel()
	_, err := clob.BuildL2Signature("not-base64!!!", "1741612800", "GET", "/fee-rate", "")
	require.Error(t, err)
}
