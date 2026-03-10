package ports

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

// ClosedPositionRecord pairs a closed Position with its exit price and close time.
type ClosedPositionRecord struct {
	Pos       *position.Position
	ExitPrice decimal.Decimal
	ClosedAt  time.Time
}

// PositionRepository persists and retrieves Position aggregates.
type PositionRepository interface {
	Save(ctx context.Context, pos *position.Position) error
	FindByID(ctx context.Context, positionID string) (*position.Position, error)
	FindByMarket(ctx context.Context, marketID string) ([]*position.Position, error)
	ListOpen(ctx context.Context) ([]*position.Position, error)
	ListClosed(ctx context.Context) ([]*position.Position, error)
	ListClosedWithExitPrice(ctx context.Context) ([]ClosedPositionRecord, error)
	Close(ctx context.Context, positionID string, exitPrice decimal.Decimal, closedAt time.Time) error
}
