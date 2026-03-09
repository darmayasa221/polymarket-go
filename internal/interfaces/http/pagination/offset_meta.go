package pagination

// OffsetMeta is the HTTP response metadata for offset pagination.
type OffsetMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// NewOffsetMeta creates HTTP pagination metadata from domain result values.
func NewOffsetMeta(page, pageSize, totalItems, totalPages int) OffsetMeta {
	return OffsetMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
