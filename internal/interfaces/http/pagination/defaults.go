package pagination

// String-form defaults used by DefaultQuery() calls in the HTTP layer.
// The domain's integer constants (pagination.DefaultPage, DefaultPageSize) are
// the authoritative values; these are the equivalent string representations.
const (
	defaultPageStr      = "1"
	defaultPageSizeStr  = "10"
	defaultDirectionStr = "forward"
)
