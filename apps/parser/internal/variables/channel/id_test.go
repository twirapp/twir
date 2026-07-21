package channel

import (
	"context"
	"testing"

	"github.com/twirapp/twir/apps/parser/internal/types"
)

func TestIDReturnsInternalChannelID(t *testing.T) {
	t.Parallel()

	result, err := ID.Handler(
		context.Background(),
		&types.VariableParseContext{
			ParseContext: &types.ParseContext{
				Channel: &types.ParseContextChannel{
					ID:          "platform-channel-id",
					DBChannelID: "internal-channel-id",
				},
			},
		},
		nil,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Result != "internal-channel-id" {
		t.Errorf("expected internal channel ID, got %q", result.Result)
	}
}
