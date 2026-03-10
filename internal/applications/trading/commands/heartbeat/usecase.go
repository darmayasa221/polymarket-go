package heartbeat

import (
	"context"

	tradingcmds "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/heartbeat/dto"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	sender tradingports.HeartbeatSender
}

// New creates a Heartbeat use case.
func New(sender tradingports.HeartbeatSender) UseCase {
	return &useCase{sender: sender}
}

// Execute sends the CLOB keepalive POST.
// Must be called every 5 seconds or all open orders auto-cancel after 10 seconds.
func (uc *useCase) Execute(ctx context.Context, _ dto.Input) (dto.Output, error) {
	if err := uc.sender.Send(ctx); err != nil {
		return dto.Output{}, errtypes.NewInternalServerError(tradingcmds.ErrHeartbeatFailed)
	}
	return dto.Output{SentAt: timeutil.Now()}, nil
}
