package messages_updater

import (
	"context"
	"errors"
	"log/slog"

	"github.com/avast/retry-go/v4"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/switchupcb/disgo"
)

func (c *MessagesUpdater) updateDiscordMessages(
	ctx context.Context,
	stream model.ChannelsStreams,
) error {
	settings, err := c.getChannelDiscordIntegration(ctx, stream.UserId)
	if err != nil {
		return err
	}

	if settings.Data.Discord == nil || len(settings.Data.Discord.Guilds) == 0 {
		return nil
	}

	twitchUsersReq, err := c.twitchClient.GetUsers(&helix.UsersParams{IDs: []string{stream.UserId}})
	if len(twitchUsersReq.Data.Users) == 0 {
		return errors.New("user not found")
	}
	twitchUser := twitchUsersReq.Data.Users[0]

	for _, guild := range settings.Data.Discord.Guilds {
		if !guild.LiveNotificationEnabled {
			continue
		}

		messages, err := c.store.GetByGuildId(ctx, guild.ID)
		if err != nil {
			return err
		}

		for _, message := range messages {
			embed := c.buildEmbed(twitchUser, stream, guild)

			editMsg := disgo.EditMessage{
				Embeds:    &[]*disgo.Embed{embed},
				MessageID: message.MessageID,
				ChannelID: message.DiscordChannelID,
			}

			_, err := retry.DoWithData(
				func() (*disgo.Message, error) {
					return editMsg.Send(c.discord.Client)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to edit message", slog.Any("err", err))
				continue
			}
		}
	}

	return nil
}
