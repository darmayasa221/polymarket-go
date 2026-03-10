package heartbeat_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/heartbeat"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/heartbeat/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// mockHeartbeatSender implements ports.HeartbeatSender.
type mockHeartbeatSender struct {
	sendErr error
	called  bool
}

func (m *mockHeartbeatSender) Send(_ context.Context) error {
	m.called = true
	return m.sendErr
}

func TestHeartbeat_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		sendErr   error
		wantErr   bool
		errTarget any
	}{
		{
			name:    "heartbeat sent successfully",
			wantErr: false,
		},
		{
			name:      "send failure returns internal error",
			sendErr:   errors.New("CLOB unreachable"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sender := &mockHeartbeatSender{sendErr: tt.sendErr}
			uc := heartbeat.New(sender)

			out, err := uc.Execute(t.Context(), dto.Input{})

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
				assert.True(t, sender.called)
			} else {
				require.NoError(t, err)
				assert.True(t, sender.called)
				assert.False(t, out.SentAt.IsZero())
			}
		})
	}
}
