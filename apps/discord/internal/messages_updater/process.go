package messages_updater

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessagesUpdater) process(ctx context.Context) {
	c.logger.Info("Start updating messages")

	messages, err := c.store.GetAll(ctx)
	if err != nil {
		c.logger.Error("Failed to get messages", slog.Any("err", err))
		return
	}

	streams := c.getStreams(ctx)
	streams = lo.Filter(
		streams,
		func(stream model.ChannelsStreams, _ int) bool {
			return stream.Channel != nil && stream.Channel.IsEnabled
		},
	)

	// DELETE messages that are not in the streams, so memory do not leak
	// for _, message := range messages {
	// 	_, ok := lo.Find(
	// 		streams,
	// 		func(stream model.ChannelsStreams) bool {
	// 			return stream.UserId == message.ChannelID
	// 		},
	// 	)
	// 	if !ok {
	// 		c.store.Delete(message.MessageID)
	// 	}
	// }

	for _, stream := range streams {
		_, err := c.store.Get(ctx, stream.UserId)
		if err != nil {
			continue
		}

		// TODO: send message, channel becomes online
		c.logger.Info(
			"Channel online, sending message",
			slog.Group(
				"channel",
				slog.String("id", stream.UserId),
				slog.String("name", stream.UserLogin),
				slog.String("title", stream.Title),
				slog.String("category", stream.GameName),
				slog.Int("viewers", stream.ViewerCount),
			),
		)

		onlineMessages, err := c.sendOnlineMessage(ctx, stream)
		if err != nil {
			c.logger.Error("Failed to send message", slog.Any("err", err))
			continue
		}

		for _, m := range onlineMessages {
			err = c.store.Add(ctx, m)
			if err != nil {
				c.logger.Error("Failed to add message to store", slog.Any("err", err))
				continue
			}
		}
	}

	for _, message := range messages {
		stream, ok := lo.Find(
			streams,
			func(stream model.ChannelsStreams) bool { return stream.UserId == message.ChannelID },
		)
		if !ok {
			if err = c.store.Delete(ctx, message.MessageID); err != nil {
				c.logger.Error("Failed to delete message from store", slog.Any("err", err))
			}
			// TODO: update message, channel becomes offline
			slog.Info("channel is offline", slog.String("id", message.ChannelID))
			continue
		}

		// TODO: update message, channel is still online, we need update values inside message embed
		c.logger.Info(
			"channel is still online",
			slog.Group(
				"channel",
				slog.String("id", stream.UserId),
				slog.String("name", stream.UserLogin),
				slog.String("title", stream.Title),
				slog.String("category", stream.GameName),
				slog.Int("viewers", stream.ViewerCount),
			),
		)
	}
}

func (c *MessagesUpdater) getStreams(
	ctx context.Context,
) []model.ChannelsStreams {
	var streams []model.ChannelsStreams
	if err := c.db.WithContext(ctx).Preload("Channel").Find(&streams).Error; err != nil {
		c.logger.Error("Failed to get streams", slog.Any("err", err))
		return nil
	}

	return streams
}
