package refreshauth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth/dto"
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	sectypes "github.com/darmayasa221/polymarket-go/internal/applications/security/types"
	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

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
var _ security.TokenManager = (*mockTokenManager)(nil)

func TestRefreshAuth_Execute(t *testing.T) {
	t.Parallel()

	const validRefreshValue = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.old-refresh-token"
	const userID = "550e8400-e29b-41d4-a716-446655440000"

	futureTime := time.Now().Add(time.Hour)

	validToken, err := token.New(token.Params{
		UserID:    userID,
		Value:     validRefreshValue,
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: futureTime,
	})
	require.NoError(t, err)

	expiredToken := token.Reconstitute(token.ReconstitutedParams{
		ID:        "some-id",
		UserID:    userID,
		Value:     validRefreshValue,
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-2 * time.Hour),
	})

	newPair := sectypes.TokenPair{
		AccessToken:           "new-access-token",
		RefreshToken:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new-refresh-token",
		AccessTokenExpiresAt:  futureTime,
		RefreshTokenExpiresAt: futureTime,
	}

	tests := []struct {
		name    string
		input   dto.Input
		setup   func(authRepo *mockAuthRepo, tm *mockTokenManager)
		wantErr bool
		errType any
	}{
		{
			name:  "successful refresh returns new token pair",
			input: dto.Input{RefreshToken: validRefreshValue},
			setup: func(authRepo *mockAuthRepo, tm *mockTokenManager) {
				authRepo.On("CheckToken", mock.Anything, token.TokenValue(validRefreshValue)).Return(validToken, nil)
				tm.On("CreateTokenPair", mock.Anything, userID).Return(newPair, nil)
				authRepo.On("DeleteByValue", mock.Anything, token.TokenValue(validRefreshValue)).Return(nil)
				authRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "empty refresh token returns client error",
			input:   dto.Input{RefreshToken: ""},
			setup:   func(_ *mockAuthRepo, _ *mockTokenManager) {},
			wantErr: true,
			errType: &errtypes.ClientError{},
		},
		{
			name:  "token not found returns authentication error",
			input: dto.Input{RefreshToken: validRefreshValue},
			setup: func(authRepo *mockAuthRepo, _ *mockTokenManager) {
				authRepo.On("CheckToken", mock.Anything, token.TokenValue(validRefreshValue)).Return((*token.Token)(nil), errors.New("not found"))
			},
			wantErr: true,
			errType: &errtypes.AuthenticationError{},
		},
		{
			name:  "expired token returns authentication error",
			input: dto.Input{RefreshToken: validRefreshValue},
			setup: func(authRepo *mockAuthRepo, _ *mockTokenManager) {
				authRepo.On("CheckToken", mock.Anything, token.TokenValue(validRefreshValue)).Return(expiredToken, nil)
			},
			wantErr: true,
			errType: &errtypes.AuthenticationError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			authRepo := new(mockAuthRepo)
			tm := new(mockTokenManager)
			tt.setup(authRepo, tm)
			uc := refreshauth.New(authRepo, tm)

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
				assert.False(t, out.AccessTokenExpiresAt.IsZero())
				assert.False(t, out.RefreshTokenExpiresAt.IsZero())
			}
			authRepo.AssertExpectations(t)
			tm.AssertExpectations(t)
		})
	}
}
