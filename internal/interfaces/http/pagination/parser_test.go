package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	httppagination "github.com/darmayasa221/polymarket-go/internal/interfaces/http/pagination"
)

func newTestContext(queryParams string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?"+queryParams, http.NoBody)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}

func TestParseOffset(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		expectedPage int
		expectedSize int
	}{
		{"defaults when empty", "", 1, 10},
		{"valid params", "page=3&page_size=20", 3, 20},
		{"invalid page falls back to default", "page=abc", 1, 10},
		{"invalid page_size falls back to default", "page_size=xyz", 1, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestContext(tt.query)
			params := httppagination.ParseOffset(c)
			assert.Equal(t, tt.expectedPage, params.Page)
			assert.Equal(t, tt.expectedSize, params.PageSize)
		})
	}
}

func TestParseCursor(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		expectedCursor  string
		expectedSize    int
		expectedForward bool
	}{
		{"defaults when empty", "", "", 10, true},
		{"with cursor and size", "cursor=abc123&page_size=5", "abc123", 5, true},
		{"backward direction", "direction=backward", "", 10, false},
		{"invalid page_size falls back to default", "page_size=xyz", "", 10, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestContext(tt.query)
			params := httppagination.ParseCursor(c)
			assert.Equal(t, tt.expectedCursor, params.Cursor)
			assert.Equal(t, tt.expectedSize, params.PageSize)
			assert.Equal(t, tt.expectedForward, params.Forward)
		})
	}
}
