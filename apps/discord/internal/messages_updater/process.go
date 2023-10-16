package messages_updater

import (
	"context"
	"errors"
	"log/slog"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *MessagesUpdater) process(ctx context.Context) {
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
	// 		c.store.Delete(ctx, message.MessageID)
	// 		// TODO: set offline message
	// 	}
	// }

	for _, stream := range streams {
		_, ok := lo.Find(
			messages,
			func(message sended_messages_store.Message) bool {
				return message.TwitchChannelID == stream.UserId
			},
		)
		if ok {
			continue
		}

		onlineMessages, err := c.sendOnlineMessage(ctx, stream)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.logger.Error("Failed to send message", slog.Any("err", err))
			}
			continue
		}

		if len(onlineMessages) == 0 {
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
			func(stream model.ChannelsStreams) bool { return stream.UserId == message.TwitchChannelID },
		)
		if !ok {
			if err = c.processOffline(ctx, message.TwitchChannelID); err != nil {
				c.logger.Error("Failed to process offline", slog.Any("err", err))
				continue
			}

			if err = c.store.DeleteByMessageId(ctx, message.MessageID); err != nil {
				c.logger.Error("Failed to delete message from store", slog.Any("err", err))
			}

			continue
		}

		if err = c.updateDiscordMessages(ctx, stream); err != nil {
			c.logger.Error("Failed to update message", slog.Any("err", err))
			continue
		}
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
