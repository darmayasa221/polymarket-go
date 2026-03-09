package dto

import "time"

// UserItem is a single user in the list result.
type UserItem struct {
	ID        string
	Username  string
	Email     string
	FullName  string
	CreatedAt time.Time
}

// OffsetOutput holds offset-paginated user list results.
type OffsetOutput struct {
	Users      []UserItem
	TotalItems int
	Page       int
	PageSize   int
	TotalPages int
}

// CursorOutput holds cursor-paginated user list results.
type CursorOutput struct {
	Users      []UserItem
	NextCursor string
	PrevCursor string
	HasNext    bool
	HasPrev    bool
}
