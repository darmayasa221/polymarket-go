package logoutuser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser/dto"
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

func TestLogoutUser_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   dto.Input
		setup   func(authRepo *mockAuthRepo)
		wantErr bool
		errType any
	}{
		{
			name:  "successful logout deletes token",
			input: dto.Input{TokenValue: "valid-refresh-token"},
			setup: func(authRepo *mockAuthRepo) {
				authRepo.On("DeleteByValue", mock.Anything, token.TokenValue("valid-refresh-token")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "empty token value returns client error",
			input:   dto.Input{TokenValue: ""},
			setup:   func(_ *mockAuthRepo) {},
			wantErr: true,
			errType: &errtypes.ClientError{},
		},
		{
			name:  "repo error is propagated",
			input: dto.Input{TokenValue: "some-token"},
			setup: func(authRepo *mockAuthRepo) {
				authRepo.On("DeleteByValue", mock.Anything, token.TokenValue("some-token")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			authRepo := new(mockAuthRepo)
			tt.setup(authRepo)
			uc := logoutuser.New(authRepo)

			// Act
			out, err := uc.Execute(t.Context(), tt.input)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.As(err, &tt.errType))
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, dto.Output{}, out)
			}
			authRepo.AssertExpectations(t)
		})
	}
}
