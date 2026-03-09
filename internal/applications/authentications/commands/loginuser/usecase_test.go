package loginuser_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	sectypes "github.com/darmayasa221/polymarket-go/internal/applications/security/types"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// mockUserRepo implements repository.User interface.
type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Add(ctx context.Context, u *user.User) error {
	return m.Called(ctx, u).Error(0)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id user.UserID) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepo) GetIDByUsername(ctx context.Context, username string) (user.UserID, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(user.UserID), args.Error(1)
}

func (m *mockUserRepo) GetPassword(ctx context.Context, username string) (user.HashedPassword, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(user.HashedPassword), args.Error(1)
}

func (m *mockUserRepo) VerifyUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepo) ListOffset(ctx context.Context, params pagination.OffsetParams) (pagination.OffsetResult[*user.User], error) {
	args := m.Called(ctx, params)
	return args.Get(0).(pagination.OffsetResult[*user.User]), args.Error(1)
}

func (m *mockUserRepo) ListCursor(ctx context.Context, params pagination.CursorParams) (pagination.CursorResult[*user.User], error) {
	args := m.Called(ctx, params)
	return args.Get(0).(pagination.CursorResult[*user.User]), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	return m.Called(ctx, u).Error(0)
}

func (m *mockUserRepo) Delete(ctx context.Context, id user.UserID) error {
	return m.Called(ctx, id).Error(0)
}

// mockAuthRepo implements repository.Authentication interface.
type mockAuthRepo struct{ mock.Mock }

func (m *mockAuthRepo) Add(ctx context.Context, t *token.Token) error {
	return m.Called(ctx, t).Error(0)
}

func (m *mockAuthRepo) CheckToken(ctx context.Context, value token.TokenValue) (*token.Token, error) {
	args := m.Called(ctx, value)
	return args.Get(0).(*token.Token), args.Error(1)
}

func (m *mockAuthRepo) DeleteByUserID(ctx context.Context, userID valueobjects.ID) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *mockAuthRepo) DeleteByValue(ctx context.Context, value token.TokenValue) error {
	return m.Called(ctx, value).Error(0)
}

// mockEncryption implements security.Encryption interface.
type mockEncryption struct{ mock.Mock }

func (m *mockEncryption) Hash(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *mockEncryption) Compare(ctx context.Context, hashedPassword, plainPassword string) error {
	return m.Called(ctx, hashedPassword, plainPassword).Error(0)
}

// mockTokenManager implements security.TokenManager interface.
type mockTokenManager struct{ mock.Mock }

func (m *mockTokenManager) CreateTokenPair(ctx context.Context, userID string) (sectypes.TokenPair, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(sectypes.TokenPair), args.Error(1)
}

func (m *mockTokenManager) VerifyAccessToken(ctx context.Context, tokenValue string) (*sectypes.TokenClaims, error) {
	args := m.Called(ctx, tokenValue)
	return args.Get(0).(*sectypes.TokenClaims), args.Error(1)
}

func (m *mockTokenManager) VerifyRefreshToken(ctx context.Context, tokenValue string) (*sectypes.TokenClaims, error) {
	args := m.Called(ctx, tokenValue)
	return args.Get(0).(*sectypes.TokenClaims), args.Error(1)
}

// compile-time interface assertions.
var (
	_ security.Encryption   = (*mockEncryption)(nil)
	_ security.TokenManager = (*mockTokenManager)(nil)
)

func TestLoginUser_Execute(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now().Add(time.Hour)
	validPair := sectypes.TokenPair{
		AccessToken:           "access-token-value",
		RefreshToken:          "refresh-token-value",
		AccessTokenExpiresAt:  fixedTime,
		RefreshTokenExpiresAt: fixedTime,
	}
	validUserID := user.UserID("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name    string
		input   dto.Input
		setup   func(repo *mockUserRepo, authRepo *mockAuthRepo, enc *mockEncryption, tm *mockTokenManager)
		wantErr bool
		errType any
	}{
		{
			name:  "successful login returns token pair",
			input: dto.Input{Username: "johndoe", Password: "SecurePass123"},
			setup: func(repo *mockUserRepo, authRepo *mockAuthRepo, enc *mockEncryption, tm *mockTokenManager) {
				repo.On("GetPassword", mock.Anything, "johndoe").Return(user.HashedPassword("$2a$12$hashed"), nil)
				enc.On("Compare", mock.Anything, "$2a$12$hashed", "SecurePass123").Return(nil)
				repo.On("GetIDByUsername", mock.Anything, "johndoe").Return(validUserID, nil)
				tm.On("CreateTokenPair", mock.Anything, validUserID.String()).Return(validPair, nil)
				authRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "empty username returns client error",
			input:   dto.Input{Username: "", Password: "SecurePass123"},
			setup:   func(_ *mockUserRepo, _ *mockAuthRepo, _ *mockEncryption, _ *mockTokenManager) {},
			wantErr: true,
			errType: &errtypes.ClientError{},
		},
		{
			name:    "empty password returns client error",
			input:   dto.Input{Username: "johndoe", Password: ""},
			setup:   func(_ *mockUserRepo, _ *mockAuthRepo, _ *mockEncryption, _ *mockTokenManager) {},
			wantErr: true,
			errType: &errtypes.ClientError{},
		},
		{
			name:  "wrong password returns authentication error",
			input: dto.Input{Username: "johndoe", Password: "WrongPass"},
			setup: func(repo *mockUserRepo, _ *mockAuthRepo, enc *mockEncryption, _ *mockTokenManager) {
				repo.On("GetPassword", mock.Anything, "johndoe").Return(user.HashedPassword("$2a$12$hashed"), nil)
				enc.On("Compare", mock.Anything, "$2a$12$hashed", "WrongPass").Return(errors.New("mismatch"))
			},
			wantErr: true,
			errType: &errtypes.AuthenticationError{},
		},
		{
			name:  "user not found returns authentication error",
			input: dto.Input{Username: "unknown", Password: "SecurePass123"},
			setup: func(repo *mockUserRepo, _ *mockAuthRepo, _ *mockEncryption, _ *mockTokenManager) {
				repo.On("GetPassword", mock.Anything, "unknown").Return(user.HashedPassword(""), errors.New("not found"))
			},
			wantErr: true,
			errType: &errtypes.AuthenticationError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repo := new(mockUserRepo)
			authRepo := new(mockAuthRepo)
			enc := new(mockEncryption)
			tm := new(mockTokenManager)
			tt.setup(repo, authRepo, enc, tm)
			uc := loginuser.New(repo, authRepo, enc, tm)

			// Act
			out, err := uc.Execute(t.Context(), tt.input)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.As(err, &tt.errType))
				}
				assert.Empty(t, out.AccessToken)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, out.AccessToken)
				assert.NotEmpty(t, out.RefreshToken)
				assert.NotEmpty(t, out.UserID)
			}
			repo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
			enc.AssertExpectations(t)
			tm.AssertExpectations(t)
		})
	}
}
