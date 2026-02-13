package messagehandler

import (
	"context"
	"strings"

	"github.com/twirapp/twir/libs/bus-core/twitch"
	channels_giveaways "github.com/twirapp/twir/libs/entities/channels_giveaways"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleGiveaways(ctx context.Context, msg twitch.TwitchChatMessage) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	giveaways, err := c.giveawaysCacher.Get(ctx, msg.BroadcasterUserId)
	if err != nil {
		return err
	}

	if len(giveaways) == 0 {
		return nil
	}

	/*
		Here we can try to check info in database but for premium users only.
	*/
	for _, giveaway := range giveaways {
		if giveaway.IsNil() {
			continue
		}

		if giveaway.StartedAt == nil {
			continue
		}

		if giveaway.StoppedAt != nil {
			continue
		}

		// Skip keyword-based giveaways if this is an ONLINE_CHATTERS giveaway
		if giveaway.Type == channels_giveaways.GiveawayTypeOnlineChatters {
			continue
		}

		// For KEYWORD giveaways, check if message contains the keyword
		if giveaway.Keyword == nil || !strings.Contains(strings.ToLower(msg.Message.Text), strings.ToLower(*giveaway.Keyword)) {
			continue
		}

		err := c.giveawaysService.TryAddParticipant(
			ctx,
			msg.ChatterUserId,
			msg.ChatterUserLogin,
			msg.ChatterUserName,
			giveaway.ID.String(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
