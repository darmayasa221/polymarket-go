// Package dto defines data transfer objects for the listusers query.
package dto

import "github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"

// Input holds pagination parameters for listing users.
type Input struct {
	Mode         pagination.Mode
	OffsetParams pagination.OffsetParams
	CursorParams pagination.CursorParams
}
