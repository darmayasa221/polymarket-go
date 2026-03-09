package listusers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers/dto"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
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

func makeTestUser(t *testing.T) *user.User {
	t.Helper()
	u, err := user.New(user.Params{
		Username:       "johndoe",
		Email:          "john@example.com",
		HashedPassword: "$2a$12$hashed",
		FullName:       "John Doe",
	})
	require.NoError(t, err)
	return u
}

func TestListUsers_ExecuteOffset(t *testing.T) {
	t.Parallel()

	offsetParams := pagination.NewOffsetParams(1, 10)

	tests := []struct {
		name      string
		input     dto.Input
		setup     func(repo *mockUserRepo)
		wantErr   bool
		wantCount int
	}{
		{
			name:  "returns paginated users successfully",
			input: dto.Input{OffsetParams: offsetParams},
			setup: func(repo *mockUserRepo) {
				u := makeTestUser(t)
				result := pagination.NewOffsetResult([]*user.User{u}, 1, 1, 10)
				repo.On("ListOffset", mock.Anything, offsetParams).Return(result, nil)
			},
			wantCount: 1,
		},
		{
			name:  "returns empty list when no users",
			input: dto.Input{OffsetParams: offsetParams},
			setup: func(repo *mockUserRepo) {
				result := pagination.NewOffsetResult([]*user.User{}, 0, 1, 10)
				repo.On("ListOffset", mock.Anything, offsetParams).Return(result, nil)
			},
			wantCount: 0,
		},
		{
			name:    "repo error propagates",
			input:   dto.Input{OffsetParams: offsetParams},
			wantErr: true,
			setup: func(repo *mockUserRepo) {
				repo.On("ListOffset", mock.Anything, offsetParams).
					Return(pagination.OffsetResult[*user.User]{}, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockUserRepo)
			tt.setup(repo)
			uc := listusers.New(repo)

			out, err := uc.ExecuteOffset(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, out.Users)
			} else {
				require.NoError(t, err)
				assert.Len(t, out.Users, tt.wantCount)
				if tt.wantCount > 0 {
					assert.NotEmpty(t, out.Users[0].ID)
					assert.Equal(t, "johndoe", out.Users[0].Username)
					assert.Equal(t, "john@example.com", out.Users[0].Email)
					assert.Equal(t, "John Doe", out.Users[0].FullName)
					assert.IsType(t, time.Time{}, out.Users[0].CreatedAt)
				}
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestListUsers_ExecuteCursor(t *testing.T) {
	t.Parallel()

	cursorParams := pagination.NewCursorParams("", 10, true)

	tests := []struct {
		name      string
		input     dto.Input
		setup     func(repo *mockUserRepo)
		wantErr   bool
		wantCount int
	}{
		{
			name:  "returns cursor-paginated users successfully",
			input: dto.Input{CursorParams: cursorParams},
			setup: func(repo *mockUserRepo) {
				u := makeTestUser(t)
				result := pagination.NewCursorResult([]*user.User{u}, "next-token", "", true, false)
				repo.On("ListCursor", mock.Anything, cursorParams).Return(result, nil)
			},
			wantCount: 1,
		},
		{
			name:  "returns empty list when no users",
			input: dto.Input{CursorParams: cursorParams},
			setup: func(repo *mockUserRepo) {
				result := pagination.NewCursorResult([]*user.User{}, "", "", false, false)
				repo.On("ListCursor", mock.Anything, cursorParams).Return(result, nil)
			},
			wantCount: 0,
		},
		{
			name:    "repo error propagates",
			input:   dto.Input{CursorParams: cursorParams},
			wantErr: true,
			setup: func(repo *mockUserRepo) {
				repo.On("ListCursor", mock.Anything, cursorParams).
					Return(pagination.CursorResult[*user.User]{}, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockUserRepo)
			tt.setup(repo)
			uc := listusers.New(repo)

			out, err := uc.ExecuteCursor(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, out.Users)
			} else {
				require.NoError(t, err)
				assert.Len(t, out.Users, tt.wantCount)
				if tt.wantCount > 0 {
					assert.NotEmpty(t, out.Users[0].ID)
					assert.Equal(t, "johndoe", out.Users[0].Username)
					assert.Equal(t, "john@example.com", out.Users[0].Email)
					assert.Equal(t, "John Doe", out.Users[0].FullName)
					assert.IsType(t, time.Time{}, out.Users[0].CreatedAt)
					assert.Equal(t, "next-token", out.NextCursor)
					assert.True(t, out.HasNext)
					assert.False(t, out.HasPrev)
				}
			}
			repo.AssertExpectations(t)
		})
	}
}
