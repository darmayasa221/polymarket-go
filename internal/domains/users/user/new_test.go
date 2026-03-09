package user_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

func TestNew_Success(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hashedpasswordhere",
		FullName:       "John Doe",
	}
	u, err := user.New(p)
	require.NoError(t, err)
	assert.False(t, u.ID().IsEmpty())
	assert.Equal(t, "johndoe", u.Username())
	assert.Equal(t, user.Email("john@example.com"), u.Email())
	assert.Equal(t, "John Doe", u.FullName())
	assert.False(t, u.CreatedAt().IsZero())
}

func TestNew_UsernameRequired(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       "John Doe",
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrUsernameRequired)
}

func TestNew_UsernameTooShort(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "ab",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       "John Doe",
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrUsernameTooShort)
}

func TestNew_EmailInvalid(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "johndoe",
		Email:          "not-an-email",
		HashedPassword: "$2a$12$hash",
		FullName:       "John Doe",
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrEmailInvalid)
}

func TestNew_PasswordRequired(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "",
		FullName:       "John Doe",
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrPasswordRequired)
}

func TestNew_UsernameTooLong(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       strings.Repeat("a", user.UsernameMaxLength+1),
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       "John Doe",
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrUsernameTooLong)
}

func TestNew_EmailTooLong(t *testing.T) {
	t.Parallel()
	localPart := strings.Repeat("a", user.EmailMaxLength)
	p := user.Params{
		Username:       "johndoe",
		Email:          localPart + "@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       "John Doe",
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrEmailTooLong)
}

func TestNew_FullNameRequired(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       "",
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrFullNameRequired)
}

func TestNew_FullNameTooShort(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       strings.Repeat("a", user.FullNameMinLength-1),
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrFullNameTooShort)
}

func TestNew_FullNameTooLong(t *testing.T) {
	t.Parallel()
	p := user.Params{
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       strings.Repeat("a", user.FullNameMaxLength+1),
	}
	_, err := user.New(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), user.ErrFullNameTooLong)
}

func TestReconstitute(t *testing.T) {
	t.Parallel()
	p := user.ReconstitutedParams{
		ID:             "existing-id-123",
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hash",
		FullName:       "John Doe",
	}
	u := user.Reconstitute(p)
	assert.Equal(t, "existing-id-123", u.ID().String())
	assert.Equal(t, "johndoe", u.Username())
}
