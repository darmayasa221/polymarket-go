package pagination_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
)

func TestNewOffsetResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		items          []int
		total          int
		page           int
		pageSize       int
		wantPageSize   int
		wantTotalPages int
	}{
		{
			name:           "pageSize zero falls back to DefaultPageSize",
			items:          []int{1, 2, 3},
			total:          30,
			page:           1,
			pageSize:       0,
			wantPageSize:   pagination.DefaultPageSize,
			wantTotalPages: 3,
		},
		{
			name:           "pageSize negative falls back to DefaultPageSize",
			items:          []int{1},
			total:          5,
			page:           1,
			pageSize:       -1,
			wantPageSize:   pagination.DefaultPageSize,
			wantTotalPages: 1,
		},
		{
			name:           "remainder causes totalPages to round up",
			items:          []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			total:          11,
			page:           1,
			pageSize:       10,
			wantPageSize:   10,
			wantTotalPages: 2,
		},
		{
			name:           "exact division produces exact totalPages",
			items:          []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			total:          20,
			page:           2,
			pageSize:       10,
			wantPageSize:   10,
			wantTotalPages: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := pagination.NewOffsetResult(tc.items, tc.total, tc.page, tc.pageSize)

			assert.Equal(t, tc.total, result.TotalItems)
			assert.Equal(t, tc.page, result.Page)
			assert.Equal(t, tc.wantPageSize, result.PageSize)
			assert.Equal(t, tc.wantTotalPages, result.TotalPages)
			assert.Equal(t, tc.items, result.Items)
		})
	}
}

func TestNewCursorResult(t *testing.T) {
	t.Parallel()

	items := []string{"a", "b", "c"}
	nextCursor := "next-token"
	prevCursor := "prev-token"

	result := pagination.NewCursorResult(items, nextCursor, prevCursor, true, false)

	assert.Equal(t, items, result.Items)
	assert.Equal(t, nextCursor, result.NextCursor)
	assert.Equal(t, prevCursor, result.PrevCursor)
	assert.True(t, result.HasNext)
	assert.False(t, result.HasPrev)
}

func TestNewOffsetParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{
			name:         "page below DefaultPage is clamped to 1",
			page:         0,
			pageSize:     10,
			wantPage:     pagination.DefaultPage,
			wantPageSize: 10,
		},
		{
			name:         "pageSize above MaxPageSize is clamped to DefaultPageSize",
			page:         1,
			pageSize:     pagination.MaxPageSize + 1,
			wantPage:     1,
			wantPageSize: pagination.DefaultPageSize,
		},
		{
			name:         "pageSize below MinPageSize is clamped to DefaultPageSize",
			page:         1,
			pageSize:     pagination.MinPageSize - 1,
			wantPage:     1,
			wantPageSize: pagination.DefaultPageSize,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			params := pagination.NewOffsetParams(tc.page, tc.pageSize)

			assert.Equal(t, tc.wantPage, params.Page)
			assert.Equal(t, tc.wantPageSize, params.PageSize)
		})
	}
}

func TestOffsetParams_Offset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantOffset int
	}{
		{
			name:       "page 1 offset is zero",
			page:       1,
			pageSize:   10,
			wantOffset: 0,
		},
		{
			name:       "page 3 with pageSize 10 gives offset 20",
			page:       3,
			pageSize:   10,
			wantOffset: 20,
		},
		{
			name:       "page 2 with pageSize 5 gives offset 5",
			page:       2,
			pageSize:   5,
			wantOffset: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			params := pagination.NewOffsetParams(tc.page, tc.pageSize)

			assert.Equal(t, tc.wantOffset, params.Offset())
		})
	}
}

func TestNewCursorParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cursor       string
		pageSize     int
		forward      bool
		wantPageSize int
	}{
		{
			name:         "pageSize above MaxPageSize is clamped to DefaultPageSize",
			cursor:       "some-cursor",
			pageSize:     pagination.MaxPageSize + 1,
			forward:      true,
			wantPageSize: pagination.DefaultPageSize,
		},
		{
			name:         "valid pageSize is preserved",
			cursor:       "cursor-abc",
			pageSize:     25,
			forward:      false,
			wantPageSize: 25,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			params := pagination.NewCursorParams(tc.cursor, tc.pageSize, tc.forward)

			assert.Equal(t, tc.cursor, params.Cursor)
			assert.Equal(t, tc.wantPageSize, params.PageSize)
			assert.Equal(t, tc.forward, params.Forward)
		})
	}
}
