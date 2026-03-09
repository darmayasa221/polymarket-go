package pagination

// OffsetParams holds parameters for offset-based pagination.
type OffsetParams struct {
	Page     int
	PageSize int
}

// CursorParams holds parameters for cursor-based pagination.
type CursorParams struct {
	Cursor   string
	PageSize int
	Forward  bool // true = forward, false = backward
}

// NewOffsetParams creates OffsetParams with defaults applied.
func NewOffsetParams(page, pageSize int) OffsetParams {
	if page < DefaultPage {
		page = DefaultPage
	}
	if pageSize < MinPageSize || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}
	return OffsetParams{Page: page, PageSize: pageSize}
}

// Offset returns the number of items to skip before the current page.
func (p OffsetParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// NewCursorParams creates CursorParams with defaults applied.
func NewCursorParams(cursor string, pageSize int, forward bool) CursorParams {
	if pageSize < MinPageSize || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}
	return CursorParams{Cursor: cursor, PageSize: pageSize, Forward: forward}
}
