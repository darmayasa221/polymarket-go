package pagination

// OffsetResult holds paginated results with offset metadata.
type OffsetResult[T any] struct {
	Items      []T
	TotalItems int
	Page       int
	PageSize   int
	TotalPages int
}

// NewOffsetResult creates an OffsetResult and calculates total pages.
func NewOffsetResult[T any](items []T, total, page, pageSize int) OffsetResult[T] {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}
	return OffsetResult[T]{
		Items:      items,
		TotalItems: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// CursorResult holds paginated results with cursor metadata.
type CursorResult[T any] struct {
	Items      []T
	NextCursor string
	PrevCursor string
	HasNext    bool
	HasPrev    bool
}

// NewCursorResult creates a CursorResult with the given items and cursor metadata.
func NewCursorResult[T any](items []T, nextCursor, prevCursor string, hasNext, hasPrev bool) CursorResult[T] {
	return CursorResult[T]{
		Items:      items,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}
