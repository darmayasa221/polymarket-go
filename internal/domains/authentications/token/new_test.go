package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
)

func TestNew_Success(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "user-id-123",
		Value:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: timeutil.AddDuration(7 * 24 * time.Hour),
	}
	tok, err := token.New(p)
	require.NoError(t, err)
	assert.False(t, tok.ID().IsEmpty())
	assert.Equal(t, tokenconst.TypeRefresh, tok.Type())
	assert.False(t, tok.IsExpired())
}

func TestNew_ValueRequired(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "user-id-123",
		Value:     "",
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: timeutil.AddDuration(time.Hour),
	}
	_, err := token.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), token.ErrValueRequired)
}

func TestToken_IsExpired(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "user-id-123",
		Value:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: timeutil.Now().Add(-time.Hour), // expired 1 hour ago
	}
	tok, err := token.New(p)
	require.NoError(t, err)
	assert.True(t, tok.IsExpired())
}

func TestNew_UserIDRequired(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "",
		Value:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: timeutil.AddDuration(time.Hour),
	}
	_, err := token.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), token.ErrUserIDRequired)
}

func TestNew_ValueTooShort(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "user-id-123",
		Value:     "short",
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: timeutil.AddDuration(time.Hour),
	}
	_, err := token.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), token.ErrValueTooShort)
}

func TestNew_TypeRequired(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "user-id-123",
		Value:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		Type:      "",
		Purpose:   "authentication",
		ExpiresAt: timeutil.AddDuration(time.Hour),
	}
	_, err := token.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), token.ErrTypeRequired)
}

func TestNew_PurposeRequired(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "user-id-123",
		Value:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		Type:      "refresh",
		Purpose:   "",
		ExpiresAt: timeutil.AddDuration(time.Hour),
	}
	_, err := token.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), token.ErrPurposeRequired)
}

func TestNew_ExpiresAtRequired(t *testing.T) {
	t.Parallel()
	p := token.Params{
		UserID:    "user-id-123",
		Value:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: time.Time{}, // zero value
	}
	_, err := token.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), token.ErrExpiresAtRequired)
}
