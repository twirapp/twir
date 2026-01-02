package discordmessagesupdater

import (
	"context"

	"github.com/avast/retry-go/v4"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

func (c *MessagesUpdater) updateDiscordMessages(
	ctx context.Context,
	stream model.ChannelsStreams,
) error {
	integrations, err := c.getChannelDiscordIntegrations(ctx, stream.UserId)
	if err != nil {
		return err
	}

	if len(integrations) == 0 {
		return nil
	}

	for _, integration := range integrations {
		if !integration.LiveNotificationEnabled {
			continue
		}

		messages, err := c.store.GetByGuildId(ctx, integration.GuildID)
		if err != nil {
			return err
		}

		for _, message := range messages {
			if message.TwitchChannelID != stream.UserId {
				continue
			}

			twitchUser, err := c.getTwitchUser(stream.UserId)
			if err != nil {
				c.logger.Error("Failed to get twitch user", logger.Error(err))
				continue
			}

			embed := c.buildEmbed(twitchUser, stream, integration)

			err = retry.Do(
				func() error {
					content := c.replaceMessageVars(
						integration.LiveNotificationMessage, replaceMessageVarsOpts{
							UserName:     stream.UserLogin,
							DisplayName:  stream.UserName,
							CategoryName: stream.GameName,
							Title:        stream.Title,
						},
					)

					return c.discord.EditMessage(
						ctx,
						message.DiscordChannelID,
						message.MessageID,
						content,
						embed,
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
