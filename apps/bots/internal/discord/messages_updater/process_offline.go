package discordmessagesupdater

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/avast/retry-go/v4"
	"github.com/twirapp/twir/libs/logger"
)

func (c *MessagesUpdater) ProcessOffline(
	ctx context.Context,
	twitchChannelId string,
) error {
	messages, err := c.store.GetByChannelId(ctx, twitchChannelId)
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	for _, message := range messages {
		integration, err := c.discordRepo.GetByChannelIDAndGuildID(
			ctx,
			twitchChannelId,
			message.GuildID,
		)
		if err != nil {
			return err
		}

		if integration.IsNil() {
			continue
		}

		twitchUser, err := c.getTwitchUser(message.TwitchChannelID)
		if err != nil {
			return err
		}

		if integration.ShouldDeleteMessageOnOffline {
			err = retry.Do(
				func() error {
					return c.discord.DeleteMessage(
						ctx,
						message.DiscordChannelID,
						message.MessageID,
						"stream offline",
					)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error(
					"Failed to delete message",
					logger.Error(err),
					slog.Group(
						"discord",
						slog.String("message_id", message.MessageID),
						slog.String("guild_id", message.GuildID),
						slog.String("channel_id", message.DiscordChannelID),
					),
					slog.String("twitch_channel_id", message.TwitchChannelID),
				)

				if err := c.store.DeleteByMessageId(ctx, message.MessageID); err != nil {
					c.logger.Error(
						"Failed to delete sended message from store when we cannot delete discord message",
						logger.Error(err),
						slog.Group(
							"discord",
							slog.String("message_id", message.MessageID),
							slog.String("guild_id", message.GuildID),
							slog.String("channel_id", message.DiscordChannelID),
						),
						slog.String("twitch_channel_id", message.TwitchChannelID),
					)
				}

				continue
			}
		} else {
			content := integration.OfflineNotificationMessage
			if content == "" {
				content = "{userName} is offline now"
			}

			content = c.replaceMessageVars(
				content,
				replaceMessageVarsOpts{
					UserName:    twitchUser.Login,
					DisplayName: twitchUser.DisplayName,
				},
			)

			err = retry.Do(
				func() error {
					return c.discord.EditMessage(
						ctx,
						message.DiscordChannelID,
						message.MessageID,
						content,
					)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to edit message", logger.Error(err))
				continue
			}
		}
	}

	return nil
}
