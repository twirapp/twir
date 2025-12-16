package messagehandler

import (
	"context"

	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleGamesVoteban(ctx context.Context, msg handleMessage) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	return c.votebanService.HandleTwitchMessage(ctx, msg.TwitchChatMessage)
}
