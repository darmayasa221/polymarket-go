package pagination

// CursorMeta is the HTTP response metadata for cursor pagination.
type CursorMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
}

// NewCursorMeta creates cursor pagination metadata from domain result values.
func NewCursorMeta(nextCursor, prevCursor string, hasNext, hasPrev bool) CursorMeta {
	return CursorMeta{
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}
