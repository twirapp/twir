package messagehandler

import (
	"context"
	"log/slog"

	channelsemotesusages "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleEmotesUsagesBatched(ctx context.Context, data []handleMessage) {
	var createEmoteUsageInputs []channelsemotesusages.ChannelEmoteUsageInput

	for _, msg := range data {
		for key, count := range msg.EnrichedData.UsedEmotesWithThirdParty {
			for i := 0; i < count; i++ {
				createEmoteUsageInputs = append(
					createEmoteUsageInputs,
					channelsemotesusages.ChannelEmoteUsageInput{
						ChannelID: msg.BroadcasterUserId,
						UserID:    msg.ChatterUserId,
						Emote:     key,
					},
				)
			}
		}
	}

	err := c.channelsEmotesUsagesRepository.CreateMany(ctx, createEmoteUsageInputs)
	if err != nil {
		c.logger.Error("cannot create emotes usages", slog.Any("err", err))
	}
}

func (c *MessageHandler) handleEmotesUsages(ctx context.Context, msg handleMessage) error {
	span := trace.SpanFromContext(ctx)
  defer span.End()
  span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	c.messagesEmotesBatcher.Add(msg)
	return nil
}
