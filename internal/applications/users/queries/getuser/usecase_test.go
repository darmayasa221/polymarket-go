package getuser_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

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

func TestGetUser_Execute(t *testing.T) {
	t.Parallel()

	existingUser, err := user.New(user.Params{
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hashed",
		FullName:       "John Doe",
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   dto.Input
		setup   func(repo *mockUserRepo)
		wantErr bool
		errType any
	}{
		{
			name:  "returns user by ID",
			input: dto.Input{UserID: "some-uuid"},
			setup: func(repo *mockUserRepo) {
				repo.On("GetByID", mock.Anything, user.UserID("some-uuid")).Return(existingUser, nil)
			},
		},
		{
			name:    "empty ID returns client error",
			input:   dto.Input{UserID: ""},
			setup:   func(repo *mockUserRepo) {},
			wantErr: true,
			errType: &errtypes.ClientError{},
		},
		{
			name:  "not found returns error from repo",
			input: dto.Input{UserID: "missing-id"},
			setup: func(repo *mockUserRepo) {
				repo.On("GetByID", mock.Anything, user.UserID("missing-id")).Return((*user.User)(nil), errtypes.NewNotFoundError("USER_REPO.NOT_FOUND"))
			},
			wantErr: true,
			errType: &errtypes.NotFoundError{},
		},
		{
			name:  "repo error propagates",
			input: dto.Input{UserID: "some-uuid"},
			setup: func(repo *mockUserRepo) {
				repo.On("GetByID", mock.Anything, user.UserID("some-uuid")).Return((*user.User)(nil), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockUserRepo)
			tt.setup(repo)
			uc := getuser.New(repo)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.As(err, &tt.errType))
				}
				assert.Empty(t, out.ID)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, out.ID)
				assert.Equal(t, existingUser.Username(), out.Username)
				assert.Equal(t, existingUser.Email().String(), out.Email)
				assert.False(t, out.CreatedAt.IsZero())
				assert.IsType(t, time.Time{}, out.UpdatedAt)
			}
			repo.AssertExpectations(t)
		})
	}
}
