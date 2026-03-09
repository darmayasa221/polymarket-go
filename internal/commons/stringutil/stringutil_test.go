package stringutil_test

import (
	"testing"

	"github.com/darmayasa221/polymarket-go/internal/commons/stringutil"
)

func TestMaskEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid email masked correctly",
			input: "user@domain.com",
			want:  "u***@domain.com",
		},
		{
			name:  "email with no at sign",
			input: "nodomain",
			want:  "***",
		},
		{
			name:  "empty local part",
			input: "@domain.com",
			want:  "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := stringutil.MaskEmail(tt.input)

			// Assert
			if got != tt.want {
				t.Errorf("MaskEmail(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMaskString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "string longer than two chars",
			input: "secret",
			want:  "s****t",
		},
		{
			name:  "string exactly two chars",
			input: "ab",
			want:  "***",
		},
		{
			name:  "single char string",
			input: "x",
			want:  "***",
		},
		{
			name:  "empty string",
			input: "",
			want:  "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := stringutil.MaskString(tt.input)

			// Assert
			if got != tt.want {
				t.Errorf("MaskString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "camelCase to snake_case",
			input: "camelCaseString",
			want:  "camel_case_string",
		},
		{
			name:  "already lowercase",
			input: "lowercase",
			want:  "lowercase",
		},
		{
			name:  "single word",
			input: "word",
			want:  "word",
		},
		{
			name:  "acronym followed by word",
			input: "HTTPRequest",
			want:  "http_request",
		},
		{
			name:  "word followed by acronym",
			input: "MyHTTPServer",
			want:  "my_http_server",
		},
		{
			name:  "PascalCase",
			input: "PascalCase",
			want:  "pascal_case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := stringutil.ToSnakeCase(tt.input)

			// Assert
			if got != tt.want {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTrimSpaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "leading and trailing spaces removed",
			input: "  hello world  ",
			want:  "hello world",
		},
		{
			name:  "no spaces unchanged",
			input: "nospaces",
			want:  "nospaces",
		},
		{
			name:  "only spaces returns empty",
			input: "   ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange — input defined in table

			// Act
			got := stringutil.TrimSpaces(tt.input)

			// Assert
			if got != tt.want {
				t.Errorf("TrimSpaces(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
