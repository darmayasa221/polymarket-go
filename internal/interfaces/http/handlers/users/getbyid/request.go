// Package getbyid defines request/response DTOs for the get by ID action.
package getbyid

// URIParams holds URI path parameters for the get by ID action.
type URIParams struct {
	ID string `uri:"id" binding:"required"`
}
