package adduser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
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

// mockEncryption implements security.Encryption interface.
type mockEncryption struct{ mock.Mock }

func (m *mockEncryption) Hash(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}
func (m *mockEncryption) Compare(ctx context.Context, hashedPassword, plainPassword string) error {
	return m.Called(ctx, hashedPassword, plainPassword).Error(0)
}

func TestAddUser_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   dto.Input
		setup   func(repo *mockUserRepo, enc *mockEncryption)
		wantErr bool
		errType any
	}{
		{
			name:  "valid user registered successfully",
			input: dto.Input{Username: "johndoe", Email: "john@example.com", Password: "SecurePass123", FullName: "John Doe"},
			setup: func(repo *mockUserRepo, enc *mockEncryption) {
				repo.On("VerifyUsername", mock.Anything, "johndoe").Return(false, nil)
				enc.On("Hash", mock.Anything, "SecurePass123").Return("$2a$12$hashed", nil)
				repo.On("Add", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "duplicate username returns conflict error",
			input: dto.Input{Username: "johndoe", Email: "john@example.com", Password: "SecurePass123", FullName: "John Doe"},
			setup: func(repo *mockUserRepo, enc *mockEncryption) {
				repo.On("VerifyUsername", mock.Anything, "johndoe").Return(true, nil)
			},
			wantErr: true,
			errType: &errtypes.ConflictError{},
		},
		{
			name:  "invalid email fails entity validation",
			input: dto.Input{Username: "johndoe", Email: "not-an-email", Password: "SecurePass123", FullName: "John Doe"},
			setup: func(repo *mockUserRepo, enc *mockEncryption) {
				repo.On("VerifyUsername", mock.Anything, "johndoe").Return(false, nil)
				enc.On("Hash", mock.Anything, "SecurePass123").Return("$2a$12$hashed", nil)
			},
			wantErr: true,
			errType: &errtypes.InvariantError{},
		},
		{
			name:  "verify username repo error returns internal error",
			input: dto.Input{Username: "johndoe", Email: "john@example.com", Password: "SecurePass123", FullName: "John Doe"},
			setup: func(repo *mockUserRepo, enc *mockEncryption) {
				repo.On("VerifyUsername", mock.Anything, "johndoe").Return(false, errors.New("db error"))
			},
			wantErr: true,
			errType: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			repo := new(mockUserRepo)
			enc := new(mockEncryption)
			tt.setup(repo, enc)
			uc := adduser.New(repo, enc)

			// Act
			out, err := uc.Execute(t.Context(), tt.input)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.As(err, &tt.errType))
				}
				assert.Empty(t, out.ID)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, out.ID)
				assert.Equal(t, tt.input.Username, out.Username)
				assert.Equal(t, tt.input.Email, out.Email)
			}
			repo.AssertExpectations(t)
			enc.AssertExpectations(t)
		})
	}
}
