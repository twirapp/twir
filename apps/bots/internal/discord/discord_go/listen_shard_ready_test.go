package discord_go

import (
	"io"
	"log/slog"
	"testing"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func TestHandleShardReadyWithoutShard(t *testing.T) {
	client := &Discord{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	client.handleShardReady(&gateway.ReadyEvent{})
}
