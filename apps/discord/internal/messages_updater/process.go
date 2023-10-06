package messages_updater

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessagesUpdater) process(ctx context.Context) {
	c.logger.Info("Start updating messages")

	messages := c.store.GetAll()
	streams := c.getStreams(ctx)

	for _, stream := range streams {
		// TODO: update message, channel becomes online
		fmt.Println(stream)
	}

	for _, message := range messages {
		stream, ok := lo.Find(
			streams,
			func(stream model.ChannelsStreams) bool { return stream.UserId == message.ChannelID },
		)
		if !ok {
			c.store.Delete(message.MessageID)
			// TODO: update message, channel becomes offline
			continue
		}

		// TODO: update message, channel is still online, we need update values inside message embed
		fmt.Println(stream)
	}
}

func (c *MessagesUpdater) getStreams(
	ctx context.Context,
) []model.ChannelsStreams {
	var streams []model.ChannelsStreams
	if err := c.db.WithContext(ctx).Find(&streams).Error; err != nil {
		c.logger.Error("Failed to get streams", slog.Any("err", err))
		return nil
	}

	return streams
}
