package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/commons/crypto"
)

func TestGenerateUUID(t *testing.T) {
	id1 := crypto.GenerateUUID()
	id2 := crypto.GenerateUUID()
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.Len(t, id1, 36) // UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
}

func TestSHA256(t *testing.T) {
	hash := crypto.SHA256("hello")
	assert.NotEmpty(t, hash)
	assert.Equal(t, crypto.SHA256("hello"), hash) // deterministic
	assert.NotEqual(t, crypto.SHA256("hello"), crypto.SHA256("world"))
}

func TestBase64EncodeDecode(t *testing.T) {
	original := []byte("test-data-123")
	encoded := crypto.Base64Encode(original)
	decoded, err := crypto.Base64Decode(encoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}
