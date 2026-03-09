package validation_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/validation"
	"github.com/darmayasa221/polymarket-go/internal/commons/validation/rules"
	"github.com/darmayasa221/polymarket-go/internal/commons/validation/types"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		checks    map[string]func() (bool, string, string)
		wantErr   bool
		wantField string
		wantCode  string
	}{
		{
			name: "all checks pass returns nil",
			checks: map[string]func() (bool, string, string){
				"email": func() (bool, string, string) {
					return true, "", ""
				},
				"name": func() (bool, string, string) {
					return true, "", ""
				},
			},
			wantErr: false,
		},
		{
			name: "one check fails returns MultiError with correct field and code",
			checks: map[string]func() (bool, string, string){
				"email": func() (bool, string, string) {
					return false, "USER.INVALID_EMAIL", "email is invalid"
				},
			},
			wantErr:   true,
			wantField: "email",
			wantCode:  "USER.INVALID_EMAIL",
		},
		{
			name:    "empty checks returns nil",
			checks:  map[string]func() (bool, string, string){},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — checks defined in table

			// Act
			err := validation.Validate(tt.checks)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				var multi *types.MultiError
				require.True(t, errors.As(err, &multi), "error must be *types.MultiError")
				require.True(t, multi.HasErrors())

				// Build field→code map to avoid ordering dependency.
				fieldMap := make(map[string]string, len(multi.Errors))
				for _, fe := range multi.Errors {
					fieldMap[fe.Field] = fe.Code
				}
				assert.Equal(t, tt.wantCode, fieldMap[tt.wantField])
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidate_MultipleFailures(t *testing.T) {
	t.Parallel()

	// Arrange
	checks := map[string]func() (bool, string, string){
		"email": func() (bool, string, string) {
			return false, "USER.INVALID_EMAIL", "email is invalid"
		},
		"password": func() (bool, string, string) {
			return false, "USER.WEAK_PASSWORD", "password is too weak"
		},
		"name": func() (bool, string, string) {
			return true, "", ""
		},
	}

	// Act
	err := validation.Validate(checks)

	// Assert
	require.Error(t, err)
	var multi *types.MultiError
	require.True(t, errors.As(err, &multi))
	assert.Len(t, multi.Errors, 2)

	fieldMap := make(map[string]string, len(multi.Errors))
	for _, fe := range multi.Errors {
		fieldMap[fe.Field] = fe.Code
	}
	assert.Equal(t, "USER.INVALID_EMAIL", fieldMap["email"])
	assert.Equal(t, "USER.WEAK_PASSWORD", fieldMap["password"])
}

func TestFormatErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want map[string]string
	}{
		{
			name: "MultiError maps field to code",
			err: func() error {
				m := &types.MultiError{}
				m.Add("email", "USER.INVALID_EMAIL", "email is invalid")
				m.Add("password", "USER.WEAK_PASSWORD", "password too weak")
				return m
			}(),
			want: map[string]string{
				"email":    "USER.INVALID_EMAIL",
				"password": "USER.WEAK_PASSWORD",
			},
		},
		{
			name: "non-MultiError returns nil",
			err:  errors.New("some generic error"),
			want: nil,
		},
		{
			name: "nil error returns nil",
			err:  nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — err defined in table

			// Act
			got := validation.FormatErrors(tt.err)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsRequired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "non-empty string is required",
			input: "hello",
			want:  true,
		},
		{
			name:  "empty string fails required",
			input: "",
			want:  false,
		},
		{
			name:  "whitespace-only string is non-empty",
			input: "   ",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := rules.IsRequired(tt.input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid email",
			input: "user@example.com",
			want:  true,
		},
		{
			name:  "email with subdomain",
			input: "user@mail.example.com",
			want:  true,
		},
		{
			name:  "email with plus tag",
			input: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "missing at sign",
			input: "userexample.com",
			want:  false,
		},
		{
			name:  "missing domain",
			input: "user@",
			want:  false,
		},
		{
			name:  "missing tld",
			input: "user@example",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := rules.IsEmail(tt.input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsStrongPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "strong password passes",
			input: "Secure1Pass",
			want:  true,
		},
		{
			name:  "too short fails",
			input: "Ab1",
			want:  false,
		},
		{
			name:  "no uppercase fails",
			input: "secure1pass",
			want:  false,
		},
		{
			name:  "no lowercase fails",
			input: "SECURE1PASS",
			want:  false,
		},
		{
			name:  "no digit fails",
			input: "SecurePass",
			want:  false,
		},
		{
			name:  "empty string fails",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := rules.IsStrongPassword(tt.input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsPhone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid international number with plus",
			input: "+14155552671",
			want:  true,
		},
		{
			name:  "valid number without plus",
			input: "14155552671",
			want:  true,
		},
		{
			name:  "too short fails",
			input: "+123",
			want:  false,
		},
		{
			name:  "starts with zero fails",
			input: "+0123456789",
			want:  false,
		},
		{
			name:  "empty string fails",
			input: "",
			want:  false,
		},
		{
			name:  "letters fail",
			input: "+1415abc5671",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := rules.IsPhone(tt.input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "uppercase lowercased",
			input: "USER@EXAMPLE.COM",
			want:  "user@example.com",
		},
		{
			name:  "leading and trailing spaces trimmed",
			input: "  user@example.com  ",
			want:  "user@example.com",
		},
		{
			name:  "mixed case and spaces normalized",
			input: "  User@Example.COM  ",
			want:  "user@example.com",
		},
		{
			name:  "already clean email unchanged",
			input: "user@example.com",
			want:  "user@example.com",
		},
		{
			name:  "empty string returns empty",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := validation.SanitizeEmail(tt.input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}
